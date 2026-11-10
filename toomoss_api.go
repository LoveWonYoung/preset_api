package main

import (
	"errors"
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
	channel    byte
	masterMode byte
	baudrate   uint32
}

type LinFrame struct {
	id            byte
	dataLen       byte
	checksum      byte //1 classic 0 enhanced
	checksumKnown byte /* 0 unknown, 1 driver or channel configuration */
	data          [8]byte
	timestamp     uint32
	timestampUs   uint64
}

func ToomossDefaultCfg() ToomossLinCfg {
	var tlc ToomossLinCfg
	r, _, _ := syscall.SyscallN(preset_toomoss_lin_default_config)
	copyRegisterResult(unsafe.Pointer(&tlc), unsafe.Sizeof(tlc), r)
	return tlc
}

type Toomoss struct {
	device uintptr
}

func NewToomoss() (*Toomoss, error) {

	var device uintptr
	if r1, _, _ := syscall.SyscallN(preset_toomoss_open, uintptr(unsafe.Pointer(&device))); r1 != 0 {
		fmt.Println(r1)
		return nil, errors.New("open fail")
	} else {
		fmt.Printf("device open Ok (%d)\n", int32(r1))
	}

	return &Toomoss{device: device}, nil
}

func (t *Toomoss) LinInit(channels ...byte) error {
	if len(channels) == 0 {
		return fmt.Errorf("channals invalid")
	}

	for _, v := range channels {
		if v > 3 {
			return fmt.Errorf("channals invalid any channel > 3")
		}
		cfg := ToomossDefaultCfg()
		cfg.channel = v
		if r1, _, _ := syscall.SyscallN(
			preset_toomoss_lin_init,
			t.device,
			uintptr(unsafe.Pointer(&cfg))); r1 != 0 {
			fmt.Println(r1)
			return errors.New("init fail")
		} else {
			fmt.Printf("device init Ok (%d)\n", r1)
		}
	}
	return nil
}

func (t *Toomoss) Break(channel byte) error {
	if r1, _, _ := syscall.SyscallN(preset_toomoss_lin_break, t.device, uintptr(channel)); r1 != 0 {
		return errors.New("toomoss break is fail")
	}
	return nil
}

func (t *Toomoss) Send(channel byte, id byte, data []byte) error {
	if r1, _, _ := syscall.SyscallN(
		preset_lin_master_write,
		t.device,
		uintptr(channel),
		uintptr(id),
		uintptr(unsafe.Pointer(&data[0])),
		uintptr(len(data))); r1 != 0 {
		return errors.New("toomoss send is fail")
	}
	return nil
}

func (t *Toomoss) Rx(channel byte, id byte, timeoutMs uint32) (LinFrame, error) {
	var frame LinFrame
	if r1, _, _ := syscall.SyscallN(
		preset_lin_master_read,
		t.device,
		uintptr(channel),
		uintptr(id),
		uintptr(timeoutMs),
		uintptr(unsafe.Pointer(&frame))); r1 != 0 {
		return LinFrame{}, errors.New("toomoss rx is fail")
	}
	return frame, nil
}

func (t *Toomoss) RxOpt(channel byte, id byte, timeoutMs uint32) (o Option[LinFrame]) {
	var f, err = t.Rx(channel, id, timeoutMs)
	if err != nil {
		None[LinFrame]()
		return
	}
	return Some(f)
}

func (t *Toomoss) Close() error {
	if t == nil || t.device == 0 {
		return nil
	}
	if rc := syscallR1(preset_device_close, uintptr(unsafe.Pointer(&t.device))); int32(rc) < 0 {
		return &PresetError{Code: int32(rc), Func: "preset_device_close", Message: LastError()}
	}
	return nil
}
