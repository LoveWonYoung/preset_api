package preset_api

import "runtime"

func ToomossOpen() (device PresetDevice, status int32) {
	status = invokeStatus("preset_toomoss_open", pointer(&device))
	runtime.KeepAlive(&device)
	return
}

func ToomossCanInit(device PresetDevice, config *PresetToomossConfig) int32 {
	status := invokeStatus("preset_toomoss_can_init", uintptr(device), pointer(config))
	runtime.KeepAlive(config)
	return status
}

func PcanOpen(config *PresetPCANConfig) (device PresetDevice, status int32) {
	status = invokeStatus("preset_pcan_open", pointer(config), pointer(&device))
	runtime.KeepAlive(config)
	return
}

func PcanCanInit(device PresetDevice, config *PresetPCANConfig) int32 {
	status := invokeStatus("preset_pcan_can_init", uintptr(device), pointer(config))
	runtime.KeepAlive(config)
	return status
}

func TsmasterOpen(config *PresetTSMasterConfig) (device PresetDevice, status int32) {
	status = invokeStatus("preset_tsmaster_open", pointer(config), pointer(&device))
	runtime.KeepAlive(config)
	return
}

func TsmasterCanInit(device PresetDevice, config *PresetTSMasterConfig) int32 {
	status := invokeStatus("preset_tsmaster_can_init", uintptr(device), pointer(config))
	runtime.KeepAlive(config)
	return status
}

func VectorOpen(config *PresetVectorConfig) (device PresetDevice, status int32) {
	status = invokeStatus("preset_vector_open", pointer(config), pointer(&device))
	runtime.KeepAlive(config)
	return
}

func VectorCanInit(device PresetDevice, config *PresetVectorConfig) int32 {
	status := invokeStatus("preset_vector_can_init", uintptr(device), pointer(config))
	runtime.KeepAlive(config)
	return status
}

// DeviceClose closes any backend and clears device on success. It returns
// PRESET_ERR_BUSY without changing the handle while a UDS client is active.
func DeviceClose(device *PresetDevice) int32 {
	status := invokeStatus("preset_device_close", pointer(device))
	runtime.KeepAlive(device)
	return status
}
