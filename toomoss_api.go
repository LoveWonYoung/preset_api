package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	lib, _                               = syscall.LoadLibrary("bin\\preset_rs.dll")
	preset_toomoss_lin_default_config, _ = syscall.GetProcAddress(lib, "preset_toomoss_lin_default_config")
	preset_toomoss_open, _               = syscall.GetProcAddress(lib, "preset_toomoss_open")
	preset_toomoss_lin_init, _           = syscall.GetProcAddress(lib, "preset_toomoss_lin_init")
	preset_toomoss_lin_break, _          = syscall.GetProcAddress(lib, "preset_toomoss_lin_break")
	preset_lin_master_write, _           = syscall.GetProcAddress(lib, "preset_lin_master_write")
	preset_device_close, _               = syscall.GetProcAddress(lib, "preset_device_close")
	preset_lin_master_read, _            = syscall.GetProcAddress(lib, "preset_lin_master_read")
	preset_lin_request_start, _          = syscall.GetProcAddress(lib, "preset_lin_request_start")
)

type ToomossLinCfg struct {
	Channel    byte
	MasterMode byte
	Baudrate   uint32
}

type LinFrame struct {
	ID              byte
	DataLen         byte
	ClassicChecksum byte // 1 classic, 0 enhanced
	ChecksumKnown   byte // 0 unknown, 1 driver or channel configuration
	Data            [8]byte
	Timestamp       uint32
	TimestampUs     uint64
}

func ToomossDefaultCfg() ToomossLinCfg {
	var tlc ToomossLinCfg
	r, _, _ := syscall.SyscallN(preset_toomoss_lin_default_config)
	copyRegisterResult(unsafe.Pointer(&tlc), unsafe.Sizeof(tlc), r)
	return tlc
}

type Toomoss struct {
	Device
}

// Device owns a preset_rs handle shared by CAN and LIN channels.
type Device struct {
	device uintptr
}

func NewToomoss() (*Toomoss, error) {

	device, err := openDevice(preset_toomoss_open, "preset_toomoss_open")
	if err != nil {
		return nil, err
	}
	return &Toomoss{Device: Device{device: device}}, nil
}

func (t *Toomoss) LinInit(channels ...byte) error {
	if len(channels) == 0 {
		return fmt.Errorf("channels invalid")
	}

	for _, v := range channels {
		if v > 3 {
			return fmt.Errorf("channels invalid: channel > 3")
		}
		cfg := ToomossDefaultCfg()
		cfg.Channel = v
		if err := presetStatus("preset_toomoss_lin_init", syscallR1(
			preset_toomoss_lin_init,
			t.device,
			uintptr(unsafe.Pointer(&cfg)))); err != nil {
			return err
		}
	}
	return nil
}

func (t *Toomoss) Break(channel byte) error {
	return presetStatus("preset_toomoss_lin_break", syscallR1(preset_toomoss_lin_break, t.device, uintptr(channel)))
}

func (t *Toomoss) Send(channel byte, id byte, data []byte) error {
	if len(data) == 0 || len(data) > 8 {
		return fmt.Errorf("LIN payload must contain 1-8 bytes")
	}
	return presetStatus("preset_lin_master_write", syscallR1(
		preset_lin_master_write,
		t.device,
		uintptr(channel),
		uintptr(id),
		bytePtr(data),
		uintptr(len(data))))
}

func (t *Toomoss) Rx(channel byte, id byte, timeoutMs uint32) (LinFrame, error) {
	var frame LinFrame
	if err := presetStatus("preset_lin_master_read", syscallR1(
		preset_lin_master_read,
		t.device,
		uintptr(channel),
		uintptr(id),
		uintptr(timeoutMs),
		uintptr(unsafe.Pointer(&frame)))); err != nil {
		return LinFrame{}, err
	}
	return frame, nil
}

func (t *Toomoss) RxOpt(channel byte, id byte, timeoutMs uint32) (o Option[LinFrame]) {
	var f, err = t.Rx(channel, id, timeoutMs)
	if err != nil {
		return None[LinFrame]()
	}
	return Some(f)
}

func (t *Device) Close() error {
	if t == nil || t.device == 0 {
		return nil
	}
	return presetStatus("preset_device_close", syscallR1(preset_device_close, uintptr(unsafe.Pointer(&t.device))))
}
