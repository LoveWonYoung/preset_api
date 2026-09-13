//go:build windows && amd64

package preset_api

import "runtime"

func CanUdsClientNew(device PresetDevice, channel uint8, config *PresetConfig) (client PresetCanUdsClient, status int32) {
	var handle uintptr
	status = withHandle(device.state, func(deviceHandle uintptr) int32 {
		return invokeStatus(
			"preset_can_uds_client_new",
			word(deviceHandle), word(uintptr(channel)), pointer(config), pointer(&handle),
		)
	})
	client.state = newHandle(handle)
	return
}

func canUdsRequest(symbol string, client PresetCanUdsClient, payload []byte, timeoutMS uint32, out []byte) (count int, status int32) {
	length := uintptr(0)
	status = withHandle(client.state, func(clientHandle uintptr) int32 {
		return invokeStatus(
			symbol,
			word(clientHandle), slicePointer(payload), word(uintptr(len(payload))), word(uintptr(timeoutMS)),
			slicePointer(out), word(uintptr(len(out))), pointer(&length),
		)
	})
	return int(length), status
}

func CanUdsRequest(client PresetCanUdsClient, payload []byte, timeoutMS uint32, out []byte) (count int, status int32) {
	return canUdsRequest("preset_can_uds_request", client, payload, timeoutMS, out)
}

func CanUdsFunctionalRequest(client PresetCanUdsClient, payload []byte, timeoutMS uint32, out []byte) (count int, status int32) {
	return canUdsRequest("preset_can_uds_functional_request", client, payload, timeoutMS, out)
}

func CanUdsSetDefaultSTMin(client PresetCanUdsClient, stMinMS uint32) int32 {
	return withHandle(client.state, func(handle uintptr) int32 {
		return invokeStatus("preset_can_uds_set_default_st_min", word(handle), word(uintptr(stMinMS)))
	})
}

func CanUdsSetDefaultBlockSize(client PresetCanUdsClient, blockSize uint32) int32 {
	return withHandle(client.state, func(handle uintptr) int32 {
		return invokeStatus("preset_can_uds_set_default_block_size", word(handle), word(uintptr(blockSize)))
	})
}

func CanUdsSetManualFlowControl(client PresetCanUdsClient, enabled bool) int32 {
	var value uintptr
	if enabled {
		value = 1
	}
	return withHandle(client.state, func(handle uintptr) int32 {
		return invokeStatus("preset_can_uds_set_manual_flow_control", word(handle), word(value))
	})
}

func CanUdsSetBRS(client PresetCanUdsClient, enabled bool) int32 {
	var value uintptr
	if enabled {
		value = 1
	}
	return withHandle(client.state, func(handle uintptr) int32 {
		return invokeStatus("preset_can_uds_set_brs", word(handle), word(value))
	})
}

func CanUdsGetBRS(client PresetCanUdsClient) (enabled bool, status int32) {
	var value uint8
	status = withHandle(client.state, func(handle uintptr) int32 {
		return invokeStatus("preset_can_uds_get_brs", word(handle), pointer(&value))
	})
	return value != 0, status
}

func CanUdsSetRawTxEcho(client PresetCanUdsClient, enabled bool) int32 {
	var value uintptr
	if enabled {
		value = 1
	}
	return withHandle(client.state, func(handle uintptr) int32 {
		return invokeStatus("preset_can_uds_set_raw_tx_echo", word(handle), word(value))
	})
}

func CanUdsWrite(client PresetCanUdsClient, id uint32, isFD bool, data []byte) int32 {
	var fd uintptr
	if isFD {
		fd = 1
	}
	return withHandle(client.state, func(handle uintptr) int32 {
		return invokeStatus(
			"preset_can_uds_write",
			word(handle), word(uintptr(id)), word(fd), slicePointer(data), word(uintptr(len(data))),
		)
	})
}

func CanUdsTryRead(client PresetCanUdsClient, frames []PresetCanFrame) (count int, status int32) {
	length := uintptr(len(frames))
	status = withHandle(client.state, func(handle uintptr) int32 {
		return invokeStatus(
			"preset_can_uds_try_read",
			word(handle), slicePointer(frames), pointer(&length),
		)
	})
	return int(length), status
}

