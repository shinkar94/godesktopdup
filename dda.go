package dda

import (
	"fmt"

	"github.com/shinkar94/godesktopdup/capture"
	"github.com/shinkar94/godesktopdup/gfx11"
)

type DesktopDuplication struct {
	device    *gfx11.Device
	deviceCtx *gfx11.DeviceContext
	capture   *capture.ScreenCapture
}

func New(outputIndex uint) (*DesktopDuplication, error) {
	device, deviceCtx, err := gfx11.NewDevice()
	if err != nil {
		return nil, fmt.Errorf("failed to create device: %w", err)
	}

	sc, err := capture.NewScreenCapture(device, deviceCtx, outputIndex)
	if err != nil {
		device.Release()
		deviceCtx.Release()
		return nil, fmt.Errorf("failed to create screen capture: %w", err)
	}

	return &DesktopDuplication{
		device:    device,
		deviceCtx: deviceCtx,
		capture:   sc,
	}, nil
}

func (dd *DesktopDuplication) GetFrameBGRA(buffer []byte, timeoutMs uint) error {
	return dd.capture.GetFrameBGRA(buffer, timeoutMs)
}

func (dd *DesktopDuplication) GetSize() (int, int, error) {
	bounds, err := dd.capture.GetBounds()
	if err != nil {
		return 0, 0, err
	}
	return int(bounds.Right - bounds.Left), int(bounds.Bottom - bounds.Top), nil
}

func (dd *DesktopDuplication) GetBounds() (int, int, int, int, error) {
	bounds, err := dd.capture.GetBounds()
	if err != nil {
		return 0, 0, 0, 0, err
	}
	return int(bounds.Left), int(bounds.Top), int(bounds.Right), int(bounds.Bottom), nil
}

func (dd *DesktopDuplication) SetMonitorBounds(left, top, right, bottom int) {
	dd.capture.SetMonitorBounds(int32(left), int32(top), int32(right), int32(bottom))
}

// SetCaptureCursor enables or disables cursor capture.
func (dd *DesktopDuplication) SetCaptureCursor(enabled bool) {
	dd.capture.SetCaptureCursor(enabled)
}

func (dd *DesktopDuplication) Release() {
	if dd.capture != nil {
		dd.capture.Release()
	}
	if dd.deviceCtx != nil {
		dd.deviceCtx.Release()
	}
	if dd.device != nil {
		dd.device.Release()
	}
}

