//go:build windows && amd64

package preset_api

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Backend identifies a CAN adapter. Values match PRESET_CAN_BACKEND_*.
type Backend int

const (
	BackendToomoss  Backend = PRESET_CAN_BACKEND_TOOMOSS
	BackendTSMaster Backend = PRESET_CAN_BACKEND_TSMASTER
	BackendPcan     Backend = PRESET_CAN_BACKEND_PCAN
	BackendVector   Backend = PRESET_CAN_BACKEND_VECTOR
)

func (backend Backend) String() string {
	switch backend {
	case BackendToomoss:
		return "Toomoss"
	case BackendTSMaster:
		return "TSMaster"
	case BackendPcan:
		return "PCAN"
	case BackendVector:
		return "Vector"
	default:
		return fmt.Sprintf("Backend(%d)", int(backend))
	}
}

// CanFrame is a raw CAN/CAN-FD frame returned by MyDevice.Rxfn.
type CanFrame struct {
	ID        uint32
	DLC       uint8
	Data      []byte
	Direction string
	IsFD      bool
	BRS       bool
}

func canFrameFromEx(frame PresetCanFrameEx) CanFrame {
	direction := "RX"
	if frame.Direction == PRESET_CAN_DIRECTION_TX {
		direction = "TX"
	}
	data := make([]byte, frame.DataLen)
	copy(data, frame.Data[:frame.DataLen])
	return CanFrame{
		ID:        frame.ID,
		DLC:       frame.DLC,
		Data:      data,
		Direction: direction,
		IsFD:      frame.IsFD != 0,
		BRS:       frame.BRS != 0,
	}
}

// MyDevice wraps DLL loading, backend open, channel init, and UDS clients.
// Callers only construct this type, Open it, then Txfn/Rxfn/Request.
type MyDevice struct {
	Backend        Backend
	PhysID         uint32
	RespID         uint32
	FuncID         uint32
	Channels       []uint8
	DeviceType     int32
	DLLPath        string
	IsFD           bool
	BRS            bool
	NominalBitrate uint32
	DataBitrate    uint32
	HardwareIndex  int32
	RawRxEnabled   bool

	device           PresetDevice
	clients          map[uint8]PresetCanUdsClient
	opened           bool
	responseCapacity int
}

// NewMyDevice builds a closed device. deviceType 0 selects the backend default
// (TSMaster TC1016, Vector ANY). Set optional fields before Open.
func NewMyDevice(backend Backend, physID, respID, funcID uint32, channels []uint8, deviceType int32) (*MyDevice, error) {
	if len(channels) == 0 {
		return nil, errors.New("channels must not be empty")
	}
	copied := make([]uint8, len(channels))
	copy(copied, channels)
	return &MyDevice{
		Backend:        backend,
		PhysID:         physID,
		RespID:         respID,
		FuncID:         funcID,
		Channels:       copied,
		DeviceType:     deviceType,
		IsFD:           true,
		BRS:            true,
		NominalBitrate: 500_000,
		DataBitrate:    2_000_000,
		RawRxEnabled:   true,
		clients:        make(map[uint8]PresetCanUdsClient),
	}, nil
}

// Open loads the DLL, opens the adapter, inits every channel, then creates one
// UDS client per channel. Extra TSMaster/PCAN/Vector channels are initialized
// before any client is created.
func (d *MyDevice) Open() (err error) {
	if d == nil {
		return errors.New("MyDevice is nil")
	}
	if d.opened {
		return nil
	}
	if len(d.Channels) == 0 {
		return errors.New("channels must not be empty")
	}
	if d.clients == nil {
		d.clients = make(map[uint8]PresetCanUdsClient)
	}

	path := d.DLLPath
	if path == "" {
		wd, wdErr := os.Getwd()
		if wdErr != nil {
			return wdErr
		}
		path = filepath.Join(wd, "preset_rs.dll")
	}
	if err = LoadDLL(path); err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = d.closeOpened()
		}
	}()

	if err = d.openBackend(); err != nil {
		return err
	}
	if err = d.initExtraChannels(); err != nil {
		return err
	}

	uds, udsErr := d.udsConfig()
	if udsErr != nil {
		return udsErr
	}
	d.responseCapacity = int(uds.MaxPDULen)
	if d.responseCapacity <= 0 {
		d.responseCapacity = 4096
	}
	for _, channel := range d.Channels {
		cfg := uds
		var client PresetCanUdsClient
		if err = checkCall("preset_can_uds_client_new", func() int32 {
			var status int32
			client, status = CanUdsClientNew(d.device, channel, &cfg)
			return status
		}); err != nil {
			return err
		}
		d.clients[channel] = client
	}
	d.opened = true
	return nil
}

