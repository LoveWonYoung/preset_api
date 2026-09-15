//go:build windows && amd64

package preset_api

import "testing"

func TestNewMyDeviceRequiresChannels(t *testing.T) {
	_, err := NewMyDevice(BackendToomoss, 0x73A, 0x7BA, 0x7DF, nil, 0)
	if err == nil {
		t.Fatal("NewMyDevice must reject empty channels")
	}
}

func TestMyDeviceMethodsRequireOpen(t *testing.T) {
	device, err := NewMyDevice(BackendToomoss, 0x73A, 0x7BA, 0x7DF, []uint8{0}, 0)
	if err != nil {
		t.Fatal(err)
	}
	if err := device.Txfn(0x73A, []byte{0x02, 0x10, 0x01}); err == nil {
		t.Fatal("Txfn must fail before Open")
	}
	if _, err := device.Rxfn(0); err == nil {
		t.Fatal("Rxfn must fail before Open")
	}
	if _, err := device.Request([]byte{0x22, 0xF1, 0x80}, 1000); err == nil {
		t.Fatal("Request must fail before Open")
	}
	if err := device.Close(); err != nil {
		t.Fatalf("Close on a never-opened device: %v", err)
	}
}

func TestBackendString(t *testing.T) {
	if BackendToomoss.String() != "Toomoss" || BackendPcan.String() != "PCAN" {
		t.Fatalf("unexpected backend names: %s %s", BackendToomoss, BackendPcan)
	}
	if BackendToomoss != PRESET_CAN_BACKEND_TOOMOSS || BackendVector != PRESET_CAN_BACKEND_VECTOR {
		t.Fatal("Backend values must match PRESET_CAN_BACKEND_*")
	}
}
