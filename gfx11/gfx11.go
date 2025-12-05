package gfx11

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/shinkar94/godesktopdup/disp"
	resultcode "github.com/shinkar94/godesktopdup/errors"
	"github.com/shinkar94/godesktopdup/interop"
	"golang.org/x/sys/windows"
)

var (
	modD3D11              = windows.NewLazySystemDLL("d3d11.dll")
	procD3D11CreateDevice = modD3D11.NewProc("D3D11CreateDevice")

	IID_Texture2D, _ = windows.GUIDFromString("{6f15aaf2-d208-4e89-9ab4-489535d34f9c}")
	IID_Debug, _     = windows.GUIDFromString("{79cf2233-7536-4948-9d36-1e4692dc5760}")
	IID_InfoQueue, _ = windows.GUIDFromString("{6543dbb6-1b48-42f5-ab82-e97ec74326f6}")
)

func createDevice(ppDevice **Device, ppDeviceContext **DeviceContext) error {
	var factory1 *disp.Factory1
	if err := disp.CreateDXGIFactory1(&factory1); err != nil {
		return fmt.Errorf("CreateDXGIFactory1: %w", err)
	}
	defer factory1.Release()

	var adapter1 *disp.Adapter1
	var desc disp.AdapterDesc1

	ai := uint32(0)
	for {
		hr := factory1.EnumAdapters1(ai, &adapter1)
		if resultcode.ResultCode(hr).Failed() {
			break
		}
		ai++

		hr = int32(adapter1.GetDesc1(&desc))
		if resultcode.ResultCode(hr).Failed() {
			adapter1.Release()
			adapter1 = nil
			continue
		}

		if (desc.Flags & disp.AdapterFlagSoftware) == 0 {
			break
		}
		adapter1.Release()
		adapter1 = nil
	}

	if adapter1 == nil {
		hr := factory1.EnumAdapters1(0, &adapter1)
		if resultcode.ResultCode(hr).Failed() {
			return fmt.Errorf("failed to fallback to default display adapter")
		}
	}

	if adapter1 == nil {
		return fmt.Errorf("no suitable adapter found")
	}

	defer adapter1.Release()

	fflags := [...]uint32{
		0xc100, // D3D_FEATURE_LEVEL_12_1
		0xc000, // D3D_FEATURE_LEVEL_12_0
		0xb100, // D3D_FEATURE_LEVEL_11_1
		0xb000, // D3D_FEATURE_LEVEL_11_0
		0xa100, // D3D_FEATURE_LEVEL_10_1
		0xa000, // D3D_FEATURE_LEVEL_10_0
	}
	featureLevel := 0x9100
	flags := 0

	const DriverTypeUnknown = 0
	ret, _, _ := syscall.SyscallN(
		procD3D11CreateDevice.Addr(),
		uintptr(unsafe.Pointer(adapter1)),
		uintptr(DriverTypeUnknown),
		uintptr(0),
		uintptr(flags),
		uintptr(unsafe.Pointer(&fflags[0])),
		uintptr(len(fflags)),
		uintptr(SDKVersion),
		uintptr(unsafe.Pointer(ppDevice)),
		uintptr(unsafe.Pointer(&featureLevel)),
		uintptr(unsafe.Pointer(ppDeviceContext)),
	)

	if ret != 0 {
		return resultcode.ResultCode(ret)
	}
	return nil
}

func NewDevice() (*Device, *DeviceContext, error) {
	var device *Device
	var deviceCtx *DeviceContext

	err := createDevice(&device, &deviceCtx)

	if err != nil || device == nil || deviceCtx == nil {
		return nil, nil, err
	}

	return device, deviceCtx, nil
}

type Texture2D struct {
	vtbl *Texture2DVtbl
}

func (obj *Texture2D) GetDesc(desc *Texture2DDesc) int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.GetDesc,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(desc)),
	)
	return int32(ret)
}

func (obj *Texture2D) Release() int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	)
	return int32(ret)
}

func (obj *Texture2D) QueryInterface(iid windows.GUID, pp interface{}) int32 {
	return interop.QueryInterface(obj, obj.vtbl.QueryInterface, &iid, pp)
}

type Device struct {
	vtbl *DeviceVtbl
}

func (obj *Device) QueryInterface(iid windows.GUID, pp interface{}) int32 {
	return interop.QueryInterface(obj, obj.vtbl.QueryInterface, &iid, pp)
}

func (obj *Device) CreateTexture2D(desc *Texture2DDesc, ppTexture2D **Texture2D) int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.CreateTexture2D,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(desc)),
		uintptr(0),
		uintptr(unsafe.Pointer(ppTexture2D)),
	)
	return int32(ret)
}

func (obj *Device) Release() int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	)
	return int32(ret)
}

type DeviceContext struct {
	vtbl *DeviceContextVtbl
}

func (obj *DeviceContext) CopyResource2D(dst, src *Texture2D) int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.CopyResource,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(dst)),
		uintptr(unsafe.Pointer(src)),
	)
	return int32(ret)
}

func (obj *DeviceContext) CopySubresourceRegion2D(dst *Texture2D, dstSubResource, dstX, dstY, dstZ uint32, src *Texture2D, srcSubResource uint32, pSrcBox *Box) int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.CopySubresourceRegion,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(dst)),
		uintptr(dstSubResource),
		uintptr(dstX),
		uintptr(dstY),
		uintptr(dstZ),
		uintptr(unsafe.Pointer(src)),
		uintptr(srcSubResource),
		uintptr(unsafe.Pointer(pSrcBox)),
	)
	return int32(ret)
}

func (obj *DeviceContext) Release() int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	)
	return int32(ret)
}

