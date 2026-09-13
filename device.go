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

// DeviceClose closes any backend and clears device on success. It returns
// PRESET_ERR_BUSY without changing the handle while a UDS client is active.
func DeviceClose(device *PresetDevice) int32 {
	if device == nil {
		return PRESET_ERR_NULL_PTR
	}
	return closeHandle(device.state, "preset_device_close")
}
