package capture

import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/shinkar94/godesktopdup/disp"
	resultcode "github.com/shinkar94/godesktopdup/errors"
	"github.com/shinkar94/godesktopdup/gfx11"
)

var ErrNoImageYet = errors.New("no image yet")

type ScreenCapture struct {
	device            *gfx11.Device
	deviceCtx         *gfx11.DeviceContext
	outputDuplication *disp.OutputDuplication
	dxgiOutput        *disp.Output5

	stagedTex  *gfx11.Texture2D
	surface    *disp.Surface
	mappedRect disp.MappedRect
	size       disp.Point

	dirtyRects    []disp.Rect
	movedRects    []disp.DuplicationMoveRect
	acquiredFrame bool
	needsSwizzle  bool
	rotation      disp.ModeRotation
	logicalWidth  int32
	logicalHeight int32

	currentFrameInfo  disp.DuplicationFrameInfo
	cursorShapeBuffer []byte
	cursorShapeInfo   disp.DuplicationPointerShapeInfo

	monitorBounds *disp.Rect
	captureCursor bool

	lastOutputPtr    uintptr
	frameInitialized bool

	cachedIsVertical  bool
	cachedContentWidth int
	cachedDataWidth    int
	cachedSize         disp.Point
}

func (sc *ScreenCapture) initializeStage(texture *gfx11.Texture2D) error {
	desc := gfx11.Texture2DDesc{}
	hr := texture.GetDesc(&desc)
	if resultcode.ResultCode(hr).Failed() {
		return resultcode.ResultCode(hr)
	}

	desc.Usage = gfx11.UsageStaging
	desc.CPUAccessFlags = gfx11.CPUAccessRead
	desc.BindFlags = 0
	desc.MipLevels = 1
	desc.ArraySize = 1
	desc.MiscFlags = 0
	desc.SampleDesc.Count = 1

	hr = sc.device.CreateTexture2D(&desc, &sc.stagedTex)
	if resultcode.ResultCode(hr).Failed() {
		return resultcode.ResultCode(hr)
	}

	hr = sc.stagedTex.QueryInterface(disp.IID_Surface, &sc.surface)
	if resultcode.ResultCode(hr).Failed() {
		return resultcode.ResultCode(hr)
	}
	sc.size = disp.Point{X: int32(desc.Width), Y: int32(desc.Height)}

	return nil
}

func (sc *ScreenCapture) Release() {
	sc.ReleaseFrame()
	if sc.stagedTex != nil {
		sc.stagedTex.Release()
		sc.stagedTex = nil
	}
	if sc.surface != nil {
		sc.surface.Release()
		sc.surface = nil
	}
	if sc.outputDuplication != nil {
		sc.outputDuplication.Release()
		sc.outputDuplication = nil
	}
	if sc.dxgiOutput != nil {
		sc.dxgiOutput.Release()
		sc.dxgiOutput = nil
	}
	sc.lastOutputPtr = 0
	sc.frameInitialized = false
}

func (sc *ScreenCapture) ReleaseFrame() {
	if sc.acquiredFrame {
		sc.outputDuplication.ReleaseFrame()
		sc.acquiredFrame = false
	}
}

