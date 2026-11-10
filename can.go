package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

// CAN / CAN-FD / UDS bindings for preset_rs (ABI v12).
//
// The structures below mirror the C layout in preset_rs.h field for field.
// Do not reorder fields: the DLL reads them directly through pointers.

const (
	PresetOK                    int32 = 0
	PresetErrNullPtr            int32 = -1
	PresetErrBufferTooSmall     int32 = -2
	PresetErrNotInit            int32 = -3
	PresetErrTimeout            int32 = -4
	PresetErrUdsNegative        int32 = -5
	PresetErrTransport          int32 = -6
	PresetErrUnsupported        int32 = -7
	PresetErrInvalidArg         int32 = -8
	PresetErrBusy               int32 = -9
	PresetErrIsotp              int32 = -10
	PresetErrPanic              int32 = -11
	PresetErrCancelled          int32 = -12
	PresetErrWaitTimeout        int32 = -13
	PresetErrNotReady           int32 = -14
	PresetErrDiagnosticOccupied int32 = -15
)

const (
	PresetCanDirectionTx            byte   = 0
	PresetCanDirectionRx            byte   = 1
	PresetCanBackendNone            byte   = 0
	PresetCanBackendToomoss         byte   = 1
	PresetCanBackendTSMaster        byte   = 2
	PresetCanBackendPCAN            byte   = 3
	PresetCanBackendVector          byte   = 4
	PresetAddressPhysical           byte   = 0
	PresetAddressFunctional         byte   = 1
	PresetOutcomeNone               uint32 = 0
	PresetOutcomePositive           uint32 = 1
	PresetOutcomeNegative           uint32 = 2
	PresetOutcomeSuppressedPositive uint32 = 3
	PresetOutcomeCancelled          uint32 = 4
	PresetOutcomeTimeout            uint32 = 5
	PresetOutcomeError              uint32 = 6
	PresetRuntimeRunning            uint32 = 0
	PresetRuntimeRecovering         uint32 = 1
	PresetRuntimeStopping           uint32 = 2
	PresetRuntimeStopped            uint32 = 3
	PresetRuntimeFaulted            uint32 = 4
	PresetRequestAccepted           uint32 = 0
	PresetRequestTransmitting       uint32 = 1
	PresetRequestWaitResponse       uint32 = 2
	PresetRequestReceiving          uint32 = 3
	PresetRequestResponsePending    uint32 = 4
	PresetRequestCancelRequested    uint32 = 5
	PresetRequestCompleted          uint32 = 6
	PresetRequestCancelled          uint32 = 7
	PresetRequestTimedOut           uint32 = 8
	PresetRequestFailed             uint32 = 9
)

