package main

import (
	"fmt"
	"unsafe"
)

const (
	LinProtocol13 byte = iota
	LinProtocol20
	LinProtocol21
	LinFunctionalNAD   byte = 0x7E
	LinBroadcastNAD    byte = 0x7F
	ToomossLinPower0V  byte = 0
	ToomossLinPower12V byte = 1
	ToomossLinPower5V  byte = 2
)

type PCANLinCfg struct {
	Channel        byte
	Reserved       byte
	HardwareHandle uint16
	Baudrate       uint32
}

type TSMasterLinCfg struct {
	ApplicationChannel byte
	HardwareChannel    byte
	Protocol           byte
	Reserved           byte
	HardwareIndex      int32
	DeviceType         int32
	Baudrate           uint32
}

type VectorLinCfg struct {
	ApplicationChannel byte
	Version            byte
	Reserved           [2]byte
	HardwareType       int32
	HardwareIndex      uint32
	HardwareChannel    uint32
	Baudrate           uint32
	RxQueueSize        uint32
}

type ToomossElinsCfg struct {
	ChannelsMask     byte
	ResistorEnabled  byte
	Version          byte
	Reserved         byte
	Baudrate         uint32
	ReceiveTimeoutUs uint32
}

type ToomossElinsMessage struct {
	DataLen         byte
	BreakBits       byte
	Status          byte
	Flags           byte
	Sync            byte
	TimestampHigh   byte
	MsgSendTimes    uint16
	Timestamp       uint32
	CmdCode         byte
	DeviceID        byte
	RegisterAddress uint16
	CRC16           uint16
	Data            [64]byte
	AckValue        [4]byte
}

func PCANLinDefaultCfg() (PCANLinCfg, error) {
	return defaultStruct[PCANLinCfg]("preset_pcan_lin_default_config")
}
func TSMasterLinDefaultCfg() (TSMasterLinCfg, error) {
	return defaultStruct[TSMasterLinCfg]("preset_tsmaster_lin_default_config")
}
func VectorLinDefaultCfg() (VectorLinCfg, error) {
	return defaultStruct[VectorLinCfg]("preset_vector_lin_default_config")
}
func ToomossElinsDefaultCfg() (ToomossElinsCfg, error) {
	return defaultStruct[ToomossElinsCfg]("preset_toomoss_elins_default_config")
}

func (t *PCAN) LinInitCfg(cfg *PCANLinCfg, slave bool) error {
	if cfg == nil {
		return fmt.Errorf("PCAN LIN config is nil")
	}
	name := "preset_pcan_lin_init"
	if slave {
		name += "_slave"
	}
	return nativeStatus(name, t.device, uintptr(unsafe.Pointer(cfg)))
}

func (t *TSMaster) LinInitCfg(cfg *TSMasterLinCfg, slave bool) error {
	if cfg == nil {
		return fmt.Errorf("TSMaster LIN config is nil")
	}
	name := "preset_tsmaster_lin_init"
	if slave {
		name += "_slave"
	}
	return nativeStatus(name, t.device, uintptr(unsafe.Pointer(cfg)))
}

func (t *Vector) LinInitCfg(cfg *VectorLinCfg, slave bool) error {
	if cfg == nil {
		return fmt.Errorf("Vector LIN config is nil")
	}
	name := "preset_vector_lin_init"
	if slave {
		name += "_slave"
	}
	return nativeStatus(name, t.device, uintptr(unsafe.Pointer(cfg)))
}

func (t *Toomoss) LinInitCfg(cfg *ToomossLinCfg) error {
	if cfg == nil {
		return fmt.Errorf("Toomoss LIN config is nil")
	}
	return nativeStatus("preset_toomoss_lin_init", t.device, uintptr(unsafe.Pointer(cfg)))
}

func (t *Toomoss) LinSetPower(channel, voltage byte) error {
	return nativeStatus("preset_toomoss_lin_set_power", t.device, uintptr(channel), uintptr(voltage))
}

func (t *Toomoss) LinTryRead(channel byte, maxFrames int) ([]LinFrame, error) {
	if maxFrames < 1 || maxFrames > 512 {
		return nil, fmt.Errorf("LIN read capacity must be 1-512")
	}
	frames := make([]LinFrame, maxFrames)
	n := uintptr(maxFrames)
	err := nativeStatus("preset_toomoss_lin_try_read", t.device, uintptr(channel), uintptr(unsafe.Pointer(&frames[0])), uintptr(unsafe.Pointer(&n)))
	if err != nil {
		return nil, err
	}
	return frames[:n], nil
}

func (t *Toomoss) ElinsInit(cfg *ToomossElinsCfg) error {
	if cfg == nil {
		return fmt.Errorf("ELINS config is nil")
	}
	return nativeStatus("preset_toomoss_elins_init", t.device, uintptr(unsafe.Pointer(cfg)))
}