func (sc *ScreenCapture) Snapshot(timeoutMs uint) (func() int32, *disp.MappedRect, *disp.Point, error) {
	if sc.outputDuplication == nil {
		return nil, nil, nil, fmt.Errorf("outputDuplication is nil")
	}

	var hr int32
	desc := disp.DuplicationDesc{}
	hr = sc.outputDuplication.GetDesc(&desc)
	if hr := resultcode.ResultCode(hr); hr.Failed() {
		return nil, nil, nil, fmt.Errorf("failed to get the description. %w", hr)
	}

	width := int32(desc.ModeDesc.Width)
	height := int32(desc.ModeDesc.Height)
	rotation := disp.ModeRotation(desc.Rotation)

	if sc.rotation == 0 {
		sc.rotation = rotation
	}

	if desc.DesktopImageInSystemMemory != 0 {
		sc.size = disp.Point{X: width, Y: height}
		hr = sc.outputDuplication.MapDesktopSurface(&sc.mappedRect)
		if hr := resultcode.ResultCode(hr); !hr.Failed() {
			return sc.outputDuplication.UnMapDesktopSurface, &sc.mappedRect, &sc.size, nil
		}
	}

	var desktop *disp.Resource
	var frameInfo disp.DuplicationFrameInfo

	sc.ReleaseFrame()
	hrF := sc.outputDuplication.AcquireNextFrame(uint32(timeoutMs), &frameInfo, &desktop)
	sc.acquiredFrame = true
	if hr := resultcode.ResultCode(hrF); hr.Failed() {
		if hr == resultcode.ErrorWaitTimeout {
			return nil, nil, nil, ErrNoImageYet
		}
		return nil, nil, nil, fmt.Errorf("failed to AcquireNextFrame. %w", resultcode.ResultCode(hrF))
	}

	defer sc.ReleaseFrame()
	defer desktop.Release()

	if frameInfo.AccumulatedFrames == 0 {
		return nil, nil, nil, ErrNoImageYet
	}

	sc.currentFrameInfo = frameInfo

	if sc.captureCursor && frameInfo.PointerShapeBufferSize > 0 {
		if err := sc.getCursorShape(); err != nil {
			_ = err
		}
	}

	var desktop2d *gfx11.Texture2D
	hr = desktop.QueryInterface(gfx11.IID_Texture2D, &desktop2d)
	if hr := resultcode.ResultCode(hr); hr.Failed() {
		return nil, nil, nil, fmt.Errorf("failed to QueryInterface(iid_Texture2D, ...). %w", hr)
	}
	defer desktop2d.Release()

	if sc.stagedTex == nil {
		err := sc.initializeStage(desktop2d)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to InitializeStage. %w", err)
		}
	}

	if frameInfo.TotalMetadataBufferSize > 0 {
		moveRectsRequired := uint32(1)
		for {
			if cap(sc.movedRects) < int(moveRectsRequired) {
				sc.movedRects = make([]disp.DuplicationMoveRect, moveRectsRequired)
			} else {
				sc.movedRects = sc.movedRects[:moveRectsRequired]
			}
			hr = sc.outputDuplication.GetFrameMoveRects(sc.movedRects, &moveRectsRequired)
			if hr := resultcode.ResultCode(hr); hr.Failed() {
				if hr == resultcode.ErrorMoreData {
					continue
				}
				return nil, nil, nil, fmt.Errorf("failed to GetFrameMoveRects. %w", hr)
			}
			sc.movedRects = sc.movedRects[:moveRectsRequired]
			break
		}

		dirtyRectsRequired := uint32(1)
		for {
			if cap(sc.dirtyRects) < int(dirtyRectsRequired) {
				sc.dirtyRects = make([]disp.Rect, dirtyRectsRequired)
			} else {
				sc.dirtyRects = sc.dirtyRects[:dirtyRectsRequired]
			}
			hr = sc.outputDuplication.GetFrameDirtyRects(sc.dirtyRects, &dirtyRectsRequired)
			if hr := resultcode.ResultCode(hr); hr.Failed() {
				if hr == resultcode.ErrorMoreData {
					continue
				}
				return nil, nil, nil, fmt.Errorf("failed to GetFrameDirtyRects. %w", hr)
			}
			sc.dirtyRects = sc.dirtyRects[:dirtyRectsRequired]
			break
		}

		box := gfx11.Box{
			Front: 0,
			Back:  1,
		}
		if len(sc.movedRects) == 0 {
			for _, rect := range sc.dirtyRects {
				box.Left = uint32(rect.Left)
				box.Top = uint32(rect.Top)
				box.Right = uint32(rect.Right)
				box.Bottom = uint32(rect.Bottom)

				sc.deviceCtx.CopySubresourceRegion2D(sc.stagedTex, 0, box.Left, box.Top, 0, desktop2d, 0, &box)
			}
		} else {
			sc.deviceCtx.CopyResource2D(sc.stagedTex, desktop2d)
		}
	} else {
		sc.deviceCtx.CopyResource2D(sc.stagedTex, desktop2d)
		if !sc.needsSwizzle {
			sc.needsSwizzle = true
		}
		sc.dirtyRects = sc.dirtyRects[:0]
		sc.movedRects = sc.movedRects[:0]
	}

	hr = sc.surface.Map(&sc.mappedRect, disp.MapRead)
	if hr := resultcode.ResultCode(hr); hr.Failed() {
		return nil, nil, nil, fmt.Errorf("failed to surface.Map(...). %v", hr)
	}
	return sc.surface.Unmap, &sc.mappedRect, &sc.size, nil
}

