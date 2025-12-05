package disp

import "go-dda/interop"

type ObjectVtbl struct {
	interop.UnknownVtbl
	SetPrivateData          uintptr
	SetPrivateDataInterface uintptr
	GetPrivateData          uintptr
	GetParent               uintptr
}

type AdapterVtbl struct {
	ObjectVtbl
	EnumOutputs           uintptr
	GetDesc               uintptr
	CheckInterfaceSupport uintptr
}

type Adapter1Vtbl struct {
	AdapterVtbl
	GetDesc1 uintptr
}

type DeviceVtbl struct {
	ObjectVtbl
	CreateSurface          uintptr
	GetAdapter             uintptr
	GetGPUThreadPriority   uintptr
	QueryResourceResidency uintptr
	SetGPUThreadPriority   uintptr
}

type Device1Vtbl struct {
	DeviceVtbl
	GetMaximumFrameLatency uintptr
	SetMaximumFrameLatency uintptr
}

type DeviceSubObjectVtbl struct {
	ObjectVtbl
	GetDevice uintptr
}

type SurfaceVtbl struct {
	DeviceSubObjectVtbl
	GetDesc uintptr
	Map     uintptr
	Unmap   uintptr
}

type ResourceVtbl struct {
	DeviceSubObjectVtbl
	GetSharedHandle     uintptr
	GetUsage            uintptr
	SetEvictionPriority uintptr
	GetEvictionPriority uintptr
}

type OutputVtbl struct {
	ObjectVtbl
	GetDesc                     uintptr
	GetDisplayModeList          uintptr
	FindClosestMatchingMode     uintptr
	WaitForVBlank               uintptr
	TakeOwnership               uintptr
	ReleaseOwnership            uintptr
	GetGammaControlCapabilities uintptr
	SetGammaControl             uintptr
	GetGammaControl             uintptr
	SetDisplaySurface           uintptr
	GetDisplaySurfaceData       uintptr
	GetFrameStatistics          uintptr
}

type Output1Vtbl struct {
	OutputVtbl
	GetDisplayModeList1      uintptr
	FindClosestMatchingMode1 uintptr
	GetDisplaySurfaceData1   uintptr
	DuplicateOutput          uintptr
}

type Output5Vtbl struct {
	OutputVtbl
	GetDisplayModeList1      uintptr
	FindClosestMatchingMode1 uintptr
	GetDisplaySurfaceData1   uintptr
	DuplicateOutput          uintptr
	DuplicateOutput1         uintptr
}

type OutputDuplicationVtbl struct {
	ObjectVtbl
	GetDesc              uintptr
	AcquireNextFrame     uintptr
	GetFrameDirtyRects   uintptr
	GetFrameMoveRects    uintptr
	GetFramePointerShape uintptr
	MapDesktopSurface    uintptr
	UnMapDesktopSurface  uintptr
	ReleaseFrame         uintptr
}

type FactoryVtbl struct {
	ObjectVtbl
	EnumAdapters          uintptr
	MakeWindowAssociation uintptr
	GetWindowAssociation  uintptr
	CreateSwapChain       uintptr
	CreateSoftwareAdapter uintptr
}

type Factory1Vtbl struct {
	FactoryVtbl
	EnumAdapters1 uintptr
	IsCurrent     uintptr
}

