package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

// All backends share the Device CAN/LIN methods. Close child runtimes and LIN
// clients before closing the device. Native DLL and vendor drivers are Windows-only.
type PCAN struct{ Device }
type TSMaster struct{ Device }
type Vector struct{ Device }
type AutoDevice struct{ Device }

type AutoCfg struct {
	Channel        byte
	Can            CanChannelConfig
	CandidateOrder [4]byte
}

type PCANCfg struct {
	Channel byte
	Can     CanChannelConfig
}

type TSMasterCfg struct {
	ApplicationChannel byte
	HardwareChannel    byte
	HardwareIndex      int32
	DeviceType         int32
	Can                CanChannelConfig
}

type VectorCfg struct {
	ApplicationChannel byte
	HardwareType       int32
	HardwareIndex      uint32
	HardwareChannel    uint32
	NominalSjw         uint32
	NominalTseg1       uint32
	NominalTseg2       uint32
	DataSjw            uint32
	DataTseg1          uint32
	DataTseg2          uint32
	DriverRxQueueSize  uint32
	Can                CanChannelConfig
}

func proc(name string) (uintptr, error) {
	if lib == 0 {
		return 0, fmt.Errorf("load bin\\preset_rs.dll before calling %s", name)
	}
	p, err := syscall.GetProcAddress(lib, name)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", name, err)
	}
	return p, nil
}

func nativeStatus(name string, args ...uintptr) error {
	p, err := proc(name)
	if err != nil {
		return err
	}
	return presetStatus(name, syscallR1(p, args...))
}

func openDevice(p uintptr, name string) (uintptr, error) {
	if p == 0 {
		return 0, fmt.Errorf("%s unavailable: check bin\\preset_rs.dll", name)
	}
	if err := checkABI(); err != nil {
		return 0, err
	}
	var handle uintptr
	if err := presetStatus(name, syscallR1(p, uintptr(unsafe.Pointer(&handle)))); err != nil {
		return 0, err
	}
	return handle, nil
}

func openNamedDevice(name string) (uintptr, error) {
	p, err := proc(name)
	if err != nil {
		return 0, err
	}
	return openDevice(p, name)
}

// Native struct-return helpers use a hidden output pointer on Windows amd64.
func defaultStruct[T any](name string) (T, error) {
	var result T
	p, err := proc(name)
	if err != nil {
		return result, err
	}
	if unsafe.Sizeof(result) > 8 {
		syscallR1(p, uintptr(unsafe.Pointer(&result)))
	} else {
		copyRegisterResult(unsafe.Pointer(&result), unsafe.Sizeof(result), syscallR1(p))
	}
	return result, nil
}

func NewPCAN() (*PCAN, error) {
	h, err := openNamedDevice("preset_pcan_open")
	if err != nil {
		return nil, err
	}
	return &PCAN{Device{device: h}}, nil
}

func NewTSMaster() (*TSMaster, error) {
	h, err := openNamedDevice("preset_tsmaster_open")
	if err != nil {
		return nil, err
	}
	return &TSMaster{Device{device: h}}, nil
}

func NewVector() (*Vector, error) {
	h, err := openNamedDevice("preset_vector_open")
	if err != nil {
		return nil, err
	}
	return &Vector{Device{device: h}}, nil
}

func AutoDefaultCfg() (AutoCfg, error) { return defaultStruct[AutoCfg]("preset_auto_default_config") }
func PCANDefaultCfg() (PCANCfg, error) { return defaultStruct[PCANCfg]("preset_pcan_default_config") }
func TSMasterDefaultCfg() (TSMasterCfg, error) {
	return defaultStruct[TSMasterCfg]("preset_tsmaster_default_config")
}
func VectorDefaultCfg() (VectorCfg, error) {
	return defaultStruct[VectorCfg]("preset_vector_default_config")
}
func PCANDefaultFdTiming() (CanFdTiming, error) {
	return defaultStruct[CanFdTiming]("preset_pcan_default_fd_timing")
}

// NewAuto selects an installed backend in the configured priority order.
func NewAuto(cfg *AutoCfg) (*AutoDevice, error) {
	if cfg == nil {
		return nil, fmt.Errorf("auto config is nil")
	}
	if err := checkABI(); err != nil {
		return nil, err
	}
	var h uintptr
	if err := nativeStatus("preset_auto_open", uintptr(unsafe.Pointer(cfg)), uintptr(unsafe.Pointer(&h))); err != nil {
		return nil, err
	}
	return &AutoDevice{Device{device: h}}, nil
}

func (t *PCAN) CanInitCfg(cfg *PCANCfg, timing *CanFdTiming) error {
	if cfg == nil {
		return fmt.Errorf("PCAN config is nil")
	}
	return nativeStatus("preset_pcan_can_init", t.device, uintptr(unsafe.Pointer(cfg)), uintptr(unsafe.Pointer(timing)))
}

func (t *TSMaster) CanInitCfg(cfg *TSMasterCfg) error {
	if cfg == nil {
		return fmt.Errorf("TSMaster config is nil")
	}
	return nativeStatus("preset_tsmaster_can_init", t.device, uintptr(unsafe.Pointer(cfg)))
}

func (t *Vector) CanInitCfg(cfg *VectorCfg) error {
	if cfg == nil {
		return fmt.Errorf("Vector config is nil")
	}
	return nativeStatus("preset_vector_can_init", t.device, uintptr(unsafe.Pointer(cfg)))
}

func (t *PCAN) CanInit(channels ...byte) error {
	for _, ch := range channels {
		cfg, err := PCANDefaultCfg()
		if err != nil {
			return err
		}
		cfg.Channel = ch
		if err := t.CanInitCfg(&cfg, nil); err != nil {
			return err
		}
	}
	return nil
}

func (t *TSMaster) CanInit(channels ...byte) error {
	for _, ch := range channels {
		cfg, err := TSMasterDefaultCfg()
		if err != nil {
			return err
		}
		cfg.ApplicationChannel, cfg.HardwareChannel = ch, ch
		if err := t.CanInitCfg(&cfg); err != nil {
			return err
		}
	}
	return nil
}

func (t *Vector) CanInit(channels ...byte) error {
	for _, ch := range channels {
		cfg, err := VectorDefaultCfg()
		if err != nil {
			return err
		}
		cfg.ApplicationChannel, cfg.HardwareChannel = ch, uint32(ch)
		if err := t.CanInitCfg(&cfg); err != nil {
			return err
		}
	}
	return nil
}