func (sc *ScreenCapture) GetFrameBGRA(buffer []byte, timeoutMs uint) error {
	if sc.outputDuplication == nil {
		return fmt.Errorf("outputDuplication is nil before Snapshot call")
	}

	unmap, mappedRect, size, err := sc.Snapshot(timeoutMs)
	if err != nil {
		return err
	}
	defer unmap()

	dataSize := int(mappedRect.Pitch) * int(size.Y)
	data := unsafe.Slice((*byte)(mappedRect.PBits), dataSize)

	var contentWidth, dataWidth int
	var isVertical bool
	if sc.cachedSize.X == size.X && sc.cachedSize.Y == size.Y {
		contentWidth = sc.cachedContentWidth
		dataWidth = sc.cachedDataWidth
		isVertical = sc.cachedIsVertical
	} else {
		contentWidth = int(size.X) * 4
		dataWidth = int(mappedRect.Pitch)
		isVertical = sc.logicalWidth > 0 && sc.logicalHeight > 0 && sc.logicalWidth < sc.logicalHeight
		
		sc.cachedContentWidth = contentWidth
		sc.cachedDataWidth = dataWidth
		sc.cachedIsVertical = isVertical
		sc.cachedSize = *size
	}

	if len(buffer) == 0 {
		return fmt.Errorf("buffer too small")
	}
	outputPtr := uintptr(unsafe.Pointer(&buffer[0]))
	if sc.lastOutputPtr != outputPtr {
		sc.frameInitialized = false
		sc.cachedSize = disp.Point{}
	}

	physicalWidth := int(size.X)
	physicalHeight := int(size.Y)
	logicalWidthInt := int(sc.logicalWidth)
	logicalHeightInt := int(sc.logicalHeight)

	if isVertical && (physicalWidth != logicalWidthInt || physicalHeight != logicalHeightInt) {
		if err := sc.copyRotatedFrame(buffer, data, *size, mappedRect.Pitch, sc.rotation, sc.logicalWidth, sc.logicalHeight); err != nil {
			return err
		}
		if sc.captureCursor {
			if err := sc.drawCursor(buffer, logicalWidthInt, logicalHeightInt); err != nil {
				_ = err
			}
		}
		sc.frameInitialized = true
		sc.lastOutputPtr = outputPtr
		return nil
	}

	hasDirtyRects := len(sc.dirtyRects) > 0
	hasMovedRects := len(sc.movedRects) > 0
	needFullCopy := !sc.frameInitialized || !hasDirtyRects || hasMovedRects
	if needFullCopy {
		if err := sc.copyFullFrame(buffer, data, *size, contentWidth, dataWidth); err != nil {
			return err
		}
	} else {
		if err := sc.copyDirtyRegions(buffer, data, *size, contentWidth, dataWidth); err != nil {
			if errFull := sc.copyFullFrame(buffer, data, *size, contentWidth, dataWidth); errFull != nil {
				return errFull
			}
		}
	}

	if sc.captureCursor {
		if err := sc.drawCursor(buffer, physicalWidth, physicalHeight); err != nil {
			_ = err
		}
	}

	sc.frameInitialized = true
	sc.lastOutputPtr = outputPtr

	return nil
}

