//go:build windows && amd64

package preset

import (
	"errors"
	"testing"
)

func TestCANFrameValidation(t *testing.T) {
	tests := []struct {
		name  string
		frame CANFrame
		ok    bool
	}{
		{
			name: "classic data",
			frame: CANFrame{
				ID:        0x7FF,
				DataLen:   8,
				Direction: PresetCANDirectionTX,
			},
			ok: true,
		},
		{
			name: "extended FD data",
			frame: CANFrame{
				ID:         0x1FFFFFFF,
				DataLen:    64,
				Direction:  PresetCANDirectionTX,
				IsFD:       1,
				IsExtended: 1,
			},
			ok: true,
		},
		{
			name: "oversized standard ID",
			frame: CANFrame{
				ID:        0x800,
				Direction: PresetCANDirectionTX,
			},
		},
		{
			name: "FD remote",
			frame: CANFrame{
				ID:        1,
				Direction: PresetCANDirectionTX,
				IsFD:      1,
				IsRemote:  1,
			},
		},
		{
			name: "TX ESI",
			frame: CANFrame{
				ID:        1,
				Direction: PresetCANDirectionTX,
				ESI:       1,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.frame.Validate()
			if (err == nil) != test.ok {
				t.Fatalf("Validate() error = %v, want success %v", err, test.ok)
			}
		})
	}
}

func TestCANFrameConstructors(t *testing.T) {
	data := make([]byte, 64)
	frame, err := NewCANFrame(0x18DAF110, data, true, true, true)
	if err != nil {
		t.Fatal(err)
	}
	if frame.DataLen != 64 || frame.IsFD != 1 || frame.BRS != 1 || frame.IsExtended != 1 {
		t.Fatalf("unexpected frame: %+v", frame)
	}
	if _, err := NewCANFrame(0x123, make([]byte, 9), false, false, false); err == nil {
		t.Fatal("classic payload larger than 8 bytes was accepted")
	}
	if _, err := NewCANRemoteFrame(0x123, 9, false); err == nil {
		t.Fatal("remote DLC larger than 8 was accepted")
	}
}

func TestValidationDoesNotRequireHardware(t *testing.T) {
	device := &Device{}
	if _, err := device.TryReadCAN(0, 0); err == nil {
		t.Fatal("zero CAN read capacity was accepted")
	}
	if err := device.WriteCAN(0, 0x800, nil, false, false); err == nil {
		t.Fatal("oversized standard CAN ID was accepted")
	}
	if _, err := (&Toomoss{}).TryReadLIN(0, 513); err == nil {
		t.Fatal("oversized LIN read capacity was accepted")
	}
	if _, err := (&CANRuntime{}).StartRequest(nil, nil, nil); err == nil {
		t.Fatal("missing diagnostic request arguments were accepted")
	}
	if err := SetLogFilter(99, nil); err == nil {
		t.Fatal("invalid log filter mode was accepted")
	}
}

func TestClosedHandle(t *testing.T) {
	var handle nativeHandle
	if err := handle.use(func(uintptr) error { return nil }); !errors.Is(err, ErrClosed) {
		t.Fatalf("use() error = %v, want ErrClosed", err)
	}
	if err := handle.close("unused"); err != nil {
		t.Fatalf("close() error = %v", err)
	}
}

func TestCopiedHandleSharesState(t *testing.T) {
	handle := newNativeHandle(1)
	copyOfHandle := handle
	handle.state.value = 0
	if err := copyOfHandle.use(func(uintptr) error { return nil }); !errors.Is(err, ErrClosed) {
		t.Fatalf("copied handle error = %v, want ErrClosed", err)
	}
}

func TestPresetErrorMatchesStatusCode(t *testing.T) {
	err := &PresetError{Code: PresetErrBusy, Func: "test"}
	if !errors.Is(err, &PresetError{Code: PresetErrBusy}) {
		t.Fatal("matching status code was not recognized")
	}
	if errors.Is(err, &PresetError{Code: PresetErrTimeout}) {
		t.Fatal("different status code was recognized")
	}
}
