package capture

import (
	"fmt"
	"sync"
	"syscall"
	"unsafe"

	"github.com/shinkar94/godesktopdup/disp"
	resultcode "github.com/shinkar94/godesktopdup/errors"
)

var (
	cursorDLLOnce sync.Once
	cursorDLL     *syscall.LazyDLL
	cursorProc    *syscall.LazyProc
)

func initCursorAPI() {
	cursorDLLOnce.Do(func() {
		cursorDLL = syscall.NewLazyDLL("user32.dll")
		cursorProc = cursorDLL.NewProc("GetCursorInfo")
	})
}

// getGlobalCursorPos retrieves global cursor position via Windows API.
func getGlobalCursorPos() (x, y int, visible bool, err error) {
	initCursorAPI()

	type POINT struct {
		X int32
		Y int32
	}

	type CURSORINFO struct {
		cbSize      uint32
		flags       uint32
		hCursor     uintptr
		ptScreenPos POINT
	}

	const CURSOR_SHOWING = 0x00000001

	var ci CURSORINFO
	ci.cbSize = uint32(unsafe.Sizeof(ci))

	ret, _, _ := cursorProc.Call(uintptr(unsafe.Pointer(&ci)))
	if ret == 0 {
		return 0, 0, false, fmt.Errorf("GetCursorInfo failed")
	}

	return int(ci.ptScreenPos.X), int(ci.ptScreenPos.Y), (ci.flags & CURSOR_SHOWING) != 0, nil
}

// getCursorShape retrieves cursor shape from current frame.
func (sc *ScreenCapture) getCursorShape() error {
	if sc.outputDuplication == nil {
		return fmt.Errorf("outputDuplication is nil")
	}

	if sc.currentFrameInfo.PointerShapeBufferSize == 0 {
		if sc.cursorShapeInfo.Width > 0 && sc.cursorShapeInfo.Height > 0 {
			return nil
		}
		return nil
	}

	bufferSize := sc.currentFrameInfo.PointerShapeBufferSize
	if cap(sc.cursorShapeBuffer) < int(bufferSize) {
		sc.cursorShapeBuffer = make([]byte, bufferSize)
	} else {
		sc.cursorShapeBuffer = sc.cursorShapeBuffer[:bufferSize]
	}

	var shapeInfo disp.DuplicationPointerShapeInfo
	var bufferSizeRequired uint32

	hr := sc.outputDuplication.GetFramePointerShape(
		bufferSize,
		sc.cursorShapeBuffer,
		&bufferSizeRequired,
		&shapeInfo,
	)

	if hrCode := resultcode.ResultCode(hr); hrCode.Failed() {
		if hrCode == resultcode.ErrorMoreData {
			sc.cursorShapeBuffer = make([]byte, bufferSizeRequired)
			hr2 := sc.outputDuplication.GetFramePointerShape(
				bufferSizeRequired,
				sc.cursorShapeBuffer,
				&bufferSizeRequired,
				&shapeInfo,
			)
			if hrCode2 := resultcode.ResultCode(hr2); hrCode2.Failed() {
				return fmt.Errorf("failed to GetFramePointerShape: %w", hrCode2)
			}
		} else {
			return fmt.Errorf("failed to GetFramePointerShape: %w", hrCode)
		}
	}

	sc.cursorShapeInfo = shapeInfo

	return nil
}