func (sc *ScreenCapture) copyFullFrame(buffer []byte, data []byte, size disp.Point, contentWidth, dataWidth int) error {
	requiredSize := contentWidth * int(size.Y)
	if len(buffer) < requiredSize {
		return fmt.Errorf("buffer too small")
	}

	height := int(size.Y)
	if contentWidth == dataWidth {
		copy(buffer[:requiredSize], data[:requiredSize])
	} else {
		if contentWidth%4 == 0 && dataWidth%4 == 0 {
			pixelsPerRow := contentWidth / 4
			dataWidthU32 := dataWidth / 4
			bufU32 := (*[1 << 30]uint32)(unsafe.Pointer(&buffer[0]))[:len(buffer)/4]
			dataU32 := (*[1 << 30]uint32)(unsafe.Pointer(&data[0]))[:len(data)/4]
			
			imgStart := 0
			dataStart := 0
			for i := 0; i < height; i++ {
				copy(bufU32[imgStart:imgStart+pixelsPerRow], dataU32[dataStart:dataStart+pixelsPerRow])
				imgStart += pixelsPerRow
				dataStart += dataWidthU32
			}
		} else {
			imgStart := 0
			dataStart := 0
			for i := 0; i < height; i++ {
				copy(buffer[imgStart:imgStart+contentWidth], data[dataStart:dataStart+contentWidth])
				imgStart += contentWidth
				dataStart += dataWidth
			}
		}
	}
	return nil
}

func (sc *ScreenCapture) copyDirtyRegions(buffer []byte, data []byte, size disp.Point, contentWidth, dataWidth int) error {
	requiredSize := contentWidth * int(size.Y)
	if len(buffer) < requiredSize {
		return fmt.Errorf("buffer too small")
	}

	maxWidth := int(size.X)
	maxHeight := int(size.Y)

	bufU32 := (*[1 << 30]uint32)(unsafe.Pointer(&buffer[0]))[:len(buffer)/4]
	dataU32 := (*[1 << 30]uint32)(unsafe.Pointer(&data[0]))[:len(data)/4]
	
	dataWidthU32 := dataWidth / 4
	contentWidthU32 := contentWidth / 4

	for i := range sc.dirtyRects {
		rect := sc.dirtyRects[i]
		left := clampInt(int(rect.Left), 0, maxWidth)
		right := clampInt(int(rect.Right), 0, maxWidth)
		top := clampInt(int(rect.Top), 0, maxHeight)
		bottom := clampInt(int(rect.Bottom), 0, maxHeight)

		if left >= right || top >= bottom {
			continue
		}

		rowBytes := (right - left) * 4
		leftTimes4 := left * 4
		srcRowStart := top*dataWidth + leftTimes4
		dstRowStart := top*contentWidth + leftTimes4
		rowCount := bottom - top

		if rowBytes%4 == 0 {
			pixelsPerRow := rowBytes / 4
			srcRowStartU32 := srcRowStart / 4
			dstRowStartU32 := dstRowStart / 4
			for y := 0; y < rowCount; y++ {
				srcRow := srcRowStartU32 + y*dataWidthU32
				dstRow := dstRowStartU32 + y*contentWidthU32
				copy(bufU32[dstRow:dstRow+pixelsPerRow], dataU32[srcRow:srcRow+pixelsPerRow])
			}
		} else {
			for y := 0; y < rowCount; y++ {
				srcRow := srcRowStart + y*dataWidth
				dstRow := dstRowStart + y*contentWidth
				copy(buffer[dstRow:dstRow+rowBytes], data[srcRow:srcRow+rowBytes])
			}
		}
	}

	return nil
}

func clampInt(value, minVal, maxVal int) int {
	if value < minVal {
		return minVal
	}
	if value > maxVal {
		return maxVal
	}
	return value
}