// Txfn writes a raw CAN frame. An omitted channel uses the first configured one.
func (d *MyDevice) Txfn(canID uint32, data []byte, channel ...uint8) error {
	client, err := d.client(channel...)
	if err != nil {
		return err
	}
	return checkCall("preset_can_uds_write", func() int32 {
		return CanUdsWrite(client, canID, d.IsFD, data)
	})
}

// Rxfn reads pending raw frames. timeoutMS <= 0 returns immediately; otherwise
// it waits until at least one frame arrives or the timeout expires.
func (d *MyDevice) Rxfn(timeoutMS int, channel ...uint8) ([]CanFrame, error) {
	client, err := d.client(channel...)
	if err != nil {
		return nil, err
	}
	if timeoutMS <= 0 {
		return d.readOnce(client, 256)
	}

	var frames []CanFrame
	deadline := time.Now().Add(time.Duration(timeoutMS) * time.Millisecond)
	for {
		batch, readErr := d.readOnce(client, 256)
		if readErr != nil {
			return frames, readErr
		}
		frames = append(frames, batch...)
		if len(frames) > 0 || !time.Now().Before(deadline) {
			return frames, nil
		}
		time.Sleep(20 * time.Millisecond)
	}
}

// Request sends a physical UDS request and returns the response payload.
func (d *MyDevice) Request(payload []byte, timeoutMS uint32, channel ...uint8) ([]byte, error) {
	return d.request(false, payload, timeoutMS, channel...)
}

// FunctionalRequest sends a functional UDS request.
func (d *MyDevice) FunctionalRequest(payload []byte, timeoutMS uint32, channel ...uint8) ([]byte, error) {
	return d.request(true, payload, timeoutMS, channel...)
}

// Close shuts down UDS clients first, then the device. It is safe on a
// never-opened or already-closed device.
func (d *MyDevice) Close() error {
	if d == nil {
		return nil
	}
	return d.closeOpened()
}

func (d *MyDevice) request(functional bool, payload []byte, timeoutMS uint32, channel ...uint8) ([]byte, error) {
	client, err := d.client(channel...)
	if err != nil {
		return nil, err
	}
	op := "preset_can_uds_request"
	call := CanUdsRequest
	if functional {
		op = "preset_can_uds_functional_request"
		call = CanUdsFunctionalRequest
	}

	out := make([]byte, d.responseCapacity)
	n, status, message := callWithLastError(func() (int, int32) {
		return call(client, payload, timeoutMS, out)
	})
	if status == PRESET_ERR_BUFFER_TOO_SMALL && n > 0 {
		out = make([]byte, n)
		n, status, message = callWithLastError(func() (int, int32) {
			return call(client, payload, timeoutMS, out)
		})
	}
	if err = statusErr(op, status, message); err != nil {
		return nil, err
	}
	result := make([]byte, n)
	copy(result, out[:n])
	return result, nil
}

func (d *MyDevice) readOnce(client PresetCanUdsClient, capacity int) ([]CanFrame, error) {
	buf := make([]PresetCanFrameEx, capacity)
	n, status, message := callWithLastError(func() (int, int32) {
		return CanUdsTryReadEx(client, buf)
	})
	if err := statusErr("preset_can_uds_try_read_ex", status, message); err != nil {
		return nil, err
	}
	if n < 0 {
		n = 0
	}
	if n > len(buf) {
		return nil, fmt.Errorf("preset_can_uds_try_read_ex: invalid length %d", n)
	}
	frames := make([]CanFrame, 0, n)
	for i := 0; i < n; i++ {
		frames = append(frames, canFrameFromEx(buf[i]))
	}
	return frames, nil
}

