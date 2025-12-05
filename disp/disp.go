package disp

import (
	"syscall"
	"unsafe"

	"github.com/shinkar94/godesktopdup/interop"
	resultcode "github.com/shinkar94/godesktopdup/errors"
	"golang.org/x/sys/windows"
)

var (
	modDXGI                = windows.NewLazySystemDLL("dxgi.dll")
	procCreateDXGIFactory1 = modDXGI.NewProc("CreateDXGIFactory1")

	IID_Device1, _   = windows.GUIDFromString("{77db970f-6276-48ba-ba28-070143b4392c}")
	IID_Adapter1, _  = windows.GUIDFromString("{29038f61-3839-4626-91fd-086879011a05}")
	IID_Output1, _   = windows.GUIDFromString("{00cddea8-939b-4b83-a340-a685226666cc}")
	IID_Output5, _   = windows.GUIDFromString("{80A07424-AB52-42EB-833C-0C42FD282D98}")
	IID_Factory1, _  = windows.GUIDFromString("{770aae78-f26f-4dba-a829-253c83d1b387}")
	IID_Surface, _  = windows.GUIDFromString("{cafcb56c-6ac3-4889-bf47-9e23bbd260ec}")
)

const (
	MapRead    = 1 << 0
	MapWrite   = 1 << 1
	MapDiscard = 1 << 2
)

type Factory1 struct {
	vtbl *Factory1Vtbl
}

func (obj *Factory1) Release() int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	)
	return int32(ret)
}

func (obj *Factory1) EnumAdapters1(adapter uint32, pp **Adapter1) int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.EnumAdapters1,
		uintptr(unsafe.Pointer(obj)),
		uintptr(adapter),
		uintptr(unsafe.Pointer(pp)),
	)
	return int32(ret)
}

func CreateDXGIFactory1(ppFactory **Factory1) error {
	ret, _, _ := syscall.SyscallN(
		procCreateDXGIFactory1.Addr(),
		uintptr(unsafe.Pointer(&IID_Factory1)),
		uintptr(unsafe.Pointer(ppFactory)),
	)
	if ret != 0 {
		return resultcode.ResultCode(ret)
	}
	return nil
}

type Adapter1 struct {
	vtbl *Adapter1Vtbl
}

func (obj *Adapter1) Release() int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	)
	return int32(ret)
}

func (obj *Adapter1) EnumOutputs(output uint32, pp **Output) uint32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.EnumOutputs,
		uintptr(unsafe.Pointer(obj)),
		uintptr(output),
		uintptr(unsafe.Pointer(pp)),
	)
	return uint32(ret)
}

func (obj *Adapter1) GetDesc1(p *AdapterDesc1) uint32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.GetDesc1,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(p)),
	)
	return uint32(ret)
}

type Adapter struct {
	vtbl *AdapterVtbl
}

func (obj *Adapter) EnumOutputs(output uint32, pp **Output) uint32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.EnumOutputs,
		uintptr(unsafe.Pointer(obj)),
		uintptr(output),
		uintptr(unsafe.Pointer(pp)),
	)
	return uint32(ret)
}

func (obj *Adapter) Release() int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	)
	return int32(ret)
}

type Device1 struct {
	vtbl *Device1Vtbl
}

func (obj *Device1) QueryInterface(iid windows.GUID, pp interface{}) int32 {
	return interop.QueryInterface(obj, obj.vtbl.QueryInterface, &iid, pp)
}

func (obj *Device1) GetParent(iid windows.GUID, pp *unsafe.Pointer) int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.GetParent,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&iid)),
		uintptr(unsafe.Pointer(pp)),
	)
	return int32(ret)
}

func (obj *Device1) Release() int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	)
	return int32(ret)
}

type Output struct {
	vtbl *OutputVtbl
}

func (obj *Output) QueryInterface(iid windows.GUID, pp interface{}) int32 {
	return interop.QueryInterface(obj, obj.vtbl.QueryInterface, &iid, pp)
}

func (obj *Output) GetParent(iid windows.GUID, pp *unsafe.Pointer) int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.GetParent,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&iid)),
		uintptr(unsafe.Pointer(pp)),
	)
	return int32(ret)
}

func (obj *Output) Release() int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	)
	return int32(ret)
}

type Output1 struct {
	vtbl *Output1Vtbl
}

func (obj *Output1) DuplicateOutput(device1 *Device1, ppOutputDuplication **OutputDuplication) int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.DuplicateOutput,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(device1)),
		uintptr(unsafe.Pointer(ppOutputDuplication)),
	)
	return int32(ret)
}

func (obj *Output1) GetParent(iid windows.GUID, pp *unsafe.Pointer) int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.GetParent,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&iid)),
		uintptr(unsafe.Pointer(pp)),
	)
	return int32(ret)
}

func (obj *Output1) Release() int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	)
	return int32(ret)
}

type Output5 struct {
	vtbl *Output5Vtbl
}

func (obj *Output5) GetDesc(desc *OutputDesc) int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.GetDesc,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(desc)),
	)
	return int32(ret)
}

