//go:build windows && amd64

package preset

import (
	"testing"
	"unsafe"
)

// These sizes/offsets must match preset_rs.h (ABI v12). The DLL reads the
// structs through raw pointers, so any layout drift is a memory-safety bug.
func TestStructLayout(t *testing.T) {
	type layout struct {
		name string
		got  uintptr
		want uintptr
	}
	checks := []layout{
		{"CANFDTiming", unsafe.Sizeof(CANFDTiming{}), 32},
		{"CANChannelConfig", unsafe.Sizeof(CANChannelConfig{}), 28},
		{"ToomossCANConfig", unsafe.Sizeof(ToomossCANConfig{}), 32},
		{"CANFrame", unsafe.Sizeof(CANFrame{}), 88},
		{"RXStats", unsafe.Sizeof(RXStats{}), 24},
		{"BusLoad", unsafe.Sizeof(BusLoad{}), 32},
		{"CANRuntimeOptions", unsafe.Sizeof(CANRuntimeOptions{}), 8},
		{"CANObserverOptions", unsafe.Sizeof(CANObserverOptions{}), 8},
		{"RequestTimeouts", unsafe.Sizeof(RequestTimeouts{}), 12},
		{"CANDiagnosticConfig", unsafe.Sizeof(CANDiagnosticConfig{}), 40},
		{"DiagnosticRequestInfo", unsafe.Sizeof(DiagnosticRequestInfo{}), 64},
		{"CANCaptureEvent", unsafe.Sizeof(CANCaptureEvent{}), 112},
		{"CANRuntimeInfo", unsafe.Sizeof(CANRuntimeInfo{}), 48},
		{"AutoConfig", unsafe.Sizeof(AutoConfig{}), 36},
		{"PCANConfig", unsafe.Sizeof(PCANConfig{}), 32},
		{"TSMasterConfig", unsafe.Sizeof(TSMasterConfig{}), 40},
		{"VectorConfig", unsafe.Sizeof(VectorConfig{}), 72},
		{"ToomossLINConfig", unsafe.Sizeof(ToomossLINConfig{}), 8},
		{"LINFrame", unsafe.Sizeof(LINFrame{}), 24},
		{"PCANLINConfig", unsafe.Sizeof(PCANLINConfig{}), 8},
		{"TSMasterLINConfig", unsafe.Sizeof(TSMasterLINConfig{}), 16},
		{"VectorLINConfig", unsafe.Sizeof(VectorLINConfig{}), 24},
		{"ToomossELINSConfig", unsafe.Sizeof(ToomossELINSConfig{}), 12},
		{"ToomossELINSMessage", unsafe.Sizeof(ToomossELINSMessage{}), 88},
	}
	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("%s size = %d, want %d", c.name, c.got, c.want)
		}
	}

	checks = []layout{
		{"ToomossCANConfig.CAN", unsafe.Offsetof(ToomossCANConfig{}.CAN), 4},
		{"CANFrame.TimestampUS", unsafe.Offsetof(CANFrame{}.TimestampUS), 16},
		{"CANFrame.Data", unsafe.Offsetof(CANFrame{}.Data), 24},
		{"CANDiagnosticConfig.Addressing", unsafe.Offsetof(CANDiagnosticConfig{}.Addressing), 32},
		{"DiagnosticRequestInfo.LastNRC", unsafe.Offsetof(DiagnosticRequestInfo{}.LastNRC), 56},
		{"CANRuntimeInfo.CallInProgress", unsafe.Offsetof(CANRuntimeInfo{}.CallInProgress), 40},
		{"CANCaptureEvent.Frame", unsafe.Offsetof(CANCaptureEvent{}.Frame), 24},
		{"AutoConfig.CAN", unsafe.Offsetof(AutoConfig{}.CAN), 4},
		{"AutoConfig.CandidateOrder", unsafe.Offsetof(AutoConfig{}.CandidateOrder), 32},
		{"TSMasterConfig.CAN", unsafe.Offsetof(TSMasterConfig{}.CAN), 12},
		{"VectorConfig.CAN", unsafe.Offsetof(VectorConfig{}.CAN), 44},
		{"LINFrame.TimestampUS", unsafe.Offsetof(LINFrame{}.TimestampUS), 16},
		{"ToomossELINSMessage.Data", unsafe.Offsetof(ToomossELINSMessage{}.Data), 18},
	}
	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("%s offset = %d, want %d", c.name, c.got, c.want)
		}
	}
}