var (
	preset_toomoss_default_config, _    = syscall.GetProcAddress(lib, "preset_toomoss_default_config")
	preset_toomoss_default_fd_timing, _ = syscall.GetProcAddress(lib, "preset_toomoss_default_fd_timing")
	preset_toomoss_can_init, _          = syscall.GetProcAddress(lib, "preset_toomoss_can_init")
	preset_device_get_backend, _        = syscall.GetProcAddress(lib, "preset_device_get_backend")

	preset_can_write_with_brs, _ = syscall.GetProcAddress(lib, "preset_can_write_with_brs")
	preset_can_write_ex, _       = syscall.GetProcAddress(lib, "preset_can_write_ex")
	preset_can_try_read_ex, _    = syscall.GetProcAddress(lib, "preset_can_try_read_ex")

	preset_can_runtime_new, _            = syscall.GetProcAddress(lib, "preset_can_runtime_new")
	preset_can_runtime_get_info, _       = syscall.GetProcAddress(lib, "preset_can_runtime_get_info")
	preset_can_runtime_copy_error, _     = syscall.GetProcAddress(lib, "preset_can_runtime_copy_error")
	preset_can_runtime_write_with_brs, _ = syscall.GetProcAddress(lib, "preset_can_runtime_write_with_brs")
	preset_can_runtime_write_ex, _       = syscall.GetProcAddress(lib, "preset_can_runtime_write_ex")
	preset_can_runtime_get_bus_load, _   = syscall.GetProcAddress(lib, "preset_can_runtime_get_bus_load")
	preset_can_runtime_stop, _           = syscall.GetProcAddress(lib, "preset_can_runtime_stop")
	preset_can_runtime_wait_stopped, _   = syscall.GetProcAddress(lib, "preset_can_runtime_wait_stopped")
	preset_can_runtime_close, _          = syscall.GetProcAddress(lib, "preset_can_runtime_close")
	preset_can_observer_new, _           = syscall.GetProcAddress(lib, "preset_can_observer_new")
	preset_can_observer_try_read, _      = syscall.GetProcAddress(lib, "preset_can_observer_try_read")
	preset_can_observer_get_stats, _     = syscall.GetProcAddress(lib, "preset_can_observer_get_stats")
	preset_can_observer_close, _         = syscall.GetProcAddress(lib, "preset_can_observer_close")
	preset_can_request_start, _          = syscall.GetProcAddress(lib, "preset_can_request_start")
	preset_request_get_info, _           = syscall.GetProcAddress(lib, "preset_request_get_info")
	preset_request_cancel, _             = syscall.GetProcAddress(lib, "preset_request_cancel")
	preset_request_wait, _               = syscall.GetProcAddress(lib, "preset_request_wait")
	preset_request_copy_response, _      = syscall.GetProcAddress(lib, "preset_request_copy_response")
	preset_request_copy_error, _         = syscall.GetProcAddress(lib, "preset_request_copy_error")
	preset_request_close, _              = syscall.GetProcAddress(lib, "preset_request_close")
	preset_last_error, _                 = syscall.GetProcAddress(lib, "preset_last_error")
)

// PresetError is a negative preset_rs status code plus the thread-local
// error text reported by preset_last_error (when available).
type PresetError struct {
	Code    int32
	Func    string
	Message string
}

func (e *PresetError) Error() string {
	text := fmt.Sprintf("%s failed: code %d", e.Func, e.Code)
	if e.Message != "" {
		text += ": " + e.Message
	}
	return text
}

func presetStatus(fn string, rc uintptr) error {
	code := int32(rc)
	if code >= 0 {
		return nil
	}
	return &PresetError{Code: code, Func: fn, Message: LastError()}
}

// LastError returns the latest synchronous error text on the calling thread.
func LastError() string {
	if preset_last_error == 0 {
		return ""
	}
	data, _, err := callCopy(preset_last_error, 0)
	if err != nil {
		return ""
	}
	return string(data)
}

// callCopy implements the two-call NULL/0 length query then copy pattern used
// by preset_last_error and preset_request_copy_*.
func callCopy(fn uintptr, handle uintptr) ([]byte, int32, error) {
	var n uintptr
	rc := int32(syscallR1(fn, handle, 0, 0, uintptr(unsafe.Pointer(&n))))
	if rc < 0 && rc != PresetErrBufferTooSmall {
		return nil, rc, &PresetError{Code: rc, Func: "preset_copy"}
	}
	if n == 0 {
		return nil, rc, nil
	}
	buf := make([]byte, n)
	rc = int32(syscallR1(fn, handle, uintptr(unsafe.Pointer(&buf[0])), n, uintptr(unsafe.Pointer(&n))))
	if rc < 0 {
		return nil, rc, &PresetError{Code: rc, Func: "preset_copy"}
	}
	return buf[:n], rc, nil
}

func bytePtr(data []byte) uintptr {
	if len(data) == 0 {
		return 0
	}
	return uintptr(unsafe.Pointer(&data[0]))
}

// syscallR1 exposes only the first return register of a native call.
func syscallR1(fn uintptr, args ...uintptr) uintptr {
	r1, _, _ := syscall.SyscallN(fn, args...)
	return r1
}

// ---- Structures ----

