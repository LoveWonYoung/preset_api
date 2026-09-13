package preset_api

import "runtime"

func CanUdsClientNew(device PresetDevice, channel uint8, config *PresetConfig) (client PresetCanUdsClient, status int32) {
	status = invokeStatus(
		"preset_can_uds_client_new",
		uintptr(device), uintptr(channel), pointer(config), pointer(&client),
	)
	runtime.KeepAlive(config)
	return
}

func canUdsRequest(symbol string, client PresetCanUdsClient, payload []byte, timeoutMS uint32, out []byte) (count int, status int32) {
	length := uintptr(0)
	status = invokeStatus(
		symbol,
		uintptr(client), slicePointer(payload), uintptr(len(payload)), uintptr(timeoutMS),
		slicePointer(out), uintptr(len(out)), pointer(&length),
	)
	runtime.KeepAlive(payload)
	runtime.KeepAlive(out)
	return int(length), status
}

func CanUdsRequest(client PresetCanUdsClient, payload []byte, timeoutMS uint32, out []byte) (count int, status int32) {
	return canUdsRequest("preset_can_uds_request", client, payload, timeoutMS, out)
}

func CanUdsFunctionalRequest(client PresetCanUdsClient, payload []byte, timeoutMS uint32, out []byte) (count int, status int32) {
	return canUdsRequest("preset_can_uds_functional_request", client, payload, timeoutMS, out)
}

func CanUdsSetDefaultSTMin(client PresetCanUdsClient, stMinMS uint32) int32 {
	return invokeStatus("preset_can_uds_set_default_st_min", uintptr(client), uintptr(stMinMS))
}

func CanUdsSetDefaultBlockSize(client PresetCanUdsClient, blockSize uint32) int32 {
	return invokeStatus("preset_can_uds_set_default_block_size", uintptr(client), uintptr(blockSize))
}

func CanUdsSetManualFlowControl(client PresetCanUdsClient, enabled bool) int32 {
	var value uintptr
	if enabled {
		value = 1
	}
	return invokeStatus("preset_can_uds_set_manual_flow_control", uintptr(client), value)
}

func CanUdsWrite(client PresetCanUdsClient, id uint32, isFD bool, data []byte) int32 {
	var fd uintptr
	if isFD {
		fd = 1
	}
	status := invokeStatus(
		"preset_can_uds_write",
		uintptr(client), uintptr(id), fd, slicePointer(data), uintptr(len(data)),
	)
	runtime.KeepAlive(data)
	return status
}

func CanUdsTryRead(client PresetCanUdsClient, frames []PresetCanFrame) (count int, status int32) {
	length := uintptr(len(frames))
	status = invokeStatus(
		"preset_can_uds_try_read",
		uintptr(client), slicePointer(frames), pointer(&length),
	)
	runtime.KeepAlive(frames)
	return int(length), status
}

func CanUdsRxGetStats(client PresetCanUdsClient) (stats PresetRxStats, status int32) {
	status = invokeStatus("preset_can_uds_rx_get_stats", uintptr(client), pointer(&stats))
	runtime.KeepAlive(&stats)
	return
}

func CanUdsSetBusLoadEnabled(client PresetCanUdsClient, enabled bool) int32 {
	var value uintptr
	if enabled {
		value = 1
	}
	return invokeStatus("preset_can_uds_set_bus_load_enabled", uintptr(client), value)
}

func CanUdsGetBusLoad(client PresetCanUdsClient) (load PresetBusLoad, status int32) {
	status = invokeStatus("preset_can_uds_get_bus_load", uintptr(client), pointer(&load))
	runtime.KeepAlive(&load)
	return
}

func CanUdsLastError(client PresetCanUdsClient, out []byte) (count int, status int32) {
	length := uintptr(0)
	status = invokeStatus(
		"preset_can_uds_last_error",
		uintptr(client), slicePointer(out), uintptr(len(out)), pointer(&length),
	)
	runtime.KeepAlive(out)
	return int(length), status
}

// CanUdsLastErrorString reads the complete asynchronous error message. An
// empty string means that the client has no recorded error or is not valid.
func CanUdsLastErrorString(client PresetCanUdsClient) string {
	length, status := CanUdsLastError(client, nil)
	if status == PRESET_OK {
		return ""
	}
	if status != PRESET_ERR_BUFFER_TOO_SMALL || length <= 0 {
		return ""
	}
	buffer := make([]byte, length)
	length, status = CanUdsLastError(client, buffer)
	if status != PRESET_OK || length > len(buffer) {
		return ""
	}
	return string(buffer[:length])
}

func LastError(out []byte) (count int, status int32) {
	length := uintptr(0)
	status = invokeStatus(
		"preset_last_error",
		slicePointer(out), uintptr(len(out)), pointer(&length),
	)
	runtime.KeepAlive(out)
	return int(length), status
}

// LastErrorString reads the complete synchronous error message for the current
// OS thread. Because preset_rs stores this error in thread-local storage,
// callers that need a guaranteed association should use runtime.LockOSThread
// around the failing call and this function.
func LastErrorString() string {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	length, status := LastError(nil)
	if status == PRESET_OK {
		return ""
	}
	if status != PRESET_ERR_BUFFER_TOO_SMALL || length <= 0 {
		return ""
	}
	buffer := make([]byte, length)
	length, status = LastError(buffer)
	if status != PRESET_OK || length > len(buffer) {
		return ""
	}
	return string(buffer[:length])
}

func CanUdsClientClose(client *PresetCanUdsClient) int32 {
	status := invokeStatus("preset_can_uds_client_close", pointer(client))
	runtime.KeepAlive(client)
	return status
}