// drawCursor draws cursor on frame.
func (sc *ScreenCapture) drawCursor(buffer []byte, width, height int) error {
	desktopCursorX, desktopCursorY, visible, err := getGlobalCursorPos()
	if err != nil || !visible {
		return nil
	}

	var bounds disp.Rect
	if sc.monitorBounds != nil {
		bounds = *sc.monitorBounds
	} else {
		var err error
		bounds, err = sc.GetBounds()
		if err != nil {
			return nil
		}
	}

	boundsLeft := int(bounds.Left)
	boundsTop := int(bounds.Top)
	boundsRight := int(bounds.Right)
	boundsBottom := int(bounds.Bottom)

	if desktopCursorX < boundsLeft || desktopCursorX >= boundsRight ||
		desktopCursorY < boundsTop || desktopCursorY >= boundsBottom {
		return nil
	}

	if sc.cursorShapeInfo.Width == 0 || sc.cursorShapeInfo.Height == 0 {
		return nil
	}

	cursorX := desktopCursorX - boundsLeft
	cursorY := desktopCursorY - boundsTop

	cursorWidth := int(sc.cursorShapeInfo.Width)
	cursorHeight := int(sc.cursorShapeInfo.Height)
	cursorPitch := int(sc.cursorShapeInfo.Pitch)
	hotSpotX := int(sc.cursorShapeInfo.HotSpot.X)
	hotSpotY := int(sc.cursorShapeInfo.HotSpot.Y)

	startX := cursorX - hotSpotX
	startY := cursorY - hotSpotY

	if startX+cursorWidth <= 0 || startX >= width || startY+cursorHeight <= 0 || startY >= height {
		return nil
	}

	switch sc.cursorShapeInfo.Type {
	case disp.DuplicationPointerShapeTypeMonochrome:
		return sc.drawMonochromeCursor(buffer, width, height, startX, startY, cursorWidth, cursorHeight, cursorPitch)
	case disp.DuplicationPointerShapeTypeColor:
		return sc.drawColorCursor(buffer, width, height, startX, startY, cursorWidth, cursorHeight, cursorPitch)
	case disp.DuplicationPointerShapeTypeMaskedColor:
		return sc.drawMaskedColorCursor(buffer, width, height, startX, startY, cursorWidth, cursorHeight, cursorPitch)
	default:
		return nil
	}
}

// drawMonochromeCursor draws monochrome cursor (AND mask + XOR mask).
func (sc *ScreenCapture) drawMonochromeCursor(buffer []byte, width, height, startX, startY, cursorWidth, cursorHeight, cursorPitch int) error {
	andMaskPitch := (cursorWidth + 7) / 8
	andMaskSize := andMaskPitch * cursorHeight
	xorMaskOffset := andMaskSize

	andMask := sc.cursorShapeBuffer[:andMaskSize]
	xorMask := sc.cursorShapeBuffer[xorMaskOffset : xorMaskOffset+andMaskSize]

	clipTop := 0
	if startY < 0 {
		clipTop = -startY
	}
	clipBottom := cursorHeight
	if startY+cursorHeight > height {
		clipBottom = height - startY
	}
	clipLeft := 0
	if startX < 0 {
		clipLeft = -startX
	}
	clipRight := cursorWidth
	if startX+cursorWidth > width {
		clipRight = width - startX
	}

	bufU32 := (*[1 << 30]uint32)(unsafe.Pointer(&buffer[0]))[:len(buffer)/4]
	widthU32 := width

	andMaskRowStart := clipTop * andMaskPitch
	for y := clipTop; y < clipBottom; y++ {
		frameY := startY + y
		dstRowStart := frameY*widthU32 + startX
		andMaskByteStart := andMaskRowStart + clipLeft/8

		for x := clipLeft; x < clipRight; x++ {
			byteIdx := andMaskByteStart + x/8
			bitIdx := 7 - (x & 7)

			andBit := (andMask[byteIdx] >> bitIdx) & 1
			if andBit == 0 {
				continue
			}

			xorBit := (xorMask[byteIdx] >> bitIdx) & 1

			dstOffset := dstRowStart + x
			if xorBit == 0 {
				bufU32[dstOffset] = 0x000000FF
			} else {
				bufU32[dstOffset] = 0xFFFFFFFF
			}
		}
		andMaskRowStart += andMaskPitch
	}

	return nil
}