type CanFdTiming struct {
	NominalBrp   uint32
	NominalTseg1 uint32
	NominalTseg2 uint32
	NominalSjw   uint32
	DataBrp      uint32
	DataTseg1    uint32
	DataTseg2    uint32
	DataSjw      uint32
}

type CanChannelConfig struct {
	IsFd                byte
	Brs                 byte
	NominalBitrate      uint32
	DataBitrate         uint32
	RxBufferSize        uint32
	PollBatchSize       uint32
	MinTxIntervalUs     uint32
	MinRxPollIntervalUs uint32
}

type ToomossCanCfg struct {
	Channel byte
	Can     CanChannelConfig
}

type CanFrame struct {
	ID          uint32
	DataLen     byte
	Direction   byte
	IsFd        byte
	Brs         byte
	IsExtended  byte
	IsRemote    byte
	Esi         byte
	Reserved    byte
	TimestampUs uint64
	Data        [64]byte
}

func (f CanFrame) Payload() []byte {
	if f.IsRemote != 0 {
		return nil
	}
	n := int(f.DataLen)
	if n > len(f.Data) {
		n = len(f.Data)
	}
	return f.Data[:n]
}

func (f CanFrame) String() string {
	return fmt.Sprintf("CanFrame(id=0x%X, len=%d, fd=%d, data=% X)",
		f.ID, f.DataLen, f.IsFd, f.Payload())
}

type RxStats struct {
	Received uint64
	Dropped  uint64
	Queued   uint32
	Capacity uint32
}

type BusLoad struct {
	Load           float64
	WindowNs       uint64
	NominalBitrate uint32
	DataBitrate    uint32
	FrameCount     uint64
}

type CanRuntimeOptions struct {
	CommandCapacity uint32
	RecoveryMs      uint32
}

type CanObserverOptions struct {
	Capacity uint32
	Rx       byte
	Tx       byte
}

type RequestTimeouts struct {
	P2TimeoutMs     uint32
	P2StarTimeoutMs uint32
	TotalTimeoutMs  uint32
}

type CanDiagnosticConfig struct {
	PhysicalID        uint32
	ResponseID        uint32
	FunctionalID      uint32
	NBsMs             uint32
	NCrMs             uint32
	StMinMs           uint32
	MaxPduLen         uint32
	TxStMinOverrideMs uint32

	Addressing             byte
	PaddingEnabled         byte
	PaddingByte            byte
	BlockSize              byte
	MaxWaitFrames          byte
	TxStMinOverrideEnabled byte
	ManualFlowControl      byte
}

type DiagnosticRequestInfo struct {
	TransactionID   uint64
	RuntimeEpoch    uint64
	ElapsedUs       uint64
	ResponseLen     uint64
	State           uint32
	Outcome         uint32
	ErrorCode       int32
	TimeoutKind     uint32
	PendingCount    uint32
	ResponseAddress uint32
	LastNrc         byte
}

type CanCaptureEvent struct {
	Sequence      uint64
	RuntimeEpoch  uint64
	TransactionID uint64
	Frame         CanFrame
}

type CanRuntimeInfo struct {
	RuntimeEpoch        uint64
	TransactionID       uint64
	RecoveryRemainingUs uint64
	State               uint32
	PhysicalID          uint32
	ResponseID          uint32
	FunctionalID        uint32
	CallInProgress      byte
}

// ---- Default configs (struct returned by value) ----

// Structs larger than 8 bytes are returned through a hidden first argument on
// the Microsoft x64 calling convention, so pass the destination pointer.
func ToomossCanDefaultCfg() ToomossCanCfg {
	var cfg ToomossCanCfg
	syscall.SyscallN(preset_toomoss_default_config, uintptr(unsafe.Pointer(&cfg)))
	return cfg
}

func ToomossDefaultFdTiming() CanFdTiming {
	var timing CanFdTiming
	syscall.SyscallN(preset_toomoss_default_fd_timing, uintptr(unsafe.Pointer(&timing)))
	return timing
}

func DefaultRuntimeOptions() CanRuntimeOptions {
	return CanRuntimeOptions{CommandCapacity: 256, RecoveryMs: 100}
}

