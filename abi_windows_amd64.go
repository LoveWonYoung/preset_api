//go:build windows && amd64

package preset_api

import "unsafe"

// Microsoft x64 returns plain C structs of 1, 2, 4, or 8 bytes in RAX. All
// other struct sizes used by this API are returned through a hidden first
// argument pointing at caller-owned storage.
func callStruct(name string, destination unsafe.Pointer, size uintptr) error {
	switch size {
	case 1, 2, 4, 8:
		result, err := invoke(name)
		if err != nil {
			return err
		}
		copyRegisterResult(destination, size, result.r1)
		return nil
	default:
		_, err := invoke(name, rawPointer(destination))
		return err
	}
}

func copyRegisterResult(destination unsafe.Pointer, size, value uintptr) {
	switch size {
	case 1:
		*(*uint8)(destination) = uint8(value)
	case 2:
		*(*uint16)(destination) = uint16(value)
	case 4:
		*(*uint32)(destination) = uint32(value)
	case 8:
		*(*uint64)(destination) = uint64(value)
	}
}
