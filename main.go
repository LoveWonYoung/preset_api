package main

import (
	"fmt"
	"time"
)

const (
	CHANNEL1 byte = iota
	CHANNEL2
	CHANNEL3
	CHANNEL4
)

func main() {
	if false {
		t, err := NewToomoss()
		if err != nil {
			return
		}
		defer t.Close()
		t.LinInit(0)
		if t.Break(CHANNEL1) != nil {
			panic("break is panic")
		}
		time.Sleep(400 * time.Millisecond)
		if l := t.RxOpt(0, 1, 10); l.IsSome() {
			fmt.Println(l.Unwarp())
		} else {
			fmt.Println("is none")
		}
		fmt.Printf("t.Send(0, 0x3c, []byte{0x7F, 0x03, 0x22, 0xF1, 0x89, 0xFF, 0xFF, 0xFF}): %v\n", t.Send(0, 0x3c, []byte{0x7F, 0x03, 0x22, 0xF1, 0x89, 0xFF, 0xFF, 0xFF}))
		fmt.Printf("t.Send(0, 0x3c, []byte{0x7F, 0x03, 0x22, 0xF1, 0x89, 0xFF, 0xFF, 0xFF}): %v\n", t.Send(0, 0x3c, []byte{0x7F, 0x03, 0x22, 0xF1, 0x89, 0xFF, 0xFF, 0xFF}))
		fmt.Printf("t.Send(0, 0x3c, []byte{0x7F, 0x03, 0x22, 0xF1, 0x89, 0xFF, 0xFF, 0xFF}): %v\n", t.Send(0, 0x3c, []byte{0x7F, 0x03, 0x22, 0xF1, 0x89, 0xFF, 0xFF, 0xFF}))
		for range 10 {
			// fmt.Println(t.Rx(0, 0x3d, 20))
			if l := t.RxOpt(0, 0x3d, 20); l.IsSome() {
				fmt.Println(l.Unwarp())
			} else {
				fmt.Println("is none")
			}
		}
	}

	dev, err := NewToomoss()
	if err != nil {
		panic(err)
	}
	defer dev.Close()

	// 默认配置：CAN-FD、500k/2M，Toomoss CAN1 = 0
	if err := dev.CanInit(0); err != nil {
		panic(err)
	}

	// 原始发送 / 接收
	// _ = dev.CanWrite(0, 0x401, make([]byte, 8), false, false)
	// frames, _ := dev.CanTryRead(0, 256)

	// UDS
	rt, err := dev.CanRuntimeNew(0) // command_capacity=256, recovery=100ms
	if err != nil {
		panic(err)
	}
	defer func() { _ = rt.Shutdown() }()

	cfg := DefaultCanDiagnosticConfig(0x73A, 0x7BA, 0x7DF)
	timeouts := DefaultRequestTimeouts()

	info, resp, err := rt.UdsRequest(&cfg, &timeouts, []byte{0x19, 0x02, 0x09}, 35000)
	fmt.Printf("outcome=%d nrc=0x%02X resp=% X err=%v\n", info.Outcome, info.LastNrc, resp, err)

	load, _ := rt.GetBusLoad()
	fmt.Printf("bus load: %.1f%%\n", load.Load*100)
}