func (d *MyDevice) client(channel ...uint8) (PresetCanUdsClient, error) {
	if d == nil || !d.opened {
		return PresetCanUdsClient{}, errors.New("device is not open; call Open() first")
	}
	var ch uint8
	switch len(channel) {
	case 0:
		ch = d.Channels[0]
	case 1:
		ch = channel[0]
	default:
		return PresetCanUdsClient{}, errors.New("at most one channel may be specified")
	}
	client, ok := d.clients[ch]
	if !ok {
		return PresetCanUdsClient{}, fmt.Errorf("channel %d is not initialized; available: %v", ch, d.Channels)
	}
	return client, nil
}

func (d *MyDevice) udsConfig() (PresetConfig, error) {
	cfg, err := DefaultConfig()
	if err != nil {
		return PresetConfig{}, err
	}
	cfg.PhysicalID = d.PhysID
	cfg.ResponseID = d.RespID
	cfg.FunctionalID = d.FuncID
	cfg.IsFD = boolByte(d.IsFD)
	cfg.RawRxEnabled = boolByte(d.RawRxEnabled)
	return cfg, nil
}

func (d *MyDevice) openBackend() error {
	first := d.Channels[0]
	switch d.Backend {
	case BackendToomoss:
		device, status, message := openWithLastError(ToomossOpen)
		if err := statusErr("preset_toomoss_open", status, message); err != nil {
			return err
		}
		d.device = device
		for _, channel := range d.Channels {
			cfg, err := d.toomossConfig(channel)
			if err != nil {
				return err
			}
			if err = checkCall("preset_toomoss_can_init", func() int32 {
				return ToomossCanInit(d.device, &cfg)
			}); err != nil {
				return err
			}
		}
		return nil
	case BackendTSMaster:
		cfg, err := d.tsmasterConfig(first)
		if err != nil {
			return err
		}
		device, status, message := openWithLastError(func() (PresetDevice, int32) {
			return TsmasterOpen(&cfg)
		})
		if err = statusErr("preset_tsmaster_open", status, message); err != nil {
			return err
		}
		d.device = device
		return nil
	case BackendPcan:
		cfg, err := d.pcanConfig(first)
		if err != nil {
			return err
		}
		device, status, message := openWithLastError(func() (PresetDevice, int32) {
			return PcanOpen(&cfg)
		})
		if err = statusErr("preset_pcan_open", status, message); err != nil {
			return err
		}
		d.device = device
		return nil
	case BackendVector:
		cfg, err := d.vectorConfig(first)
		if err != nil {
			return err
		}
		device, status, message := openWithLastError(func() (PresetDevice, int32) {
			return VectorOpen(&cfg)
		})
		if err = statusErr("preset_vector_open", status, message); err != nil {
			return err
		}
		d.device = device
		return nil
	default:
		return fmt.Errorf("unsupported backend: %s", d.Backend)
	}
}

