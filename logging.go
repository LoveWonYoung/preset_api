//go:build windows && amd64

package preset_api

import (
	"syscall"
)

// LogInit starts the process-wide asynchronous frame logger. Frame output must
// still be enabled with SetPrintLog.
func LogInit(name string) int32 {
	value, err := syscall.BytePtrFromString(name)
	if err != nil {
		return PRESET_ERR_INVALID_ARG
	}
	return invokeStatus("preset_log_init", pointer(value))
}

func LogShutdown() error {
	_, err := invoke("preset_log_shutdown")
	return err
}

func SetPrintLog(enable bool) error {
	var value uintptr
	if enable {
		value = 1
	}
	_, err := invoke("preset_set_print_log", word(value))
	return err
}

func LogDroppedCount() (uint64, error) {
	result, err := invoke("preset_log_dropped_count")
	return uint64(result.r1), err
}

func SetLogFilter(mode uint8, ids []uint32) int32 {
	status := invokeStatus(
		"preset_set_log_filter",
		word(uintptr(mode)),
		slicePointer(ids),
		word(uintptr(len(ids))),
	)
	return status
}
