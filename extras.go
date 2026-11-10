//go:build windows && amd64

package preset

import (
	"fmt"
	"strings"
	"unsafe"
)

const (
	PresetABIVersion     uint32 = 12
	PresetCapTSMaster    uint64 = 1 << 9
	PresetCapPCAN        uint64 = 1 << 10
	PresetCapVector      uint64 = 1 << 11
	PresetCapAutoDriver  uint64 = 1 << 20
	PresetCapPCANLIN     uint64 = 1 << 21
	PresetCapTSMasterLIN uint64 = 1 << 22
	PresetCapVectorLIN   uint64 = 1 << 23
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
	if strings.IndexByte(name, 0) >= 0 {
		return fmt.Errorf("log name contains NUL")
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
	if mode != PresetLogFilterOff && mode != PresetLogFilterList {
		return fmt.Errorf("invalid log filter mode %d", mode)
	}
	if len(ids) > 2048 {
		return fmt.Errorf("log filter supports at most 2048 IDs")
	}
	for _, id := range ids {
		if id > 0x7FF {
			return fmt.Errorf("log filter ID 0x%X is not a standard 11-bit CAN ID", id)
		}
	}
	var ptr uintptr
	if len(ids) > 0 {
		ptr = uintptr(unsafe.Pointer(&ids[0]))
	}
	return nativeStatus("preset_set_log_filter", uintptr(mode), ptr, uintptr(len(ids)))
}

type CANPeriodic struct{ handle nativeHandle }

func (r *CANRuntime) StartPeriodic(frames []CANFrame, periodMs uint32) (*CANPeriodic, error) {
	if len(frames) < 1 || len(frames) > 64 || periodMs == 0 {
		return nil, fmt.Errorf("periodic requires 1-64 frames and a positive period")
	}
	for i := range frames {
		if err := frames[i].Validate(); err != nil {
			return nil, fmt.Errorf("periodic frame %d: %w", i, err)
		}
	}
	return useHandle(&r.handle, func(runtime uintptr) (*CANPeriodic, error) {
		var handle uintptr
		if err := nativeStatus("preset_can_periodic_start", runtime, uintptr(unsafe.Pointer(&frames[0])), uintptr(len(frames)), uintptr(periodMs), uintptr(unsafe.Pointer(&handle))); err != nil {
			return nil, err
		}
		return &CANPeriodic{handle: newNativeHandle(handle)}, nil
	})
}

func (p *CANPeriodic) Cancel() error {
	if p == nil {
		return nil
	}
	return p.handle.use(func(handle uintptr) error {
		return nativeStatus("preset_can_periodic_cancel", handle)
	})
}

func (p *CANPeriodic) Close() error {
	if p == nil {
		return nil
	}
	return p.handle.close("preset_can_periodic_close")
}
