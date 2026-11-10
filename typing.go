package main

import "unsafe"

// 数字类型约束
type Number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

func NumberPtr[T Number](array []T) uintptr {
	return uintptr(unsafe.Pointer(&array[0]))
}

func RawPtr[T any](value T) uintptr {
	return uintptr(unsafe.Pointer(&value))
}