func DefaultRequestTimeouts() RequestTimeouts {
	return RequestTimeouts{P2TimeoutMs: 1000, P2StarTimeoutMs: 5000, TotalTimeoutMs: 30000}
}

func DefaultCanDiagnosticConfig(physicalID, responseID, functionalID uint32) CanDiagnosticConfig {
	return CanDiagnosticConfig{
		PhysicalID:     physicalID,
		ResponseID:     responseID,
		FunctionalID:   functionalID,
		NBsMs:          1000,
		NCrMs:          1000,
		StMinMs:        0,
		MaxPduLen:      4095,
		Addressing:     PresetAddressPhysical,
		PaddingEnabled: 1,
		PaddingByte:    0xCC,
		BlockSize:      0,
		MaxWaitFrames:  10,
	}
}

// ---- Device CAN setup ----

// CanInit initializes each supplied logical CAN channel with the library's
// default Toomoss configuration.
func (t *Toomoss) CanInit(channels ...byte) error {
	for _, v := range channels {
		cfg := ToomossCanDefaultCfg()
		cfg.Channel = v
		if err := t.CanInitCfg(&cfg, nil); err != nil {
			return err
		}
	}
	return nil
}

// CanInitCfg initializes one CAN channel. A nil timing pointer lets the
// library derive controller timing from the configured bitrates.
func (t *Toomoss) CanInitCfg(cfg *ToomossCanCfg, timing *CanFdTiming) error {
	var timingPtr uintptr
	if timing != nil {
		timingPtr = uintptr(unsafe.Pointer(timing))
	}
	rc := syscallR1(preset_toomoss_can_init, t.device, uintptr(unsafe.Pointer(cfg)), timingPtr)
	return presetStatus("preset_toomoss_can_init", rc)
}

// Backend returns PRESET_CAN_BACKEND_* for this device.
func (t *Toomoss) Backend() (byte, error) {
	var backend byte
	rc := syscallR1(preset_device_get_backend, t.device, uintptr(unsafe.Pointer(&backend)))
	if err := presetStatus("preset_device_get_backend", rc); err != nil {
		return 0, err
	}
	return backend, nil
}

// ---- Raw CAN I/O ----

// CanWrite sends a raw CAN/CAN-FD frame. Standard 11-bit IDs only; use
// CanWriteFrame for extended or remote frames.
func (t *Toomoss) CanWrite(channel byte, id uint32, data []byte, isFD, brs bool) error {
	rc := syscallR1(
		preset_can_write_with_brs,
		t.device,
		uintptr(channel),
		uintptr(id),
		uintptr(boolByte(isFD)),
		uintptr(boolByte(brs)),
		bytePtr(data),
		uintptr(len(data)))
	return presetStatus("preset_can_write_with_brs", rc)
}

// CanWriteFrame sends a structured frame (extended IDs, classic remote frames).
func (t *Toomoss) CanWriteFrame(channel byte, frame *CanFrame) error {
	rc := syscallR1(preset_can_write_ex, t.device, uintptr(channel), uintptr(unsafe.Pointer(frame)))
	return presetStatus("preset_can_write_ex", rc)
}

// CanTryRead drains up to maxFrames frames from the channel queue.
func (t *Toomoss) CanTryRead(channel byte, maxFrames int) ([]CanFrame, error) {
	if maxFrames <= 0 {
		return nil, nil
	}
	frames := make([]CanFrame, maxFrames)
	n := uintptr(maxFrames)
	rc := syscallR1(
		preset_can_try_read_ex,
		t.device,
		uintptr(channel),
		uintptr(unsafe.Pointer(&frames[0])),
		uintptr(unsafe.Pointer(&n)))
	if err := presetStatus("preset_can_try_read_ex", rc); err != nil {
		return nil, err
	}
	return frames[:n], nil
}

func boolByte(v bool) byte {
	if v {
		return 1
	}
	return 0
}

