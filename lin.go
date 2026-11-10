//go:build windows && amd64

package preset

import (
	"fmt"
	"unsafe"
)

const (
	LINProtocol13 byte = iota
	LINProtocol20
	LINProtocol21
	LINFunctionalNAD   byte = 0x7E
	LINBroadcastNAD    byte = 0x7F
	ToomossLINPower0V  byte = 0
	ToomossLINPower12V byte = 1
	ToomossLINPower5V  byte = 2
)

type ToomossLINConfig struct {
	Channel    byte
	MasterMode byte
	Baudrate   uint32
}

type PCANLINConfig struct {
	Channel        byte
	Reserved       byte
	HardwareHandle uint16
	Baudrate       uint32
}

type TSMasterLINConfig struct {
	ApplicationChannel byte
	HardwareChannel    byte
	Protocol           byte
	Reserved           byte
	HardwareIndex      int32
	DeviceType         int32
	Baudrate           uint32
}

type VectorLINConfig struct {
	ApplicationChannel byte
	Version            byte
	Reserved           [2]byte
	HardwareType       int32
	HardwareIndex      uint32
	HardwareChannel    uint32
	Baudrate           uint32
	RXQueueSize        uint32
}

type LINFrame struct {
	ID              byte
	DataLen         byte
	ClassicChecksum byte // 1 classic, 0 enhanced
	ChecksumKnown   byte // 0 unknown, 1 driver or channel configuration
	Data            [8]byte
	Timestamp       uint32
	TimestampUS     uint64
}

type ToomossELINSConfig struct {
	ChannelsMask     byte
	ResistorEnabled  byte
	Version          byte
	Reserved         byte
	Baudrate         uint32
	ReceiveTimeoutUs uint32
}

