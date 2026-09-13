package preset_api

import (
	"runtime"
	"syscall"
)

// LogInit starts the process-wide asynchronous frame logger. Frame output must
// still be enabled with SetPrintLog.
func LogInit(name string) int32 {
	value, err := syscall.BytePtrFromString(name)
	if err != nil {
		return PRESET_ERR_INVALID_ARG
	}
	status := invokeStatus("preset_log_init", pointer(value))
	runtime.KeepAlive(value)
	return status
}

func LogShutdown() {
	_, _ = invoke("preset_log_shutdown")
}

func SetPrintLog(enable bool) {
	var value uintptr
	if enable {
		value = 1
	}
	_, _ = invoke("preset_set_print_log", value)
}

func LogDroppedCount() uint64 {
	value, _ := invoke("preset_log_dropped_count")
	return uint64(value)
}

func SetLogFilter(mode uint8, ids []uint32) int32 {
	status := invokeStatus(
		"preset_set_log_filter",
		uintptr(mode),
		slicePointer(ids),
		uintptr(len(ids)),
	)
	runtime.KeepAlive(ids)
	return status
}