// MakeCanFrame fills a TX frame structure, mirroring the Python helper.
func MakeCanFrame(id uint32, data []byte, isFD, brs, isExtended, isRemote bool, dlc ...byte) CanFrame {
	frame := CanFrame{
		ID:         id,
		Direction:  PresetCanDirectionTx,
		IsFd:       boolByte(isFD),
		Brs:        boolByte(brs),
		IsExtended: boolByte(isExtended),
		IsRemote:   boolByte(isRemote),
	}
	n := copy(frame.Data[:], data)
	frame.DataLen = byte(n)
	if len(dlc) > 0 {
		frame.DataLen = dlc[0]
	}
	return frame
}

// ---- CAN runtime ----

type CanRuntime struct {
	handle uintptr
}

// CanRuntimeNew creates a diagnostic runtime that borrows an initialized
// channel. Pass no options to use 256 commands / 100 ms recovery.
func (t *Toomoss) CanRuntimeNew(channel byte, options ...CanRuntimeOptions) (*CanRuntime, error) {
	opts := DefaultRuntimeOptions()
	if len(options) > 0 {
		opts = options[0]
	}
	var handle uintptr
	rc := syscallR1(
		preset_can_runtime_new,
		t.device,
		uintptr(channel),
		uintptr(unsafe.Pointer(&opts)),
		uintptr(unsafe.Pointer(&handle)))
	if err := presetStatus("preset_can_runtime_new", rc); err != nil {
		return nil, err
	}
	return &CanRuntime{handle: handle}, nil
}

func (r *CanRuntime) GetInfo() (CanRuntimeInfo, error) {
	var info CanRuntimeInfo
	rc := syscallR1(preset_can_runtime_get_info, r.handle, uintptr(unsafe.Pointer(&info)))
	if err := presetStatus("preset_can_runtime_get_info", rc); err != nil {
		return CanRuntimeInfo{}, err
	}
	return info, nil
}

