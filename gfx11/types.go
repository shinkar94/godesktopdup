package gfx11

import "go-dda/disp"

type Box struct {
	Left, Top, Front, Right, Bottom, Back uint32
}

type Texture2DDesc struct {
	Width          uint32
	Height         uint32
	MipLevels      uint32
	ArraySize      uint32
	Format         uint32
	SampleDesc     disp.SampleDesc
	Usage          Usage
	BindFlags      uint32
	CPUAccessFlags uint32
	MiscFlags      uint32
}

type Usage uint32

const (
	UsageDefault   Usage = 0
	UsageImmutable Usage = 1
	UsageDynamic   Usage = 2
	UsageStaging   Usage = 3
)

const (
	CPUAccessRead = 0x20000
	SDKVersion    = 7
)