func (obj *Output5) DuplicateOutput1(device1 *Device1, flags uint32, pSupportedFormats []PixelFormat, ppOutputDuplication **OutputDuplication) int32 {
	pFormats := &pSupportedFormats[0]
	var ppDup unsafe.Pointer
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.DuplicateOutput1,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(device1)),
		uintptr(flags),
		uintptr(len(pSupportedFormats)),
		uintptr(unsafe.Pointer(pFormats)),
		uintptr(unsafe.Pointer(&ppDup)),
	)
	if ret == 0 && ppDup != nil {
		*ppOutputDuplication = (*OutputDuplication)(ppDup)
	}
	return int32(ret)
}

func (obj *Output5) GetParent(iid windows.GUID, pp *unsafe.Pointer) int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.GetParent,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&iid)),
		uintptr(unsafe.Pointer(pp)),
	)
	return int32(ret)
}

func (obj *Output5) Release() int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	)
	return int32(ret)
}

type Resource struct {
	vtbl *ResourceVtbl
}

func (obj *Resource) QueryInterface(iid windows.GUID, pp interface{}) int32 {
	return interop.QueryInterface(obj, obj.vtbl.QueryInterface, &iid, pp)
}

func (obj *Resource) Release() int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	)
	return int32(ret)
}

type Surface struct {
	vtbl *SurfaceVtbl
}

func (obj *Surface) QueryInterface(iid windows.GUID, pp interface{}) int32 {
	return interop.QueryInterface(obj, obj.vtbl.QueryInterface, &iid, pp)
}

func (obj *Surface) Map(pLockedRect *MappedRect, mapFlags uint32) int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.Map,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(pLockedRect)),
		uintptr(mapFlags),
	)
	return int32(ret)
}

func (obj *Surface) Unmap() int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.Unmap,
		uintptr(unsafe.Pointer(obj)),
	)
	return int32(ret)
}

func (obj *Surface) Release() int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	)
	return int32(ret)
}

type OutputDuplication struct {
	vtbl *OutputDuplicationVtbl
}

func (obj *OutputDuplication) GetFrameMoveRects(buffer []DuplicationMoveRect, rectsRequired *uint32) int32 {
	var buf *DuplicationMoveRect
	if len(buffer) > 0 {
		buf = &buffer[0]
	}
	size := uint32(len(buffer) * 24)
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.GetFrameMoveRects,
		uintptr(unsafe.Pointer(obj)),
		uintptr(size),
		uintptr(unsafe.Pointer(buf)),
		uintptr(unsafe.Pointer(rectsRequired)),
	)
	*rectsRequired = *rectsRequired / 24
	return int32(ret)
}

func (obj *OutputDuplication) GetFrameDirtyRects(buffer []Rect, rectsRequired *uint32) int32 {
	var buf *Rect
	if len(buffer) > 0 {
		buf = &buffer[0]
	}
	size := uint32(len(buffer) * 16)
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.GetFrameDirtyRects,
		uintptr(unsafe.Pointer(obj)),
		uintptr(size),
		uintptr(unsafe.Pointer(buf)),
		uintptr(unsafe.Pointer(rectsRequired)),
	)
	*rectsRequired = *rectsRequired / 16
	return int32(ret)
}

func (obj *OutputDuplication) GetFramePointerShape(pointerShapeBufferSize uint32,
	pPointerShapeBuffer []byte,
	pPointerShapeBufferSizeRequired *uint32,
	pPointerShapeInfo *DuplicationPointerShapeInfo) int32 {

	var buf *byte
	if len(pPointerShapeBuffer) > 0 {
		buf = &pPointerShapeBuffer[0]
	}

	ret, _, _ := syscall.SyscallN(
		obj.vtbl.GetFramePointerShape,
		uintptr(unsafe.Pointer(obj)),
		uintptr(pointerShapeBufferSize),
		uintptr(unsafe.Pointer(buf)),
		uintptr(unsafe.Pointer(pPointerShapeBufferSizeRequired)),
		uintptr(unsafe.Pointer(pPointerShapeInfo)),
	)

	return int32(ret)
}

func (obj *OutputDuplication) GetDesc(desc *DuplicationDesc) int32 {
	if obj == nil || obj.vtbl == nil {
		return -2147024809 // E_INVALIDARG
	}
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.GetDesc,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(desc)),
	)
	return int32(ret)
}

func (obj *OutputDuplication) MapDesktopSurface(pLockedRect *MappedRect) int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.MapDesktopSurface,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(pLockedRect)),
	)
	return int32(ret)
}

func (obj *OutputDuplication) UnMapDesktopSurface() int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.UnMapDesktopSurface,
		uintptr(unsafe.Pointer(obj)),
	)
	return int32(ret)
}

func (obj *OutputDuplication) AddRef() uint32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.AddRef,
		uintptr(unsafe.Pointer(obj)),
	)
	return uint32(ret)
}

func (obj *OutputDuplication) Release() uint32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	)
	return uint32(ret)
}

func (obj *OutputDuplication) AcquireNextFrame(timeoutMs uint32, pFrameInfo *DuplicationFrameInfo, ppDesktopResource **Resource) uint32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.AcquireNextFrame,
		uintptr(unsafe.Pointer(obj)),
		uintptr(timeoutMs),
		uintptr(unsafe.Pointer(pFrameInfo)),
		uintptr(unsafe.Pointer(ppDesktopResource)),
	)
	return uint32(ret)
}

func (obj *OutputDuplication) ReleaseFrame() uint32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.ReleaseFrame,
		uintptr(unsafe.Pointer(obj)),
	)
	return uint32(ret)
}

