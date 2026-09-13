package preset_api

import "runtime"

func ToomossLinInit(device PresetDevice, config *PresetToomossLinConfig) int32 {
	status := invokeStatus("preset_toomoss_lin_init", uintptr(device), pointer(config))
	runtime.KeepAlive(config)
	return status
}

func ToomossLinWrite(device PresetDevice, channel, frameID uint8, data []byte) int32 {
	status := invokeStatus(
		"preset_toomoss_lin_write",
		uintptr(device), uintptr(channel), uintptr(frameID),
		slicePointer(data), uintptr(len(data)),
	)
	runtime.KeepAlive(data)
	return status
}

func ToomossLinRead(device PresetDevice, channel, frameID uint8) (frame PresetToomossLinFrame, status int32) {
	status = invokeStatus(
		"preset_toomoss_lin_read",
		uintptr(device), uintptr(channel), uintptr(frameID), pointer(&frame),
	)
	runtime.KeepAlive(&frame)
	return
}

func ToomossLinTryRead(device PresetDevice, channel uint8, frames []PresetToomossLinFrame) (count int, status int32) {
	length := uintptr(len(frames))
	status = invokeStatus(
		"preset_toomoss_lin_try_read",
		uintptr(device), uintptr(channel), slicePointer(frames), pointer(&length),
	)
	runtime.KeepAlive(frames)
	return int(length), status
}

func ToomossLinBreak(device PresetDevice, channel uint8) int32 {
	return invokeStatus("preset_toomoss_lin_break", uintptr(device), uintptr(channel))
}

func ToomossLinSetPower(device PresetDevice, channel, voltage uint8) int32 {
	return invokeStatus("preset_toomoss_lin_set_power", uintptr(device), uintptr(channel), uintptr(voltage))
}

func ToomossElinsInit(device PresetDevice, config *PresetToomossElinsConfig) int32 {
	status := invokeStatus("preset_toomoss_elins_init", uintptr(device), pointer(config))
	runtime.KeepAlive(config)
	return status
}

func ToomossElinsRead(device PresetDevice, channel uint8, messages []PresetToomossElinsMessage) (count int, status int32) {
	length := uintptr(len(messages))
	status = invokeStatus(
		"preset_toomoss_elins_read",
		uintptr(device), uintptr(channel), slicePointer(messages), pointer(&length),
	)
	runtime.KeepAlive(messages)
	return int(length), status
}

func ToomossLinUdsClientNew(device PresetDevice, channel, nad uint8) (client PresetLinUdsClient, status int32) {
	status = invokeStatus(
		"preset_toomoss_lin_uds_client_new",
		uintptr(device), uintptr(channel), uintptr(nad), pointer(&client),
	)
	runtime.KeepAlive(&client)
	return
}

func LinUdsRequest(client PresetLinUdsClient, payload []byte, timeoutMS uint32, out []byte) (nad uint8, count int, status int32) {
	length := uintptr(0)
	status = invokeStatus(
		"preset_lin_uds_request",
		uintptr(client), slicePointer(payload), uintptr(len(payload)), uintptr(timeoutMS),
		pointer(&nad), slicePointer(out), uintptr(len(out)), pointer(&length),
	)
	runtime.KeepAlive(payload)
	runtime.KeepAlive(out)
	return nad, int(length), status
}

func LinUdsSetPollInterval(client PresetLinUdsClient, intervalMS uint32) int32 {
	return invokeStatus("preset_lin_uds_set_poll_interval", uintptr(client), uintptr(intervalMS))
}

func LinUdsClientClose(client *PresetLinUdsClient) int32 {
	status := invokeStatus("preset_lin_uds_client_close", pointer(client))
	runtime.KeepAlive(client)
	return status
}
