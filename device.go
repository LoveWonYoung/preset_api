//go:build windows && amd64

package preset_api

func AutoOpen(config *PresetAutoConfig) (device PresetDevice, status int32) {
	var handle uintptr
	status = invokeStatus("preset_auto_open", pointer(config), pointer(&handle))
	device.state = newHandle(handle)
	return
}

func DeviceGetBackend(device PresetDevice) (backend uint8, status int32) {
	status = withHandle(device.state, func(handle uintptr) int32 {
		return invokeStatus("preset_device_get_backend", word(handle), pointer(&backend))
	})
	return
}

func ToomossOpen() (device PresetDevice, status int32) {
	var handle uintptr
	status = invokeStatus("preset_toomoss_open", pointer(&handle))
	device.state = newHandle(handle)
	return
}

func ToomossCanInit(device PresetDevice, config *PresetToomossConfig) int32 {
	return withHandle(device.state, func(handle uintptr) int32 {
		return invokeStatus("preset_toomoss_can_init", word(handle), pointer(config))
	})
}

func ToomossCanInitWithTiming(device PresetDevice, config *PresetToomossConfig, timing *PresetCanFdTiming) int32 {
	return withHandle(device.state, func(handle uintptr) int32 {
		return invokeStatus("preset_toomoss_can_init_with_timing", word(handle), pointer(config), pointer(timing))
	})
}

func PcanOpen(config *PresetPCANConfig) (device PresetDevice, status int32) {
	var handle uintptr
	status = invokeStatus("preset_pcan_open", pointer(config), pointer(&handle))
	device.state = newHandle(handle)
	return
}

func PcanOpenWithTiming(config *PresetPCANConfig, timing *PresetCanFdTiming) (device PresetDevice, status int32) {
	var handle uintptr
	status = invokeStatus("preset_pcan_open_with_timing", pointer(config), pointer(timing), pointer(&handle))
	device.state = newHandle(handle)
	return
}

func PcanCanInit(device PresetDevice, config *PresetPCANConfig) int32 {
	return withHandle(device.state, func(handle uintptr) int32 {
		return invokeStatus("preset_pcan_can_init", word(handle), pointer(config))
	})
}

func PcanCanInitWithTiming(device PresetDevice, config *PresetPCANConfig, timing *PresetCanFdTiming) int32 {
	return withHandle(device.state, func(handle uintptr) int32 {
		return invokeStatus("preset_pcan_can_init_with_timing", word(handle), pointer(config), pointer(timing))
	})
}

func TsmasterOpen(config *PresetTSMasterConfig) (device PresetDevice, status int32) {
	var handle uintptr
	status = invokeStatus("preset_tsmaster_open", pointer(config), pointer(&handle))
	device.state = newHandle(handle)
	return
}

func TsmasterCanInit(device PresetDevice, config *PresetTSMasterConfig) int32 {
	return withHandle(device.state, func(handle uintptr) int32 {
		return invokeStatus("preset_tsmaster_can_init", word(handle), pointer(config))
	})
}

func VectorOpen(config *PresetVectorConfig) (device PresetDevice, status int32) {
	var handle uintptr
	status = invokeStatus("preset_vector_open", pointer(config), pointer(&handle))
	device.state = newHandle(handle)
	return
}

func VectorCanInit(device PresetDevice, config *PresetVectorConfig) int32 {
	return withHandle(device.state, func(handle uintptr) int32 {
		return invokeStatus("preset_vector_can_init", word(handle), pointer(config))
	})
}

// DeviceCanWrite sends a raw CAN/CAN-FD frame directly on an initialized
// channel without creating a UDS client.
func DeviceCanWrite(device PresetDevice, channel uint8, id uint32, isFD bool, data []byte) int32 {
	var fd uintptr
	if isFD {
		fd = 1
	}
	return withHandle(device.state, func(handle uintptr) int32 {
		return invokeStatus(
			"preset_can_write",
			word(handle), word(uintptr(channel)), word(uintptr(id)), word(fd),
			slicePointer(data), word(uintptr(len(data))),
		)
	})
}

// DeviceCanTryRead reads raw received frames directly from an initialized
// channel. It returns zero frames when the backend is not due for another poll.
func DeviceCanTryRead(device PresetDevice, channel uint8, frames []PresetCanFrame) (count int, status int32) {
	length := uintptr(len(frames))
	status = withHandle(device.state, func(handle uintptr) int32 {
		return invokeStatus(
			"preset_can_try_read",
			word(handle), word(uintptr(channel)), slicePointer(frames), pointer(&length),
		)
	})
	return int(length), status
}

// DeviceCanTryReadEx is DeviceCanTryRead with DLC, BRS, direction, and
// hardware timestamp metadata.
func DeviceCanTryReadEx(device PresetDevice, channel uint8, frames []PresetCanFrameEx) (count int, status int32) {
	length := uintptr(len(frames))
	status = withHandle(device.state, func(handle uintptr) int32 {
		return invokeStatus(
			"preset_can_try_read_ex",
			word(handle), word(uintptr(channel)), slicePointer(frames), pointer(&length),
		)
	})
	return int(length), status
}

func CanTpWriteSingleFrame(device PresetDevice, channel uint8, id uint32, config *PresetTpFrameConfig, data []byte) int32 {
	return withHandle(device.state, func(handle uintptr) int32 {
		return invokeStatus(
			"preset_can_tp_write_single_frame",
			word(handle), word(uintptr(channel)), word(uintptr(id)), pointer(config),
			slicePointer(data), word(uintptr(len(data))),
		)
	})
}

func CanTpWriteFirstFrame(device PresetDevice, channel uint8, id uint32, config *PresetTpFrameConfig, firstChunk []byte, totalMessageSize uint32) int32 {
	return withHandle(device.state, func(handle uintptr) int32 {
		return invokeStatus(
			"preset_can_tp_write_first_frame",
			word(handle), word(uintptr(channel)), word(uintptr(id)), pointer(config),
			slicePointer(firstChunk), word(uintptr(len(firstChunk))), word(uintptr(totalMessageSize)),
		)
	})
}

func CanTpWriteFlowControlFrame(device PresetDevice, channel uint8, id uint32, config *PresetTpFrameConfig, flowStatus, blockSize, stMin uint8) int32 {
	return withHandle(device.state, func(handle uintptr) int32 {
		return invokeStatus(
			"preset_can_tp_write_flow_control_frame",
			word(handle), word(uintptr(channel)), word(uintptr(id)), pointer(config),
			word(uintptr(flowStatus)), word(uintptr(blockSize)), word(uintptr(stMin)),
		)
	})
}

func CanTpWriteConsecutiveFrame(device PresetDevice, channel uint8, id uint32, config *PresetTpFrameConfig, dataChunk []byte, sequenceNumber uint8) int32 {
	return withHandle(device.state, func(handle uintptr) int32 {
		return invokeStatus(
			"preset_can_tp_write_consecutive_frame",
			word(handle), word(uintptr(channel)), word(uintptr(id)), pointer(config),
			slicePointer(dataChunk), word(uintptr(len(dataChunk))), word(uintptr(sequenceNumber)),
		)
	})
}

// DeviceClose closes any backend and clears device on success. It returns
// PRESET_ERR_BUSY without changing the handle while a UDS client is active.
func DeviceClose(device *PresetDevice) int32 {
	if device == nil {
		return PRESET_ERR_NULL_PTR
	}
	return closeHandle(device.state, "preset_device_close")
}