type ToomossELINSMessage struct {
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

func DefaultToomossLINConfig() (ToomossLINConfig, error) {
	return defaultStruct[ToomossLINConfig]("preset_toomoss_lin_default_config")
}
func DefaultPCANLINConfig() (PCANLINConfig, error) {
	return defaultStruct[PCANLINConfig]("preset_pcan_lin_default_config")
}
func DefaultTSMasterLINConfig() (TSMasterLINConfig, error) {
	return defaultStruct[TSMasterLINConfig]("preset_tsmaster_lin_default_config")
}
func DefaultVectorLINConfig() (VectorLINConfig, error) {
	return defaultStruct[VectorLINConfig]("preset_vector_lin_default_config")
}
func DefaultToomossELINSConfig() (ToomossELINSConfig, error) {
	return defaultStruct[ToomossELINSConfig]("preset_toomoss_elins_default_config")
}

func (t *PCAN) InitLINWithConfig(cfg *PCANLINConfig, slave bool) error {
	if cfg == nil {
		return fmt.Errorf("PCAN LIN config is nil")
	}
	name := "preset_pcan_lin_init"
	if slave {
		name += "_slave"
	}
	return t.handle.use(func(handle uintptr) error {
		return nativeStatus(name, handle, uintptr(unsafe.Pointer(cfg)))
	})
}

func (t *TSMaster) InitLINWithConfig(cfg *TSMasterLINConfig, slave bool) error {
	if cfg == nil {
		return fmt.Errorf("TSMaster LIN config is nil")
	}
	name := "preset_tsmaster_lin_init"
	if slave {
		name += "_slave"
	}
	return t.handle.use(func(handle uintptr) error {
		return nativeStatus(name, handle, uintptr(unsafe.Pointer(cfg)))
	})
}

func (t *Vector) InitLINWithConfig(cfg *VectorLINConfig, slave bool) error {
	if cfg == nil {
		return fmt.Errorf("Vector LIN config is nil")
	}
	name := "preset_vector_lin_init"
	if slave {
		name += "_slave"
	}
	return t.handle.use(func(handle uintptr) error {
		return nativeStatus(name, handle, uintptr(unsafe.Pointer(cfg)))
	})
}

func (t *Toomoss) InitLINWithConfig(cfg *ToomossLINConfig) error {
	if cfg == nil {
		return fmt.Errorf("Toomoss LIN config is nil")
	}
	return t.handle.use(func(handle uintptr) error {
		return nativeStatus("preset_toomoss_lin_init", handle, uintptr(unsafe.Pointer(cfg)))
	})
}

// InitLIN initializes each supplied zero-based Toomoss LIN channel with the
// library's default configuration.
func (t *Toomoss) InitLIN(channels ...byte) error {
	for _, channel := range channels {
		if channel > 3 {
			return fmt.Errorf("Toomoss LIN channel must be in range 0-3")
		}
		cfg, err := DefaultToomossLINConfig()
		if err != nil {
			return err
		}
		cfg.Channel = channel
		if err := t.InitLINWithConfig(&cfg); err != nil {
			return err
		}
	}
	return nil
}

func (t *Toomoss) SetLINPower(channel, voltage byte) error {
	return t.handle.use(func(handle uintptr) error {
		return nativeStatus("preset_toomoss_lin_set_power", handle, uintptr(channel), uintptr(voltage))
	})
}

func (t *Toomoss) TryReadLIN(channel byte, maxFrames int) ([]LINFrame, error) {
	if maxFrames < 1 || maxFrames > 512 {
		return nil, fmt.Errorf("LIN read capacity must be 1-512")
	}
	frames := make([]LINFrame, maxFrames)
	n := uintptr(maxFrames)
	err := t.handle.use(func(handle uintptr) error {
		return nativeStatus("preset_toomoss_lin_try_read", handle, uintptr(channel), uintptr(unsafe.Pointer(&frames[0])), uintptr(unsafe.Pointer(&n)))
	})
	if err != nil {
		return nil, err
	}
	if n > uintptr(len(frames)) {
		return nil, fmt.Errorf("preset_toomoss_lin_try_read returned invalid length %d", n)
	}
	return frames[:n], nil
}

func (t *Toomoss) InitELINS(cfg *ToomossELINSConfig) error {
	if cfg == nil {
		return fmt.Errorf("ELINS config is nil")
	}
	return t.handle.use(func(handle uintptr) error {
		return nativeStatus("preset_toomoss_elins_init", handle, uintptr(unsafe.Pointer(cfg)))
	})
}

func (t *Toomoss) ReadELINS(channel byte, maxMessages int) ([]ToomossELINSMessage, error) {
	if maxMessages < 1 {
		return nil, fmt.Errorf("ELINS read capacity must be positive")
	}
	msgs := make([]ToomossELINSMessage, maxMessages)
	n := uintptr(maxMessages)
	err := t.handle.use(func(handle uintptr) error {
		return nativeStatus("preset_toomoss_elins_read", handle, uintptr(channel), uintptr(unsafe.Pointer(&msgs[0])), uintptr(unsafe.Pointer(&n)))
	})
	if err != nil {
		return nil, err
	}
	if n > uintptr(len(msgs)) {
		return nil, fmt.Errorf("preset_toomoss_elins_read returned invalid length %d", n)
	}
	return msgs[:n], nil
}

// WriteLIN and ReadLIN work on initialized LIN master channels on any backend.
func (t *Device) WriteLIN(channel, id byte, data []byte) error {
	if len(data) < 1 || len(data) > 8 {
		return fmt.Errorf("LIN payload must contain 1-8 bytes")
	}
	return t.handle.use(func(handle uintptr) error {
		return nativeStatus("preset_lin_master_write", handle, uintptr(channel), uintptr(id), bytePtr(data), uintptr(len(data)))
	})
}

func (t *Device) ReadLIN(channel, id byte, timeoutMs uint32) (LINFrame, error) {
	return useHandle(&t.handle, func(handle uintptr) (LINFrame, error) {
		var frame LINFrame
		err := nativeStatus("preset_lin_master_read", handle, uintptr(channel), uintptr(id), uintptr(timeoutMs), uintptr(unsafe.Pointer(&frame)))
		return frame, err
	})
}

// WakeLIN wakes an initialized LIN channel. On Toomoss this sends a master break.
func (t *Device) WakeLIN(channel byte) error {
	return t.handle.use(func(handle uintptr) error {
		return nativeStatus("preset_lin_wakeup", handle, uintptr(channel))
	})
}

func (t *Device) SetLINSlaveResponse(channel, id byte, data []byte, classicChecksum, singleShot bool) error {
	if len(data) < 1 || len(data) > 8 {
		return fmt.Errorf("LIN response must contain 1-8 bytes")
	}
	return t.handle.use(func(handle uintptr) error {
		return nativeStatus("preset_lin_slave_set_response", handle, uintptr(channel), uintptr(id), bytePtr(data), uintptr(len(data)), uintptr(boolByte(classicChecksum)), uintptr(boolByte(singleShot)))
	})
}

func (t *Device) ClearLINSlaveResponse(channel, id byte) error {
	return t.handle.use(func(handle uintptr) error {
		return nativeStatus("preset_lin_slave_clear_response", handle, uintptr(channel), uintptr(id))
	})
}

type LINUDSClient struct{ handle nativeHandle }

// NewLINUDSClient borrows an initialized master channel; close it before Device.Close.
func (t *Device) NewLINUDSClient(channel, nad byte) (*LINUDSClient, error) {
	return useHandle(&t.handle, func(device uintptr) (*LINUDSClient, error) {
		var handle uintptr
		if err := nativeStatus("preset_lin_uds_client_new", device, uintptr(channel), uintptr(nad), uintptr(unsafe.Pointer(&handle))); err != nil {
			return nil, err
		}
		return &LINUDSClient{handle: newNativeHandle(handle)}, nil
	})
}

func (c *LINUDSClient) SetPollInterval(ms uint32) error {
	return c.handle.use(func(handle uintptr) error {
		return nativeStatus("preset_lin_uds_set_poll_interval", handle, uintptr(ms))
	})
}

func (c *LINUDSClient) SetRecoveryInterval(ms uint32) error {
	return c.handle.use(func(handle uintptr) error {
		return nativeStatus("preset_lin_uds_set_recovery_interval", handle, uintptr(ms))
	})
}

func (c *LINUDSClient) WaitIdle(ms uint32) (bool, error) {
	return useHandle(&c.handle, func(handle uintptr) (bool, error) {
		return nativeWait("preset_lin_uds_wait_idle", handle, uintptr(ms))
	})
}

func (c *LINUDSClient) StartRequest(timeouts *RequestTimeouts, pdu []byte) (*DiagnosticRequest, error) {
	if timeouts == nil || len(pdu) == 0 {
		return nil, fmt.Errorf("timeouts and PDU are required")
	}
	return useHandle(&c.handle, func(client uintptr) (*DiagnosticRequest, error) {
		var handle uintptr
		if err := nativeStatus("preset_lin_request_start", client, uintptr(unsafe.Pointer(timeouts)), bytePtr(pdu), uintptr(len(pdu)), uintptr(unsafe.Pointer(&handle))); err != nil {
			return nil, err
		}
		return &DiagnosticRequest{handle: newNativeHandle(handle)}, nil
	})
}

func (c *LINUDSClient) UDSRequest(timeouts *RequestTimeouts, pdu []byte, waitMs uint32) (DiagnosticRequestInfo, []byte, error) {
	q, err := c.StartRequest(timeouts, pdu)
	if err != nil {
		return DiagnosticRequestInfo{}, nil, err
	}
	return finishRequest(q, waitMs)
}

func (c *LINUDSClient) Close() error {
	if c == nil {
		return nil
	}
	return c.handle.close("preset_lin_uds_client_close")
}
