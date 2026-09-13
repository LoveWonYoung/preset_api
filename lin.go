//go:build windows && amd64

package preset_api

func ToomossLinInit(device PresetDevice, config *PresetToomossLinConfig) int32 {
	return withHandle(device.state, func(handle uintptr) int32 {
		return invokeStatus("preset_toomoss_lin_init", word(handle), pointer(config))
	})
}

func ToomossLinWrite(device PresetDevice, channel, frameID uint8, data []byte) int32 {
	return withHandle(device.state, func(handle uintptr) int32 {
		return invokeStatus(
			"preset_toomoss_lin_write",
			word(handle), word(uintptr(channel)), word(uintptr(frameID)),
			slicePointer(data), word(uintptr(len(data))),
		)
	})
}

func ToomossLinRead(device PresetDevice, channel, frameID uint8) (frame PresetToomossLinFrame, status int32) {
	status = withHandle(device.state, func(handle uintptr) int32 {
		return invokeStatus(
			"preset_toomoss_lin_read",
			word(handle), word(uintptr(channel)), word(uintptr(frameID)), pointer(&frame),
		)
	})
	return
}

func ToomossLinTryRead(device PresetDevice, channel uint8, frames []PresetToomossLinFrame) (count int, status int32) {
	length := uintptr(len(frames))
	status = withHandle(device.state, func(handle uintptr) int32 {
		return invokeStatus(
			"preset_toomoss_lin_try_read",
			word(handle), word(uintptr(channel)), slicePointer(frames), pointer(&length),
		)
	})
	return int(length), status
}

func ToomossLinBreak(device PresetDevice, channel uint8) int32 {
	return withHandle(device.state, func(handle uintptr) int32 {
		return invokeStatus("preset_toomoss_lin_break", word(handle), word(uintptr(channel)))
	})
}

func ToomossLinSetPower(device PresetDevice, channel, voltage uint8) int32 {
	return withHandle(device.state, func(handle uintptr) int32 {
		return invokeStatus(
			"preset_toomoss_lin_set_power",
			word(handle), word(uintptr(channel)), word(uintptr(voltage)),
		)
	})
}

func ToomossElinsInit(device PresetDevice, config *PresetToomossElinsConfig) int32 {
	return withHandle(device.state, func(handle uintptr) int32 {
		return invokeStatus("preset_toomoss_elins_init", word(handle), pointer(config))
	})
}

func ToomossElinsRead(device PresetDevice, channel uint8, messages []PresetToomossElinsMessage) (count int, status int32) {
	length := uintptr(len(messages))
	status = withHandle(device.state, func(handle uintptr) int32 {
		return invokeStatus(
			"preset_toomoss_elins_read",
			word(handle), word(uintptr(channel)), slicePointer(messages), pointer(&length),
		)
	})
	return int(length), status
}

func PcanLinOpen(config *PresetPCANLinConfig) (device PresetDevice, status int32) {
	var handle uintptr
	status = invokeStatus("preset_pcan_lin_open", pointer(config), pointer(&handle))
	device.state = newHandle(handle)
	return
}

func PcanLinInit(device PresetDevice, config *PresetPCANLinConfig) int32 {
	return withHandle(device.state, func(handle uintptr) int32 {
		return invokeStatus("preset_pcan_lin_init", word(handle), pointer(config))
	})
}

func TsmasterLinOpen(config *PresetTSMasterLinConfig) (device PresetDevice, status int32) {
	var handle uintptr
	status = invokeStatus("preset_tsmaster_lin_open", pointer(config), pointer(&handle))
	device.state = newHandle(handle)
	return
}

func TsmasterLinInit(device PresetDevice, config *PresetTSMasterLinConfig) int32 {
	return withHandle(device.state, func(handle uintptr) int32 {
		return invokeStatus("preset_tsmaster_lin_init", word(handle), pointer(config))
	})
}

func VectorLinOpen(config *PresetVectorLinConfig) (device PresetDevice, status int32) {
	var handle uintptr
	status = invokeStatus("preset_vector_lin_open", pointer(config), pointer(&handle))
	device.state = newHandle(handle)
	return
}

func VectorLinInit(device PresetDevice, config *PresetVectorLinConfig) int32 {
	return withHandle(device.state, func(handle uintptr) int32 {
		return invokeStatus("preset_vector_lin_init", word(handle), pointer(config))
	})
}

func LinMasterWrite(device PresetDevice, channel, frameID uint8, data []byte) int32 {
	return withHandle(device.state, func(handle uintptr) int32 {
		return invokeStatus(
			"preset_lin_master_write",
			word(handle), word(uintptr(channel)), word(uintptr(frameID)),
			slicePointer(data), word(uintptr(len(data))),
		)
	})
}

func LinMasterRead(device PresetDevice, channel, frameID uint8, timeoutMS uint32) (frame PresetLinFrame, status int32) {
	status = withHandle(device.state, func(handle uintptr) int32 {
		return invokeStatus(
			"preset_lin_master_read",
			word(handle), word(uintptr(channel)), word(uintptr(frameID)),
			word(uintptr(timeoutMS)), pointer(&frame),
		)
	})
	return
}

func LinUdsClientNew(device PresetDevice, channel, nad uint8) (client PresetLinUdsClient, status int32) {
	var handle uintptr
	status = withHandle(device.state, func(deviceHandle uintptr) int32 {
		return invokeStatus(
			"preset_lin_uds_client_new",
			word(deviceHandle), word(uintptr(channel)), word(uintptr(nad)), pointer(&handle),
		)
	})
	client.state = newHandle(handle)
	return
}

// ToomossLinUdsClientNew is kept for source compatibility. The ABI v5 LIN UDS
// constructor is backend-independent, so new code should use LinUdsClientNew.
func ToomossLinUdsClientNew(device PresetDevice, channel, nad uint8) (client PresetLinUdsClient, status int32) {
	return LinUdsClientNew(device, channel, nad)
}

func LinUdsRequest(client PresetLinUdsClient, payload []byte, timeoutMS uint32, out []byte) (nad uint8, count int, status int32) {
	length := uintptr(0)
	status = withHandle(client.state, func(handle uintptr) int32 {
		return invokeStatus(
			"preset_lin_uds_request",
			word(handle), slicePointer(payload), word(uintptr(len(payload))), word(uintptr(timeoutMS)),
			pointer(&nad), slicePointer(out), word(uintptr(len(out))), pointer(&length),
		)
	})
	return nad, int(length), status
}

func LinUdsSetPollInterval(client PresetLinUdsClient, intervalMS uint32) int32 {
	return withHandle(client.state, func(handle uintptr) int32 {
		return invokeStatus("preset_lin_uds_set_poll_interval", word(handle), word(uintptr(intervalMS)))
	})
}

func LinUdsClientClose(client *PresetLinUdsClient) int32 {
	if client == nil {
		return PRESET_ERR_NULL_PTR
	}
	return closeHandle(client.state, "preset_lin_uds_client_close")
}
