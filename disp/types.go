package disp

import "unsafe"

type UInt32 = uint32
type SizeT = uintptr
type ULong = uint32
type Long = int32

type Rational struct {
	Numerator   uint32
	Denominator uint32
}

type ModeRotation uint32

const (
	ModeRotationIdentity   ModeRotation = 0 // No rotation
	ModeRotationRotate90   ModeRotation = 1 // 90° clockwise
	ModeRotationRotate180  ModeRotation = 2 // 180°
	ModeRotationRotate270  ModeRotation = 3 // 270° clockwise (90° counter-clockwise)
)

type OutputDesc struct {
	DeviceName         [32]uint16
	DesktopCoordinates Rect
	AttachedToDesktop  uint32
	Rotation           ModeRotation
	Monitor            uintptr
}

type ModeDesc struct {
	Width            uint32
	Height           uint32
	Rational         Rational
	Format           uint32
	ScanlineOrdering uint32
	Scaling          uint32
}

type DuplicationDesc struct {
	ModeDesc                   ModeDesc
	Rotation                   uint32
	DesktopImageInSystemMemory uint32
}

type SampleDesc struct {
	Count   uint32
	Quality uint32
}

type Point struct {
	X int32
	Y int32
}

type Rect struct {
	Left, Top, Right, Bottom int32
}

type DuplicationMoveRect struct {
	Src  Point
	Dest Rect
}

type DuplicationPointerPosition struct {
	Position Point
	Visible  uint32
}

type DuplicationFrameInfo struct {
	LastPresentTime           int64
	LastMouseUpdateTime       int64
	AccumulatedFrames         uint32
	RectsCoalesced            uint32
	ProtectedContentMaskedOut uint32
	PointerPosition           DuplicationPointerPosition
	TotalMetadataBufferSize   uint32
	PointerShapeBufferSize    uint32
}

type MappedRect struct {
	Pitch int32
	PBits unsafe.Pointer
}

type DuplicationPointerShapeType uint32

const (
	DuplicationPointerShapeTypeMonochrome   DuplicationPointerShapeType = 1
	DuplicationPointerShapeTypeColor        DuplicationPointerShapeType = 2
	DuplicationPointerShapeTypeMaskedColor  DuplicationPointerShapeType = 4
)

type DuplicationPointerShapeInfo struct {
	Type    DuplicationPointerShapeType
	Width   uint32
	Height  uint32
	Pitch   uint32
	HotSpot Point
}

type Luid struct {
	LowPart  ULong
	HighPart Long
}

type AdapterFlag uint32

const (
	AdapterFlagNone     AdapterFlag = 0
	AdapterFlagRemote   AdapterFlag = 1
	AdapterFlagSoftware AdapterFlag = 2
)

type AdapterDesc1 struct {
	Description           [128]uint16
	VendorId              UInt32
	DeviceId              UInt32
	SubSysId              UInt32
	Revision              UInt32
	DedicatedVideoMemory  SizeT
	DedicatedSystemMemory SizeT
	SharedSystemMemory    SizeT
	AdapterLuid           Luid
	Flags                 AdapterFlag
}

type PixelFormat uint32

const (
	PixelFormatUnknown                PixelFormat = 0
	PixelFormatB8G8R8A8Unorm          PixelFormat = 87
	PixelFormatR8G8B8A8Unorm         PixelFormat = 28
)

