package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"

	preset "github.com/LoveWonYoung/preset_api"
)

const (
	physicalID   uint32 = 0x700
	responseID   uint32 = 0x708
	functionalID uint32 = 0x7DF
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	dllPath := os.Getenv("PRESET_RS_DLL")
	if dllPath == "" {
		dllPath = "preset_rs.dll"
	}

	device, err := preset.NewMyDevice(
		preset.BackendToomoss,
		physicalID,
		responseID,
		functionalID,
		[]uint8{0},
		0,
	)
	if err != nil {
		return err
	}
	device.DLLPath = dllPath
	if err := device.Open(); err != nil {
		return err
	}
	defer device.Close()

	vin, err := device.Request([]byte{0x22, 0xF1, 0x80}, 1000)
	if err != nil {
		fmt.Println("UDS request:", err)
	} else {
		fmt.Println("UDS 22 F1 80 ->", spacedHex(vin))
	}

	frames, err := device.Rxfn(2000)
	if err != nil {
		return err
	}
	for _, frame := range frames {
		fmt.Printf(
			"%s id=0x%X dlc=%d fd=%d brs=%d data=%s\n",
			frame.Direction,
			frame.ID,
			frame.DLC,
			boolBit(frame.IsFD),
			boolBit(frame.BRS),
			spacedHex(frame.Data),
		)
	}
	return nil
}

func spacedHex(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	encoded := strings.ToUpper(hex.EncodeToString(data))
	parts := make([]string, 0, len(data))
	for i := 0; i < len(encoded); i += 2 {
		parts = append(parts, encoded[i:i+2])
	}
	return strings.Join(parts, " ")
}

func boolBit(value bool) int {
	if value {
		return 1
	}
	return 0
}