func (t *Toomoss) ElinsRead(channel byte, maxMessages int) ([]ToomossElinsMessage, error) {
	if maxMessages < 1 {
		return nil, fmt.Errorf("ELINS read capacity must be positive")
	}
	msgs := make([]ToomossElinsMessage, maxMessages)
	n := uintptr(maxMessages)
	err := nativeStatus("preset_toomoss_elins_read", t.device, uintptr(channel), uintptr(unsafe.Pointer(&msgs[0])), uintptr(unsafe.Pointer(&n)))
	if err != nil {
		return nil, err
	}
	return msgs[:n], nil
}

// LinWrite and LinRead work on initialized LIN master channels on any backend.
func (t *Device) LinWrite(channel, id byte, data []byte) error {
	if len(data) < 1 || len(data) > 8 {
		return fmt.Errorf("LIN payload must contain 1-8 bytes")
	}
	return nativeStatus("preset_lin_master_write", t.device, uintptr(channel), uintptr(id), bytePtr(data), uintptr(len(data)))
}

func (t *Device) LinRead(channel, id byte, timeoutMs uint32) (LinFrame, error) {
	var frame LinFrame
	err := nativeStatus("preset_lin_master_read", t.device, uintptr(channel), uintptr(id), uintptr(timeoutMs), uintptr(unsafe.Pointer(&frame)))
	return frame, err
}

func (t *Device) LinWakeup(channel byte) error {
	return nativeStatus("preset_lin_wakeup", t.device, uintptr(channel))
}

func (t *Device) LinSlaveSetResponse(channel, id byte, data []byte, classicChecksum, singleShot bool) error {
	if len(data) < 1 || len(data) > 8 {
		return fmt.Errorf("LIN response must contain 1-8 bytes")
	}
	return nativeStatus("preset_lin_slave_set_response", t.device, uintptr(channel), uintptr(id), bytePtr(data), uintptr(len(data)), uintptr(boolByte(classicChecksum)), uintptr(boolByte(singleShot)))
}

func (t *Device) LinSlaveClearResponse(channel, id byte) error {
	return nativeStatus("preset_lin_slave_clear_response", t.device, uintptr(channel), uintptr(id))
}

type LinUdsClient struct{ handle uintptr }

// LinUdsClientNew borrows an initialized master channel; close it before Device.Close.
func (t *Device) LinUdsClientNew(channel, nad byte) (*LinUdsClient, error) {
	var handle uintptr
	if err := nativeStatus("preset_lin_uds_client_new", t.device, uintptr(channel), uintptr(nad), uintptr(unsafe.Pointer(&handle))); err != nil {
		return nil, err
	}
	return &LinUdsClient{handle: handle}, nil
}

func (c *LinUdsClient) SetPollInterval(ms uint32) error {
	return nativeStatus("preset_lin_uds_set_poll_interval", c.handle, uintptr(ms))
}

func (c *LinUdsClient) SetRecoveryInterval(ms uint32) error {
	return nativeStatus("preset_lin_uds_set_recovery_interval", c.handle, uintptr(ms))
}

func (c *LinUdsClient) WaitIdle(ms uint32) (bool, error) {
	p, err := proc("preset_lin_uds_wait_idle")
	if err != nil {
		return false, err
	}
	rc := int32(syscallR1(p, c.handle, uintptr(ms)))
	if rc == PresetErrWaitTimeout {
		return false, nil
	}
	if rc < 0 {
		return false, &PresetError{Code: rc, Func: "preset_lin_uds_wait_idle", Message: LastError()}
	}
	return true, nil
}

func (c *LinUdsClient) StartRequest(timeouts *RequestTimeouts, pdu []byte) (*DiagnosticRequest, error) {
	if timeouts == nil || len(pdu) == 0 {
		return nil, fmt.Errorf("timeouts and PDU are required")
	}
	var handle uintptr
	if err := nativeStatus("preset_lin_request_start", c.handle, uintptr(unsafe.Pointer(timeouts)), bytePtr(pdu), uintptr(len(pdu)), uintptr(unsafe.Pointer(&handle))); err != nil {
		return nil, err
	}
	return &DiagnosticRequest{handle: handle}, nil
}

func (c *LinUdsClient) UdsRequest(timeouts *RequestTimeouts, pdu []byte, waitMs uint32) (DiagnosticRequestInfo, []byte, error) {
	q, err := c.StartRequest(timeouts, pdu)
	if err != nil {
		return DiagnosticRequestInfo{}, nil, err
	}
	return finishRequest(q, waitMs)
}

func (c *LinUdsClient) Close() error {
	if c == nil || c.handle == 0 {
		return nil
	}
	return nativeStatus("preset_lin_uds_client_close", uintptr(unsafe.Pointer(&c.handle)))
}
