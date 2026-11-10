//go:build windows && amd64

package main

import (
	"fmt"

	preset "github.com/LoveWonYoung/preset_api"
)

func main() {
	dev, err := preset.NewToomoss()
	if err != nil {
		panic(err)
	}
	defer dev.Close()

	// Default configuration: CAN-FD, 500 kbit/s / 2 Mbit/s, Toomoss CAN1.
	if err := dev.InitCAN(0); err != nil {
		panic(err)
	}

	runtime, err := dev.NewCANRuntime(0)
	if err != nil {
		panic(err)
	}
	defer func() { _ = runtime.Shutdown() }()

	config := preset.DefaultCANDiagnosticConfig(0x73A, 0x7BA, 0x7DF)
	timeouts := preset.DefaultRequestTimeouts()
	info, response, err := runtime.UDSRequest(
		&config,
		&timeouts,
		[]byte{0x19, 0x02, 0x09},
		35_000,
	)
	fmt.Printf("outcome=%d nrc=0x%02X response=% X err=%v\n", info.Outcome, info.LastNRC, response, err)

	load, err := runtime.GetBusLoad()
	if err != nil {
		panic(err)
	}
	fmt.Printf("bus load: %.1f%%\n", load.Load*100)
}