// copyRotatedFrame transposes frame data for rotated monitors.
// DDA returns data in physical panel orientation (e.g., 1920x1080),
// but we need logical orientation (e.g., 1080x1920).
func (sc *ScreenCapture) copyRotatedFrame(dst, src []byte, physicalSize disp.Point, pitch int32, rotation disp.ModeRotation, logicalWidth, logicalHeight int32) error {
	physHeight := int(physicalSize.Y)
	logWidth := int(logicalWidth)
	logHeight := int(logicalHeight)

	expectedSrcSize := physHeight * int(pitch)
	expectedDstSize := logWidth * logHeight * 4
	if len(src) < expectedSrcSize {
		return fmt.Errorf("source buffer too small: %d < %d", len(src), expectedSrcSize)
	}
	if len(dst) < expectedDstSize {
		return fmt.Errorf("destination buffer too small: %d < %d", len(dst), expectedDstSize)
	}

	pitchInt := int(pitch)
	pitchU32 := pitchInt / 4
	dstU32 := (*[1 << 30]uint32)(unsafe.Pointer(&dst[0]))[:len(dst)/4]
	srcU32 := (*[1 << 30]uint32)(unsafe.Pointer(&src[0]))[:len(src)/4]

	if rotation == disp.ModeRotationRotate90 {
		physHeightMinus1 := physHeight - 1
		
		for logY := 0; logY < logHeight; logY++ {
			dstRowStart := logY * logWidth
			physX := logY
			for logX := 0; logX < logWidth; logX++ {
				physY := physHeightMinus1 - logX
				srcOffset := physY*pitchU32 + physX
				dstOffset := dstRowStart + logX
				dstU32[dstOffset] = srcU32[srcOffset]
			}
		}
	} else if rotation == disp.ModeRotationRotate180 {
		physHeightMinus1 := physHeight - 1
		for logY := 0; logY < logHeight; logY++ {
			dstRowStart := logY * logWidth
			physX := logY
			for logX := 0; logX < logWidth; logX++ {
				physY := physHeightMinus1 - logX
				srcOffset := physY*pitchU32 + physX
				dstOffset := dstRowStart + logX
				dstU32[dstOffset] = srcU32[srcOffset]
			}
		}
	} else if rotation == disp.ModeRotationRotate270 {
		physHeightMinus1 := physHeight - 1
		for logY := 0; logY < logHeight; logY++ {
			dstRowStart := logY * logWidth
			physX := physHeightMinus1 - logY
			for logX := 0; logX < logWidth; logX++ {
				physY := logX
				srcOffset := physY*pitchU32 + physX
				dstOffset := dstRowStart + logX
				dstU32[dstOffset] = srcU32[srcOffset]
			}
		}
	}

	return nil
}

func (sc *ScreenCapture) GetBounds() (disp.Rect, error) {
	if sc.monitorBounds != nil {
		return *sc.monitorBounds, nil
	}

	if sc.outputDuplication == nil {
		return disp.Rect{}, fmt.Errorf("outputDuplication is nil in GetBounds (this should not happen)")
	}

	if sc.dxgiOutput == nil {
		return disp.Rect{}, fmt.Errorf("dxgiOutput is nil")
	}

	desc := disp.OutputDesc{}
	hr := sc.dxgiOutput.GetDesc(&desc)
	if hr := resultcode.ResultCode(hr); hr.Failed() {
		return disp.Rect{}, fmt.Errorf("failed at dxgiOutput.GetDesc. %w", hr)
	}

	return desc.DesktopCoordinates, nil
}

// SetMonitorBounds sets monitor coordinates from MonitorInfo.
// These coordinates will be used in drawCursor instead of GetBounds().
func (sc *ScreenCapture) SetMonitorBounds(left, top, right, bottom int32) {
	sc.monitorBounds = &disp.Rect{
		Left:   left,
		Top:    top,
		Right:  right,
		Bottom: bottom,
	}
}

// SetCaptureCursor enables or disables cursor capture.
func (sc *ScreenCapture) SetCaptureCursor(enabled bool) {
	sc.captureCursor = enabled
}