// CopyError returns the runtime's last CAN/driver error text.
func (r *CanRuntime) CopyError() (string, error) {
	data, _, err := callCopy(preset_can_runtime_copy_error, r.handle)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (r *CanRuntime) Write(id uint32, data []byte, isFD, brs bool) error {
	rc := syscallR1(
		preset_can_runtime_write_with_brs,
		r.handle,
		uintptr(id),
		uintptr(boolByte(isFD)),
		uintptr(boolByte(brs)),
		bytePtr(data),
		uintptr(len(data)))
	return presetStatus("preset_can_runtime_write_with_brs", rc)
}

func (r *CanRuntime) WriteFrame(frame *CanFrame) error {
	rc := syscallR1(preset_can_runtime_write_ex, r.handle, uintptr(unsafe.Pointer(frame)))
	return presetStatus("preset_can_runtime_write_ex", rc)
}

func (r *CanRuntime) GetBusLoad() (BusLoad, error) {
	var load BusLoad
	rc := syscallR1(preset_can_runtime_get_bus_load, r.handle, uintptr(unsafe.Pointer(&load)))
	if err := presetStatus("preset_can_runtime_get_bus_load", rc); err != nil {
		return BusLoad{}, err
	}
	return load, nil
}

func (r *CanRuntime) Stop() error {
	return presetStatus("preset_can_runtime_stop", syscallR1(preset_can_runtime_stop, r.handle))
}

// WaitStopped returns true when stopped, false when only the wait expired
// (PRESET_ERR_WAIT_TIMEOUT).
func (r *CanRuntime) WaitStopped(timeoutMs uint32) (bool, error) {
	rc := int32(syscallR1(preset_can_runtime_wait_stopped, r.handle, uintptr(timeoutMs)))
	if rc == PresetErrWaitTimeout {
		return false, nil
	}
<<<<<<< Updated upstream
	return withHandle(client.state, func(handle uintptr) int32 {
		return invokeStatus(
			"preset_can_uds_write",
			word(handle), word(uintptr(id)), word(fd), slicePointer(payload), word(uintptr(len(payload))),
		)
	})
}

func CanUdsWriteWithBrs(
	client PresetCanUdsClient,
	id uint32,
	isFD bool,
	brs bool,
	data []byte) int32 {
	var fd, _brs = 0, 0
	if isFD {
		fd = 1
	}
	if brs {
		_brs = 1
	}
	return withHandle(client.state, func(handle uintptr) int32 {
		return invokeStatus(
			"preset_can_uds_write_with_brs",
			word(handle),
			word(uintptr(id)),
			word(uintptr(fd)),
			word(uintptr(_brs)),
			slicePointer(data),
			word(uintptr(len(data))),
		)
	})
=======
	if rc < 0 {
		return false, &PresetError{Code: rc, Func: "preset_can_runtime_wait_stopped", Message: LastError()}
	}
	return true, nil
>>>>>>> Stashed changes
}

func (r *CanRuntime) Close() error {
	if r == nil || r.handle == 0 {
		return nil
	}
	rc := syscallR1(preset_can_runtime_close, uintptr(unsafe.Pointer(&r.handle)))
	return presetStatus("preset_can_runtime_close", rc)
}

// Shutdown performs stop + wait_stopped + close.
func (r *CanRuntime) Shutdown(waitTimeoutMs ...uint32) error {
	if r == nil || r.handle == 0 {
		return nil
	}
	waitMs := uint32(2000)
	if len(waitTimeoutMs) > 0 {
		waitMs = waitTimeoutMs[0]
	}
	if err := r.Stop(); err != nil {
		return err
	}
	ok, err := r.WaitStopped(waitMs)
	if err != nil {
		return err
	}
	if !ok {
		return &PresetError{Code: PresetErrWaitTimeout, Func: "can_runtime_shutdown", Message: "runtime did not stop in time"}
	}
	return r.Close()
}

// ---- Observer ----

type CanObserver struct {
	handle uintptr
}

func (r *CanRuntime) ObserverNew(options ...CanObserverOptions) (*CanObserver, error) {
	opts := CanObserverOptions{Capacity: 4096, Rx: 1, Tx: 1}
	if len(options) > 0 {
		opts = options[0]
	}
	var handle uintptr
	rc := syscallR1(
		preset_can_observer_new,
		r.handle,
		uintptr(unsafe.Pointer(&opts)),
		uintptr(unsafe.Pointer(&handle)))
	if err := presetStatus("preset_can_observer_new", rc); err != nil {
		return nil, err
	}
	return &CanObserver{handle: handle}, nil
}

func (o *CanObserver) TryRead(maxEvents int) ([]CanCaptureEvent, error) {
	if maxEvents <= 0 {
		return nil, nil
	}
	events := make([]CanCaptureEvent, maxEvents)
	var n uintptr
	rc := syscallR1(
		preset_can_observer_try_read,
		o.handle,
		uintptr(unsafe.Pointer(&events[0])),
		uintptr(maxEvents),
		uintptr(unsafe.Pointer(&n)))
	if err := presetStatus("preset_can_observer_try_read", rc); err != nil {
		return nil, err
	}
	return events[:n], nil
}

func (o *CanObserver) GetStats() (RxStats, error) {
	var stats RxStats
	rc := syscallR1(preset_can_observer_get_stats, o.handle, uintptr(unsafe.Pointer(&stats)))
	if err := presetStatus("preset_can_observer_get_stats", rc); err != nil {
		return RxStats{}, err
	}
	return stats, nil
}

func (o *CanObserver) Close() error {
	if o == nil || o.handle == 0 {
		return nil
	}
	rc := syscallR1(preset_can_observer_close, uintptr(unsafe.Pointer(&o.handle)))
	return presetStatus("preset_can_observer_close", rc)
}

// ---- Diagnostic requests ----

type DiagnosticRequest struct {
	handle uintptr
}

// StartRequest copies the PDU and configuration, then returns without waiting
// for the ECU.
func (r *CanRuntime) StartRequest(cfg *CanDiagnosticConfig, timeouts *RequestTimeouts, pdu []byte) (*DiagnosticRequest, error) {
	var handle uintptr
	rc := syscallR1(
		preset_can_request_start,
		r.handle,
		uintptr(unsafe.Pointer(cfg)),
		uintptr(unsafe.Pointer(timeouts)),
		bytePtr(pdu),
		uintptr(len(pdu)),
		uintptr(unsafe.Pointer(&handle)))
	if err := presetStatus("preset_can_request_start", rc); err != nil {
		return nil, err
	}
	return &DiagnosticRequest{handle: handle}, nil
}

func (q *DiagnosticRequest) GetInfo() (DiagnosticRequestInfo, error) {
	var info DiagnosticRequestInfo
	rc := syscallR1(preset_request_get_info, q.handle, uintptr(unsafe.Pointer(&info)))
	if err := presetStatus("preset_request_get_info", rc); err != nil {
		return DiagnosticRequestInfo{}, err
	}
	return info, nil
}

func (q *DiagnosticRequest) Cancel() error {
	return presetStatus("preset_request_cancel", syscallR1(preset_request_cancel, q.handle))
}

// Wait returns true when the request reached a terminal state, false when only
// this wait expired (PRESET_ERR_WAIT_TIMEOUT).
func (q *DiagnosticRequest) Wait(timeoutMs uint32) (bool, error) {
	rc := int32(syscallR1(preset_request_wait, q.handle, uintptr(timeoutMs)))
	if rc == PresetErrWaitTimeout {
		return false, nil
	}
	if rc < 0 {
		return false, &PresetError{Code: rc, Func: "preset_request_wait", Message: LastError()}
	}
	return true, nil
}

// CopyResponse returns the response bytes, or nil when the request is not
// terminal yet (PRESET_ERR_NOT_READY).
func (q *DiagnosticRequest) CopyResponse() ([]byte, error) {
	data, rc, err := callCopy(preset_request_copy_response, q.handle)
	if rc == PresetErrNotReady {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (q *DiagnosticRequest) CopyError() (string, error) {
	data, _, err := callCopy(preset_request_copy_error, q.handle)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (q *DiagnosticRequest) Close() error {
	if q == nil || q.handle == 0 {
		return nil
	}
	rc := syscallR1(preset_request_close, uintptr(unsafe.Pointer(&q.handle)))
	return presetStatus("preset_request_close", rc)
}

// UdsRequest starts a CAN UDS request, waits for the result, reads the
// response and releases the request handle.
//
// A UDS negative response (7F xx NRC) is returned normally: check info.Outcome
// and info.LastNrc. Transport errors, timeouts and cancellation return an
// error together with the request info and whatever response was collected.
func (r *CanRuntime) UdsRequest(cfg *CanDiagnosticConfig, timeouts *RequestTimeouts, pdu []byte, waitMs uint32) (DiagnosticRequestInfo, []byte, error) {
	req, err := r.StartRequest(cfg, timeouts, pdu)
	if err != nil {
		return DiagnosticRequestInfo{}, nil, err
	}
	return finishRequest(req, waitMs)
}

func finishRequest(req *DiagnosticRequest, waitMs uint32) (DiagnosticRequestInfo, []byte, error) {
	defer req.Close()

	done, err := req.Wait(waitMs)
	if err != nil {
		return DiagnosticRequestInfo{}, nil, err
	}
	if !done {
		_ = req.Cancel()
		_, _ = req.Wait(1000)
	}

	info, err := req.GetInfo()
	if err != nil {
		return DiagnosticRequestInfo{}, nil, err
	}
	response, err := req.CopyResponse()
	if err != nil {
		return info, nil, err
	}

	if (info.Outcome == PresetOutcomeError ||
		info.Outcome == PresetOutcomeTimeout ||
		info.Outcome == PresetOutcomeCancelled) && info.ErrorCode < 0 {
		text, _ := req.CopyError()
		return info, response, &PresetError{Code: info.ErrorCode, Func: "uds_request", Message: text}
	}
	return info, response, nil
}
