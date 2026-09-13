package preset_api

import (
	"os"
	"testing"
	"unsafe"
)

func TestABILayouts(t *testing.T) {
	tests := []struct {
		name string
		got  uintptr
		want uintptr
	}{
		{"PresetConfig", unsafe.Sizeof(PresetConfig{}), 40},
		{"PresetToomossConfig", unsafe.Sizeof(PresetToomossConfig{}), 28},
		{"PresetToomossLinConfig", unsafe.Sizeof(PresetToomossLinConfig{}), 8},
		{"PresetToomossElinsConfig", unsafe.Sizeof(PresetToomossElinsConfig{}), 12},
		{"PresetToomossLinFrame", unsafe.Sizeof(PresetToomossLinFrame{}), 16},
		{"PresetToomossElinsMessage", unsafe.Sizeof(PresetToomossElinsMessage{}), 88},
		{"PresetPCANConfig", unsafe.Sizeof(PresetPCANConfig{}), 28},
		{"PresetTSMasterConfig", unsafe.Sizeof(PresetTSMasterConfig{}), 36},
		{"PresetVectorConfig", unsafe.Sizeof(PresetVectorConfig{}), 68},
		{"PresetCanFrame", unsafe.Sizeof(PresetCanFrame{}), 72},
		{"PresetRxStats", unsafe.Sizeof(PresetRxStats{}), 24},
		{"PresetBusLoad", unsafe.Sizeof(PresetBusLoad{}), 32},
	}
	for _, test := range tests {
		if test.got != test.want {
			t.Errorf("sizeof(%s) = %d, want %d", test.name, test.got, test.want)
		}
	}
}

func TestToomossChannelMasks(t *testing.T) {
	if PRESET_TOOMOSS_CHANNEL_1 != 1 || PRESET_TOOMOSS_CHANNEL_2 != 2 ||
		PRESET_TOOMOSS_CHANNEL_3 != 4 || PRESET_TOOMOSS_CHANNEL_4 != 8 {
		t.Fatal("Toomoss channel constants must be bit masks")
	}
}

func TestDLLSmoke(t *testing.T) {
	path := os.Getenv("PRESET_RS_DLL")
	if path == "" {
		t.Skip("set PRESET_RS_DLL to run the DLL integration test")
	}
	if err := LoadDLL(path); err != nil {
		t.Fatal(err)
	}
	if version := Version(); version == "" {
		t.Fatalf("Version returned empty string: %v", DLLLoadError())
	}
	if ABIVersion() == 0 {
		t.Fatal("ABIVersion returned zero")
	}
	if Capabilities()&PRESET_CAP_CLASSIC_CAN == 0 {
		t.Fatal("classic CAN capability is missing")
	}

	config := DefaultConfig()
	if config.PhysicalID != 0x7e0 || config.ResponseID != 0x7e8 || config.MaxPDULen == 0 {
		t.Fatalf("invalid default config: %+v", config)
	}
	lin := ToomossLinDefaultConfig()
	if lin.ChannelsMask != PRESET_TOOMOSS_CHANNEL_1 || lin.MasterMode != 1 || lin.Baudrate != 19_200 {
		t.Fatalf("invalid Toomoss LIN default config: %+v", lin)
	}
	toomoss := ToomossDefaultConfig()
	if toomoss.NominalBitrate != 500_000 || toomoss.DataBitrate != 2_000_000 {
		t.Fatalf("invalid Toomoss default config: %+v", toomoss)
	}
	elins := ToomossElinsDefaultConfig()
	if elins.ChannelsMask != PRESET_TOOMOSS_CHANNEL_1 || elins.ResistorEnabled != 1 ||
		elins.Version != PRESET_TOOMOSS_ELINS_VER_IND83220 || elins.ReceiveTimeoutUS != 1_000 {
		t.Fatalf("invalid Toomoss ELINS default config: %+v", elins)
	}
	if got := PcanDefaultConfig().NominalBitrate; got != 500_000 {
		t.Fatalf("PCAN nominal bitrate = %d, want 500000", got)
	}
	if config := TsmasterDefaultConfig(); config.DeviceType != PRESET_TSMASTER_TC1016 || config.NominalBitrate != 500_000 {
		t.Fatalf("invalid TSMaster default config: %+v", config)
	}
	if config := VectorDefaultConfig(); config.HardwareType != PRESET_VECTOR_HWTYPE_ANY ||
		config.DriverRxQueueSize != 16_384 || config.NominalBitrate != 500_000 {
		t.Fatalf("invalid Vector default config: %+v", config)
	}
	if status := CanUdsSetDefaultSTMin(0, 0); status != PRESET_ERR_NOT_INIT {
		t.Fatalf("null CAN client status = %d, want %d", status, PRESET_ERR_NOT_INIT)
	}
	if _, status := CanUdsRequest(0, []byte{0x10, 0x01}, 100, make([]byte, 8)); status != PRESET_ERR_NOT_INIT {
		t.Fatalf("null CAN request status = %d, want %d", status, PRESET_ERR_NOT_INIT)
	}
	if _, _, status := LinUdsRequest(0, []byte{0x10, 0x01}, 100, make([]byte, 8)); status != PRESET_ERR_NOT_INIT {
		t.Fatalf("null LIN request status = %d, want %d", status, PRESET_ERR_NOT_INIT)
	}
	var device PresetDevice
	if status := DeviceClose(&device); status != PRESET_OK || device != 0 {
		t.Fatalf("closing null device returned status %d and handle %#x", status, device)
	}
	var canClient PresetCanUdsClient
	if status := CanUdsClientClose(&canClient); status != PRESET_OK || canClient != 0 {
		t.Fatalf("closing null CAN client returned status %d and handle %#x", status, canClient)
	}
	var linClient PresetLinUdsClient
	if status := LinUdsClientClose(&linClient); status != PRESET_OK || linClient != 0 {
		t.Fatalf("closing null LIN client returned status %d and handle %#x", status, linClient)
	}
}
