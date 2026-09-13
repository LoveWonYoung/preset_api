//go:build windows && amd64

package preset_api

import (
	"os"
	"runtime"
	"testing"
	"unsafe"
)

func TestABILayouts(t *testing.T) {
	if runtime.GOOS != "windows" || runtime.GOARCH != "amd64" || unsafe.Sizeof(uintptr(0)) != 8 {
		t.Fatalf("unsupported target %s/%s", runtime.GOOS, runtime.GOARCH)
	}
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

	offsets := []struct {
		name string
		got  uintptr
		want uintptr
	}{
		{"PresetConfig.IsFD", unsafe.Offsetof(PresetConfig{}.IsFD), 32},
		{"PresetConfig.Reserved", unsafe.Offsetof(PresetConfig{}.Reserved), 38},
		{"PresetToomossConfig.NominalBitrate", unsafe.Offsetof(PresetToomossConfig{}.NominalBitrate), 4},
		{"PresetToomossConfig.MinRxPollIntervalUS", unsafe.Offsetof(PresetToomossConfig{}.MinRxPollIntervalUS), 24},
		{"PresetToomossLinConfig.Baudrate", unsafe.Offsetof(PresetToomossLinConfig{}.Baudrate), 4},
		{"PresetToomossElinsConfig.ReceiveTimeoutUS", unsafe.Offsetof(PresetToomossElinsConfig{}.ReceiveTimeoutUS), 8},
		{"PresetToomossLinFrame.Data", unsafe.Offsetof(PresetToomossLinFrame{}.Data), 8},
		{"PresetToomossElinsMessage.MsgSendTimes", unsafe.Offsetof(PresetToomossElinsMessage{}.MsgSendTimes), 6},
		{"PresetToomossElinsMessage.Timestamp", unsafe.Offsetof(PresetToomossElinsMessage{}.Timestamp), 8},
		{"PresetToomossElinsMessage.Data", unsafe.Offsetof(PresetToomossElinsMessage{}.Data), 18},
		{"PresetToomossElinsMessage.AckValue", unsafe.Offsetof(PresetToomossElinsMessage{}.AckValue), 82},
		{"PresetPCANConfig.NominalBitrate", unsafe.Offsetof(PresetPCANConfig{}.NominalBitrate), 4},
		{"PresetTSMasterConfig.HardwareIndex", unsafe.Offsetof(PresetTSMasterConfig{}.HardwareIndex), 4},
		{"PresetTSMasterConfig.MinRxPollIntervalUS", unsafe.Offsetof(PresetTSMasterConfig{}.MinRxPollIntervalUS), 32},
		{"PresetVectorConfig.HardwareType", unsafe.Offsetof(PresetVectorConfig{}.HardwareType), 4},
		{"PresetVectorConfig.DriverRxQueueSize", unsafe.Offsetof(PresetVectorConfig{}.DriverRxQueueSize), 48},
		{"PresetVectorConfig.MinRxPollIntervalUS", unsafe.Offsetof(PresetVectorConfig{}.MinRxPollIntervalUS), 64},
		{"PresetCanFrame.Data", unsafe.Offsetof(PresetCanFrame{}.Data), 8},
		{"PresetRxStats.Queued", unsafe.Offsetof(PresetRxStats{}.Queued), 16},
		{"PresetBusLoad.FrameCount", unsafe.Offsetof(PresetBusLoad{}.FrameCount), 24},
	}
	for _, test := range offsets {
		if test.got != test.want {
			t.Errorf("offsetof(%s) = %d, want %d", test.name, test.got, test.want)
		}
	}
}

func TestToomossChannelMasks(t *testing.T) {
	if PRESET_TOOMOSS_CHANNEL_1 != 1 || PRESET_TOOMOSS_CHANNEL_2 != 2 ||
		PRESET_TOOMOSS_CHANNEL_3 != 4 || PRESET_TOOMOSS_CHANNEL_4 != 8 {
		t.Fatal("Toomoss channel constants must be bit masks")
	}
}

func TestABIConstants(t *testing.T) {
	if SupportedABIVersion != 4 {
		t.Fatalf("supported ABI = %d, want 4", SupportedABIVersion)
	}
	statuses := [...]int32{
		PRESET_OK,
		PRESET_ERR_NULL_PTR,
		PRESET_ERR_BUFFER_TOO_SMALL,
		PRESET_ERR_NOT_INIT,
		PRESET_ERR_TIMEOUT,
		PRESET_ERR_UDS_NEGATIVE,
		PRESET_ERR_TRANSPORT,
		PRESET_ERR_UNSUPPORTED,
		PRESET_ERR_INVALID_ARG,
		PRESET_ERR_BUSY,
		PRESET_ERR_ISOTP,
		PRESET_ERR_PANIC,
	}
	for index, status := range statuses {
		if want := -int32(index); status != want {
			t.Fatalf("status %d = %d, want %d", index, status, want)
		}
	}
	if PRESET_CAP_CLASSIC_CAN != 1 || PRESET_CAP_BUS_LOAD != 1<<16 {
		t.Fatal("capability bit values do not match ABI v4")
	}
}