func CanUdsTryReadEx(client PresetCanUdsClient, frames []PresetCanFrameEx) (count int, status int32) {
	length := uintptr(len(frames))
	status = withHandle(client.state, func(handle uintptr) int32 {
		return invokeStatus(
			"preset_can_uds_try_read_ex",
			word(handle), slicePointer(frames), pointer(&length),
		)
	})
	return int(length), status
}

func CanUdsRxGetStats(client PresetCanUdsClient) (stats PresetRxStats, status int32) {
	status = withHandle(client.state, func(handle uintptr) int32 {
		return invokeStatus("preset_can_uds_rx_get_stats", word(handle), pointer(&stats))
	})
	return
}

func CanUdsSetBusLoadEnabled(client PresetCanUdsClient, enabled bool) int32 {
	var value uintptr
	if enabled {
		value = 1
	}
	return withHandle(client.state, func(handle uintptr) int32 {
		return invokeStatus("preset_can_uds_set_bus_load_enabled", word(handle), word(value))
	})
}

func CanUdsGetBusLoad(client PresetCanUdsClient) (load PresetBusLoad, status int32) {
	status = withHandle(client.state, func(handle uintptr) int32 {
		return invokeStatus("preset_can_uds_get_bus_load", word(handle), pointer(&load))
	})
	return
}

func CanUdsLastError(client PresetCanUdsClient, out []byte) (count int, status int32) {
	return canUdsLastError(client.state, out)
}

func canUdsLastError(state *handleState, out []byte) (count int, status int32) {
	length := uintptr(0)
	status = withHandle(state, func(handle uintptr) int32 {
		return invokeStatus(
			"preset_can_uds_last_error",
			word(handle), slicePointer(out), word(uintptr(len(out))), pointer(&length),
		)
	})
	return int(length), status
}

// CanUdsLastErrorString reads the complete asynchronous error message. An
// empty string means that the client has no recorded error or is not valid.
func CanUdsLastErrorString(client PresetCanUdsClient) string {
	var message string
	withHandle(client.state, func(handle uintptr) int32 {
		message = canUdsLastErrorStringLocked(handle)
		return PRESET_OK
	})
	return message
}

func canUdsLastErrorStringLocked(handle uintptr) string {
	length := uintptr(0)
	status := invokeStatus(
		"preset_can_uds_last_error",
		word(handle), callArgument{}, word(0), pointer(&length),
	)
	if status == PRESET_OK || status != PRESET_ERR_BUFFER_TOO_SMALL || length == 0 {
		return ""
	}
	buffer := make([]byte, int(length))
	status = invokeStatus(
		"preset_can_uds_last_error",
		word(handle), slicePointer(buffer), word(uintptr(len(buffer))), pointer(&length),
	)
	if status != PRESET_OK || length > uintptr(len(buffer)) {
		return ""
	}
	return string(buffer[:int(length)])
}

func LastError(out []byte) (count int, status int32) {
	length := uintptr(0)
	status = invokeStatus(
		"preset_last_error",
		slicePointer(out), word(uintptr(len(out))), pointer(&length),
	)
	return int(length), status
}

// WithLastError executes call and captures preset_rs's thread-local error before
// Go can move the goroutine to another OS thread. Use it when the numeric status
// alone is insufficient.
func WithLastError(call func() int32) (status int32, message string) {
	if call == nil {
		return PRESET_ERR_INVALID_ARG, "call is nil"
	}
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	status = call()
	if status != PRESET_OK {
		if status == PRESET_ERR_TRANSPORT {
			if err := DLLLoadError(); err != nil {
				message = err.Error()
			}
		}
		if message == "" {
			message = lastErrorStringLocked()
		}
	}
	return
}

// LastErrorString reads the current OS thread's complete synchronous error.
// Prefer WithLastError when the message must be associated with one call.
func LastErrorString() string {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return lastErrorStringLocked()
}

func lastErrorStringLocked() string {
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
	if client == nil {
		return PRESET_ERR_NULL_PTR
	}
	return closeHandle(client.state, "preset_can_uds_client_close")
}
