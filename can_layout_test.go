//go:build windows && amd64

package main

import (
	"testing"
	"unsafe"
)

// These sizes/offsets must match preset_rs.h (ABI v12). The DLL reads the
// structs through raw pointers, so any layout drift is a memory-safety bug.
func TestCanStructLayout(t *testing.T) {
	type layout struct {
		name string
		got  uintptr
		want uintptr
	}
	checks := []layout{
		{"CanFdTiming", unsafe.Sizeof(CanFdTiming{}), 32},
		{"CanChannelConfig", unsafe.Sizeof(CanChannelConfig{}), 28},
		{"ToomossCanCfg", unsafe.Sizeof(ToomossCanCfg{}), 32},
		{"CanFrame", unsafe.Sizeof(CanFrame{}), 88},
		{"RxStats", unsafe.Sizeof(RxStats{}), 24},
		{"BusLoad", unsafe.Sizeof(BusLoad{}), 32},
		{"CanRuntimeOptions", unsafe.Sizeof(CanRuntimeOptions{}), 8},
		{"CanObserverOptions", unsafe.Sizeof(CanObserverOptions{}), 8},
		{"RequestTimeouts", unsafe.Sizeof(RequestTimeouts{}), 12},
		{"CanDiagnosticConfig", unsafe.Sizeof(CanDiagnosticConfig{}), 40},
		{"DiagnosticRequestInfo", unsafe.Sizeof(DiagnosticRequestInfo{}), 64},
		{"CanCaptureEvent", unsafe.Sizeof(CanCaptureEvent{}), 112},
		{"CanRuntimeInfo", unsafe.Sizeof(CanRuntimeInfo{}), 48},
		{"AutoCfg", unsafe.Sizeof(AutoCfg{}), 36},
		{"PCANCfg", unsafe.Sizeof(PCANCfg{}), 32},
		{"TSMasterCfg", unsafe.Sizeof(TSMasterCfg{}), 40},
		{"VectorCfg", unsafe.Sizeof(VectorCfg{}), 72},
		{"ToomossLinCfg", unsafe.Sizeof(ToomossLinCfg{}), 8},
		{"LinFrame", unsafe.Sizeof(LinFrame{}), 24},
		{"PCANLinCfg", unsafe.Sizeof(PCANLinCfg{}), 8},
		{"TSMasterLinCfg", unsafe.Sizeof(TSMasterLinCfg{}), 16},
		{"VectorLinCfg", unsafe.Sizeof(VectorLinCfg{}), 24},
		{"ToomossElinsCfg", unsafe.Sizeof(ToomossElinsCfg{}), 12},
		{"ToomossElinsMessage", unsafe.Sizeof(ToomossElinsMessage{}), 88},
	}
	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("%s size = %d, want %d", c.name, c.got, c.want)
		}
	}

	checks = []layout{
		{"CanChannelConfig.Can", unsafe.Offsetof(ToomossCanCfg{}.Can), 4},
		{"CanFrame.TimestampUs", unsafe.Offsetof(CanFrame{}.TimestampUs), 16},
		{"CanFrame.Data", unsafe.Offsetof(CanFrame{}.Data), 24},
		{"CanDiagnosticConfig.Addressing", unsafe.Offsetof(CanDiagnosticConfig{}.Addressing), 32},
		{"DiagnosticRequestInfo.LastNrc", unsafe.Offsetof(DiagnosticRequestInfo{}.LastNrc), 56},
		{"CanRuntimeInfo.CallInProgress", unsafe.Offsetof(CanRuntimeInfo{}.CallInProgress), 40},
		{"CanCaptureEvent.Frame", unsafe.Offsetof(CanCaptureEvent{}.Frame), 24},
		{"AutoCfg.Can", unsafe.Offsetof(AutoCfg{}.Can), 4},
		{"AutoCfg.CandidateOrder", unsafe.Offsetof(AutoCfg{}.CandidateOrder), 32},
		{"TSMasterCfg.Can", unsafe.Offsetof(TSMasterCfg{}.Can), 12},
		{"VectorCfg.Can", unsafe.Offsetof(VectorCfg{}.Can), 44},
		{"LinFrame.TimestampUs", unsafe.Offsetof(LinFrame{}.TimestampUs), 16},
		{"ToomossElinsMessage.Data", unsafe.Offsetof(ToomossElinsMessage{}.Data), 18},
	}
	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("%s offset = %d, want %d", c.name, c.got, c.want)
		}
	}
}