// drawColorCursor draws color cursor with alpha blending.
func (sc *ScreenCapture) drawColorCursor(buffer []byte, width, height, startX, startY, cursorWidth, cursorHeight, cursorPitch int) error {
	clipTop := 0
	if startY < 0 {
		clipTop = -startY
	}
	clipBottom := cursorHeight
	if startY+cursorHeight > height {
		clipBottom = height - startY
	}
	clipLeft := 0
	if startX < 0 {
		clipLeft = -startX
	}
	clipRight := cursorWidth
	if startX+cursorWidth > width {
		clipRight = width - startX
	}

	bufU32 := (*[1 << 30]uint32)(unsafe.Pointer(&buffer[0]))[:len(buffer)/4]
	cursorU32 := (*[1 << 30]uint32)(unsafe.Pointer(&sc.cursorShapeBuffer[0]))[:len(sc.cursorShapeBuffer)/4]
	widthU32 := width
	cursorPitchU32 := cursorPitch / 4
	
	maxCursorBufferU32 := len(sc.cursorShapeBuffer) / 4

	for y := clipTop; y < clipBottom; y++ {
		frameY := startY + y
		dstRowStart := frameY*widthU32 + startX
		cursorRowStart := y * cursorPitchU32

		maxCursorOffsetU32 := cursorRowStart + clipRight
		if maxCursorOffsetU32 > maxCursorBufferU32 {
			maxCursorOffsetU32 = maxCursorBufferU32
		}

		for x := clipLeft; x < clipRight; x++ {
			cursorOffsetU32 := cursorRowStart + x
			if cursorOffsetU32 >= maxCursorOffsetU32 {
				continue
			}

			cursorPixel := cursorU32[cursorOffsetU32]
			a := byte(cursorPixel >> 24)
			if a == 0 {
				continue
			}

			dstOffset := dstRowStart + x
			
			if a == 255 {
				bufU32[dstOffset] = cursorPixel | 0xFF000000
			} else {
				alpha := uint16(a)
				invAlpha := 255 - alpha
				bgPixel := bufU32[dstOffset]
				
				bgB := uint16(bgPixel & 0xFF)
				bgG := uint16((bgPixel >> 8) & 0xFF)
				bgR := uint16((bgPixel >> 16) & 0xFF)
				
				curB := uint16(cursorPixel & 0xFF)
				curG := uint16((cursorPixel >> 8) & 0xFF)
				curR := uint16((cursorPixel >> 16) & 0xFF)
				
				newB := byte((bgB*invAlpha + curB*alpha) / 255)
				newG := byte((bgG*invAlpha + curG*alpha) / 255)
				newR := byte((bgR*invAlpha + curR*alpha) / 255)
				
				bufU32[dstOffset] = uint32(newB) | (uint32(newG) << 8) | (uint32(newR) << 16) | 0xFF000000
			}
		}
	}

	return nil
}

// drawMaskedColorCursor draws masked color cursor (XOR mask + AND mask).
func (sc *ScreenCapture) drawMaskedColorCursor(buffer []byte, width, height, startX, startY, cursorWidth, cursorHeight, cursorPitch int) error {
	xorMaskPitch := cursorPitch
	xorMaskSize := xorMaskPitch * cursorHeight
	andMaskPitch := (cursorWidth + 7) / 8
	andMaskSize := andMaskPitch * cursorHeight
	andMaskOffset := xorMaskSize

	xorMask := sc.cursorShapeBuffer[:xorMaskSize]
	andMask := sc.cursorShapeBuffer[andMaskOffset : andMaskOffset+andMaskSize]

	clipTop := 0
	if startY < 0 {
		clipTop = -startY
	}
	clipBottom := cursorHeight
	if startY+cursorHeight > height {
		clipBottom = height - startY
	}
	clipLeft := 0
	if startX < 0 {
		clipLeft = -startX
	}
	clipRight := cursorWidth
	if startX+cursorWidth > width {
		clipRight = width - startX
	}

	bufU32 := (*[1 << 30]uint32)(unsafe.Pointer(&buffer[0]))[:len(buffer)/4]
	xorU32 := (*[1 << 30]uint32)(unsafe.Pointer(&xorMask[0]))[:len(xorMask)/4]
	widthU32 := width
	xorMaskPitchU32 := xorMaskPitch / 4

	andMaskRowStart := clipTop * andMaskPitch
	for y := clipTop; y < clipBottom; y++ {
		frameY := startY + y
		dstRowStart := frameY*widthU32 + startX
		cursorRowStart := y * xorMaskPitchU32
		andMaskByteStart := andMaskRowStart + clipLeft/8

		for x := clipLeft; x < clipRight; x++ {
			byteIdx := andMaskByteStart + x/8
			bitIdx := 7 - (x & 7)
			andBit := (andMask[byteIdx] >> bitIdx) & 1

			if andBit == 0 {
				continue
			}

			dstOffset := dstRowStart + x
			cursorOffset := cursorRowStart + x
			bufU32[dstOffset] = xorU32[cursorOffset]
		}
		andMaskRowStart += andMaskPitch
	}

	return nil
}
