package main

import (
	"fmt"
)

func main() {
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
