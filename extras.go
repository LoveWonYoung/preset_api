package main

import (
	"fmt"
	"unsafe"
)

const (
	PresetABIVersion     uint32 = 12
	PresetCapTSMaster    uint64 = 1 << 9
	PresetCapPCAN        uint64 = 1 << 10
	PresetCapVector      uint64 = 1 << 11
	PresetCapAuto        uint64 = 1 << 20
	PresetCapPCANLin     uint64 = 1 << 21
	PresetCapTSMasterLin uint64 = 1 << 22
	PresetCapVectorLin   uint64 = 1 << 23
	PresetLogFilterOff   byte   = 0
	PresetLogFilterList  byte   = 1
)

func checkABI() error {
	version, err := ABIVersion()
	if err != nil {
		return err
	}
	if version != PresetABIVersion {
		return fmt.Errorf("preset_rs ABI version %d, expected %d", version, PresetABIVersion)
	}
	return nil
}

func ABIVersion() (uint32, error) {
	p, err := proc("preset_abi_version")
	if err != nil {
		return 0, err
	}
	return uint32(syscallR1(p)), nil
}

func Capabilities() (uint64, error) {
	p, err := proc("preset_get_capabilities")
	if err != nil {
		return 0, err
	}
	return uint64(syscallR1(p)), nil
}

// LogInit starts process-wide asynchronous logging, initially disabled.
func LogInit(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("log name is empty")
	}
	data := append([]byte(name), 0)
	return nativeStatus("preset_log_init", bytePtr(data))
}

func LogShutdown() error {
	p, err := proc("preset_log_shutdown")
	if err == nil {
		syscallR1(p)
	}
	return err
}

func SetPrintLog(enable bool) error {
	p, err := proc("preset_set_print_log")
	if err == nil {
		syscallR1(p, uintptr(boolByte(enable)))
	}
	return err
}

func LogDroppedCount() (uint64, error) {
	p, err := proc("preset_log_dropped_count")
	if err != nil {
		return 0, err
	}
	return uint64(syscallR1(p)), nil
}

func SetLogFilter(mode byte, ids []uint32) error {
	if len(ids) > 2048 {
		return fmt.Errorf("log filter supports at most 2048 IDs")
	}
	var ptr uintptr
	if len(ids) > 0 {
		ptr = uintptr(unsafe.Pointer(&ids[0]))
	}
	return nativeStatus("preset_set_log_filter", uintptr(mode), ptr, uintptr(len(ids)))
}

type CanPeriodic struct{ handle uintptr }

func (r *CanRuntime) PeriodicStart(frames []CanFrame, periodMs uint32) (*CanPeriodic, error) {
	if len(frames) < 1 || len(frames) > 64 || periodMs == 0 {
		return nil, fmt.Errorf("periodic requires 1-64 frames and a positive period")
	}
	var handle uintptr
	err := nativeStatus("preset_can_periodic_start", r.handle, uintptr(unsafe.Pointer(&frames[0])), uintptr(len(frames)), uintptr(periodMs), uintptr(unsafe.Pointer(&handle)))
	if err != nil {
		return nil, err
	}
	return &CanPeriodic{handle: handle}, nil
}

func (p *CanPeriodic) Cancel() error {
	if p == nil || p.handle == 0 {
		return nil
	}
	return nativeStatus("preset_can_periodic_cancel", p.handle)
}

func (p *CanPeriodic) Close() error {
	if p == nil || p.handle == 0 {
		return nil
	}
	return nativeStatus("preset_can_periodic_close", uintptr(unsafe.Pointer(&p.handle)))
}
