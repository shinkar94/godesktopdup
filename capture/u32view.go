package capture

import "unsafe"

// viewU32 reinterprets a byte slice as uint32 words (len(b)/4 elements).
// unsafe.Slice is portable across architectures: the previous pointer-cast
// idiom (*[1<<30]uint32)(unsafe.Pointer(&b[0])) exceeds the 32-bit address
// space and does not compile on windows/386. Behavior is identical for all
// real buffer sizes.
func viewU32(b []byte) []uint32 {
	return unsafe.Slice((*uint32)(unsafe.Pointer(&b[0])), len(b)/4)
}