func newScreenCaptureFormat(device *gfx11.Device, deviceCtx *gfx11.DeviceContext, output uint, format disp.PixelFormat) (*ScreenCapture, error) {
	var hr int32

	var dxgiDevice1 *disp.Device1
	hr = device.QueryInterface(disp.IID_Device1, &dxgiDevice1)
	if hr := resultcode.ResultCode(hr); hr.Failed() {
		return nil, fmt.Errorf("failed at device.QueryInterface. %w", hr)
	}

	var pdxgiAdapter unsafe.Pointer
	hr = dxgiDevice1.GetParent(disp.IID_Adapter1, &pdxgiAdapter)
	if hr := resultcode.ResultCode(hr); hr.Failed() {
		dxgiDevice1.Release()
		return nil, fmt.Errorf("failed at dxgiDevice1.GetAdapter. %w", hr)
	}
	dxgiAdapter := (*disp.Adapter1)(pdxgiAdapter)

	var dxgiOutput *disp.Output
	hr = int32(dxgiAdapter.EnumOutputs(uint32(output), &dxgiOutput))
	if hr := resultcode.ResultCode(hr); hr.Failed() {
		dxgiDevice1.Release()
		return nil, fmt.Errorf("failed at dxgiAdapter.EnumOutputs. %w", hr)
	}

	var dxgiOutput5 *disp.Output5
	hr = dxgiOutput.QueryInterface(disp.IID_Output5, &dxgiOutput5)
	if hr := resultcode.ResultCode(hr); hr.Failed() {
		dxgiOutput.Release()
		dxgiDevice1.Release()
		return nil, fmt.Errorf("failed at dxgiOutput.QueryInterface for Output5. %w", hr)
	}

	var dup *disp.OutputDuplication
	dup = nil

	hr = dxgiOutput5.DuplicateOutput1(dxgiDevice1, 0, []disp.PixelFormat{
		format,
	}, &dup)
	needsSwizzle := false

	hrCode := resultcode.ResultCode(hr)
	if hrCode.Failed() || dup == nil {
		needsSwizzle = true
		var dxgiOutput1 *disp.Output1
		hr := dxgiOutput.QueryInterface(disp.IID_Output1, &dxgiOutput1)
		if hr := resultcode.ResultCode(hr); hr.Failed() {
			dxgiOutput.Release()
			dxgiOutput5.Release()
			dxgiDevice1.Release()
			return nil, fmt.Errorf("failed at dxgiOutput.QueryInterface for Output1. %w", hr)
		}

		dup = nil
		hr = dxgiOutput1.DuplicateOutput(dxgiDevice1, &dup)
		hrCode = resultcode.ResultCode(hr)
		if hrCode.Failed() || dup == nil {
			dxgiOutput1.Release()
			dxgiOutput.Release()
			dxgiOutput5.Release()
			dxgiDevice1.Release()
			if hrCode.Failed() {
				return nil, fmt.Errorf("failed at dxgiOutput1.DuplicateOutput. %w (HRESULT: %v)", hrCode, hr)
			}
			return nil, fmt.Errorf("dxgiOutput1.DuplicateOutput returned nil pointer (HRESULT: %v was successful but pointer is nil)", hr)
		}
		dxgiOutput1.Release()
	}

	if dup == nil {
		dxgiOutput.Release()
		dxgiOutput5.Release()
		dxgiDevice1.Release()
		return nil, fmt.Errorf("DuplicateOutput1/DuplicateOutput returned nil pointer (HRESULT: %v was successful but pointer is nil)", hr)
	}

	sc := &ScreenCapture{
		device:            device,
		deviceCtx:         deviceCtx,
		outputDuplication: dup,
		needsSwizzle:      needsSwizzle,
		dxgiOutput:        dxgiOutput5,
	}

	if sc.outputDuplication == nil {
		if dxgiOutput5 != nil {
			dxgiOutput5.Release()
		}
		if dup != nil {
			dup.Release()
		}
		return nil, fmt.Errorf("outputDuplication became nil after ScreenCapture creation")
	}

	outputDesc := disp.OutputDesc{}
	hr = dxgiOutput5.GetDesc(&outputDesc)
	if hr := resultcode.ResultCode(hr); !hr.Failed() {
		sc.logicalWidth = outputDesc.DesktopCoordinates.Right - outputDesc.DesktopCoordinates.Left
		sc.logicalHeight = outputDesc.DesktopCoordinates.Bottom - outputDesc.DesktopCoordinates.Top
	}

	return sc, nil
}

func NewScreenCapture(device *gfx11.Device, deviceCtx *gfx11.DeviceContext, output uint) (*ScreenCapture, error) {
	return newScreenCaptureFormat(device, deviceCtx, output, disp.PixelFormatB8G8R8A8Unorm)
}