func TestHandleCopiesShareCloseState(t *testing.T) {
	state := newHandle(0x1234)
	first := PresetDevice{state: state}
	second := first
	if first.IsClosed() || second.IsClosed() {
		t.Fatal("non-zero handles must be open")
	}
	state.Lock()
	state.value = 0
	state.Unlock()
	if !first.IsClosed() || !second.IsClosed() {
		t.Fatal("all handle copies must observe close")
	}
}

func TestLocalValidationWithoutDLL(t *testing.T) {
	if err := LoadDLL(""); err == nil {
		t.Fatal("LoadDLL must reject an empty path")
	}
	status, message := WithLastError(nil)
	if status != PRESET_ERR_INVALID_ARG || message == "" {
		t.Fatalf("WithLastError(nil) = (%d, %q)", status, message)
	}
	var device PresetDevice
	if status := DeviceClose(&device); status != PRESET_OK {
		t.Fatalf("closing a zero-value device = %d, want %d", status, PRESET_OK)
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
	version, err := Version()
	if err != nil || version == "" {
		t.Fatalf("Version failed: version=%q error=%v", version, err)
	}
	abiVersion, err := ABIVersion()
	if err != nil || abiVersion == 0 {
		t.Fatal("ABIVersion returned zero")
	}
	capabilities, err := Capabilities()
	if err != nil || capabilities&PRESET_CAP_CLASSIC_CAN == 0 {
		t.Fatal("classic CAN capability is missing")
	}

	config := mustValue(t, DefaultConfig)
	if config.PhysicalID != 0x7e0 || config.ResponseID != 0x7e8 || config.MaxPDULen == 0 {
		t.Fatalf("invalid default config: %+v", config)
	}
	lin := mustValue(t, ToomossLinDefaultConfig)
	if lin.ChannelsMask != PRESET_TOOMOSS_CHANNEL_1 || lin.MasterMode != 1 || lin.Baudrate != 19_200 {
		t.Fatalf("invalid Toomoss LIN default config: %+v", lin)
	}
	toomoss := mustValue(t, ToomossDefaultConfig)
	if toomoss.NominalBitrate != 500_000 || toomoss.DataBitrate != 2_000_000 {
		t.Fatalf("invalid Toomoss default config: %+v", toomoss)
	}
	elins := mustValue(t, ToomossElinsDefaultConfig)
	if elins.ChannelsMask != PRESET_TOOMOSS_CHANNEL_1 || elins.ResistorEnabled != 1 ||
		elins.Version != PRESET_TOOMOSS_ELINS_VER_IND83220 || elins.ReceiveTimeoutUS != 1_000 {
		t.Fatalf("invalid Toomoss ELINS default config: %+v", elins)
	}
	if got := mustValue(t, PcanDefaultConfig).NominalBitrate; got != 500_000 {
		t.Fatalf("PCAN nominal bitrate = %d, want 500000", got)
	}
	if config := mustValue(t, TsmasterDefaultConfig); config.DeviceType != PRESET_TSMASTER_TC1016 || config.NominalBitrate != 500_000 {
		t.Fatalf("invalid TSMaster default config: %+v", config)
	}
	if config := mustValue(t, VectorDefaultConfig); config.HardwareType != PRESET_VECTOR_HWTYPE_ANY ||
		config.DriverRxQueueSize != 16_384 || config.NominalBitrate != 500_000 {
		t.Fatalf("invalid Vector default config: %+v", config)
	}
	if status := CanUdsSetDefaultSTMin(PresetCanUdsClient{}, 0); status != PRESET_ERR_NOT_INIT {
		t.Fatalf("null CAN client status = %d, want %d", status, PRESET_ERR_NOT_INIT)
	}
	if _, status := CanUdsRequest(PresetCanUdsClient{}, []byte{0x10, 0x01}, 100, make([]byte, 8)); status != PRESET_ERR_NOT_INIT {
		t.Fatalf("null CAN request status = %d, want %d", status, PRESET_ERR_NOT_INIT)
	}
	if _, _, status := LinUdsRequest(PresetLinUdsClient{}, []byte{0x10, 0x01}, 100, make([]byte, 8)); status != PRESET_ERR_NOT_INIT {
		t.Fatalf("null LIN request status = %d, want %d", status, PRESET_ERR_NOT_INIT)
	}
	var device PresetDevice
	if status := DeviceClose(&device); status != PRESET_OK || !device.IsClosed() {
		t.Fatalf("closing null device returned status %d", status)
	}
	var canClient PresetCanUdsClient
	if status := CanUdsClientClose(&canClient); status != PRESET_OK || !canClient.IsClosed() {
		t.Fatalf("closing null CAN client returned status %d", status)
	}
	var linClient PresetLinUdsClient
	if status := LinUdsClientClose(&linClient); status != PRESET_OK || !linClient.IsClosed() {
		t.Fatalf("closing null LIN client returned status %d", status)
	}
}

func mustValue[T any](t *testing.T, load func() (T, error)) T {
	t.Helper()
	value, err := load()
	if err != nil {
		t.Fatal(err)
	}
	return value
}