func (d *MyDevice) initExtraChannels() error {
	if len(d.Channels) <= 1 {
		return nil
	}
	switch d.Backend {
	case BackendToomoss:
		return nil
	case BackendTSMaster:
		for _, channel := range d.Channels[1:] {
			cfg, err := d.tsmasterConfig(channel)
			if err != nil {
				return err
			}
			if err = checkCall("preset_tsmaster_can_init", func() int32 {
				return TsmasterCanInit(d.device, &cfg)
			}); err != nil {
				return err
			}
		}
	case BackendPcan:
		for _, channel := range d.Channels[1:] {
			cfg, err := d.pcanConfig(channel)
			if err != nil {
				return err
			}
			if err = checkCall("preset_pcan_can_init", func() int32 {
				return PcanCanInit(d.device, &cfg)
			}); err != nil {
				return err
			}
		}
	case BackendVector:
		for _, channel := range d.Channels[1:] {
			cfg, err := d.vectorConfig(channel)
			if err != nil {
				return err
			}
			if err = checkCall("preset_vector_can_init", func() int32 {
				return VectorCanInit(d.device, &cfg)
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *MyDevice) toomossConfig(channel uint8) (PresetToomossConfig, error) {
	cfg, err := ToomossDefaultConfig()
	if err != nil {
		return PresetToomossConfig{}, err
	}
	cfg.Channel = channel
	cfg.BRS = boolByte(d.BRS)
	cfg.NominalBitrate = d.NominalBitrate
	cfg.DataBitrate = d.DataBitrate
	return cfg, nil
}

func (d *MyDevice) tsmasterConfig(channel uint8) (PresetTSMasterConfig, error) {
	cfg, err := TsmasterDefaultConfig()
	if err != nil {
		return PresetTSMasterConfig{}, err
	}
	cfg.ApplicationChannel = channel
	cfg.HardwareChannel = channel
	cfg.HardwareIndex = d.HardwareIndex
	if d.DeviceType == 0 {
		cfg.DeviceType = PRESET_TSMASTER_TC1016
	} else {
		cfg.DeviceType = d.DeviceType
	}
	cfg.IsFD = boolByte(d.IsFD)
	cfg.BRS = boolByte(d.BRS)
	cfg.NominalBitrate = d.NominalBitrate
	cfg.DataBitrate = d.DataBitrate
	return cfg, nil
}

func (d *MyDevice) pcanConfig(channel uint8) (PresetPCANConfig, error) {
	cfg, err := PcanDefaultConfig()
	if err != nil {
		return PresetPCANConfig{}, err
	}
	cfg.Channel = channel
	cfg.IsFD = boolByte(d.IsFD)
	cfg.BRS = boolByte(d.BRS)
	cfg.NominalBitrate = d.NominalBitrate
	cfg.DataBitrate = d.DataBitrate
	return cfg, nil
}

func (d *MyDevice) vectorConfig(channel uint8) (PresetVectorConfig, error) {
	cfg, err := VectorDefaultConfig()
	if err != nil {
		return PresetVectorConfig{}, err
	}
	cfg.ApplicationChannel = channel
	cfg.HardwareChannel = uint32(channel)
	if d.HardwareIndex < 0 {
		cfg.HardwareIndex = 0
	} else {
		cfg.HardwareIndex = uint32(d.HardwareIndex)
	}
	if d.DeviceType == 0 {
		cfg.HardwareType = PRESET_VECTOR_HWTYPE_ANY
	} else {
		cfg.HardwareType = d.DeviceType
	}
	cfg.IsFD = boolByte(d.IsFD)
	cfg.BRS = boolByte(d.BRS)
	cfg.NominalBitrate = d.NominalBitrate
	cfg.DataBitrate = d.DataBitrate
	return cfg, nil
}

func (d *MyDevice) closeOpened() error {
	var first error
	for channel, client := range d.clients {
		c := client
		if status := CanUdsClientClose(&c); status != PRESET_OK && first == nil {
			first = fmt.Errorf("close channel %d: status %d", channel, status)
		}
		delete(d.clients, channel)
	}
	if !d.device.IsClosed() {
		if status := DeviceClose(&d.device); status != PRESET_OK && first == nil {
			first = fmt.Errorf("close device: status %d", status)
		}
	}
	d.opened = false
	d.responseCapacity = 0
	return first
}

func boolByte(value bool) uint8 {
	if value {
		return 1
	}
	return 0
}

func checkCall(op string, call func() int32) error {
	status, message := WithLastError(call)
	return statusErr(op, status, message)
}

func statusErr(op string, status int32, message string) error {
	if status == PRESET_OK {
		return nil
	}
	if message != "" {
		return fmt.Errorf("%s: status %d (%s)", op, status, message)
	}
	return fmt.Errorf("%s: status %d", op, status)
}

func openWithLastError(call func() (PresetDevice, int32)) (PresetDevice, int32, string) {
	var device PresetDevice
	status, message := WithLastError(func() int32 {
		var openStatus int32
		device, openStatus = call()
		return openStatus
	})
	return device, status, message
}

func callWithLastError(call func() (int, int32)) (int, int32, string) {
	var n int
	status, message := WithLastError(func() int32 {
		var callStatus int32
		n, callStatus = call()
		return callStatus
	})
	return n, status, message
}
