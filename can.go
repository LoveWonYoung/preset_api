//go:build windows && amd64

package preset

import (
	"fmt"
	"runtime"
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
	PresetErrUDSNegative        int32 = -5
	PresetErrTransport          int32 = -6
	PresetErrUnsupported        int32 = -7
	PresetErrInvalidArg         int32 = -8
	PresetErrBusy               int32 = -9
	PresetErrISOTP              int32 = -10
	PresetErrPanic              int32 = -11
	PresetErrCancelled          int32 = -12
	PresetErrWaitTimeout        int32 = -13
	PresetErrNotReady           int32 = -14
	PresetErrDiagnosticOccupied int32 = -15
)

const (
	PresetCANDirectionTX            byte   = 0
	PresetCANDirectionRX            byte   = 1
	PresetCANBackendNone            byte   = 0
	PresetCANBackendToomoss         byte   = 1
	PresetCANBackendTSMaster        byte   = 2
	PresetCANBackendPCAN            byte   = 3
	PresetCANBackendVector          byte   = 4
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

// Is matches PresetError values by native status code.
func (e *PresetError) Is(target error) bool {
	other, ok := target.(*PresetError)
	return ok && e.Code == other.Code
}

func presetStatus(fn string, rc uintptr) error {
	code := int32(rc)
	if code >= 0 {
		return nil
	}
	return &PresetError{Code: code, Func: fn, Message: lastError()}
}

func lastError() string {
	fn, err := proc("preset_last_error")
	if err != nil {
		return ""
	}
	var n uintptr
	rc := int32(syscallR1(fn, 0, 0, uintptr(unsafe.Pointer(&n))))
	if rc < 0 && rc != PresetErrBufferTooSmall || n == 0 {
		return ""
	}
	data := make([]byte, n)
	if int32(syscallR1(fn, bytePtr(data), n, uintptr(unsafe.Pointer(&n)))) < 0 {
		return ""
	}
	return string(data[:n])
}

// callCopy implements the two-call NULL/0 length query then copy pattern used
// by preset_last_error and preset_request_copy_*.
func callCopy(name string, handle uintptr) ([]byte, int32, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	fn, err := proc(name)
	if err != nil {
		return nil, 0, err
	}
	var n uintptr
	rc := int32(syscallR1(fn, handle, 0, 0, uintptr(unsafe.Pointer(&n))))
	if rc < 0 && rc != PresetErrBufferTooSmall {
		return nil, rc, &PresetError{Code: rc, Func: name, Message: lastError()}
	}
	if n == 0 {
		return nil, rc, nil
	}
	buf := make([]byte, n)
	rc = int32(syscallR1(fn, handle, uintptr(unsafe.Pointer(&buf[0])), n, uintptr(unsafe.Pointer(&n))))
	if rc < 0 {
		return nil, rc, &PresetError{Code: rc, Func: name, Message: lastError()}
	}
	return buf[:n], rc, nil
}

func nativeWait(name string, args ...uintptr) (bool, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	result, err := nativeResult(name, args...)
	if err != nil {
		return false, err
	}
	code := int32(result)
	if code == PresetErrWaitTimeout {
		return false, nil
	}
	if code < 0 {
		return false, &PresetError{Code: code, Func: name, Message: lastError()}
	}
	return true, nil
}

func bytePtr(data []byte) uintptr {
	if len(data) == 0 {
		return 0
	}
	return uintptr(unsafe.Pointer(&data[0]))
}

// ---- Structures ----

type CANFDTiming struct {
	NominalBrp   uint32
	NominalTseg1 uint32
	NominalTseg2 uint32
	NominalSjw   uint32
	DataBrp      uint32
	DataTseg1    uint32
	DataTseg2    uint32
	DataSjw      uint32
}

type CANChannelConfig struct {
	IsFD                byte
	BRS                 byte
	NominalBitrate      uint32
	DataBitrate         uint32
	RxBufferSize        uint32
	PollBatchSize       uint32
	MinTXIntervalUS     uint32
	MinRXPollIntervalUS uint32
}

type ToomossCANConfig struct {
	Channel byte
	CAN     CANChannelConfig
}

type CANFrame struct {
	ID          uint32
	DataLen     byte
	Direction   byte
	IsFD        byte
	BRS         byte
	IsExtended  byte
	IsRemote    byte
	ESI         byte
	Reserved    byte
	TimestampUS uint64
	Data        [64]byte
}

func (f CANFrame) Payload() []byte {
	if f.IsRemote != 0 {
		return nil
	}
	n := int(f.DataLen)
	if n > len(f.Data) {
		n = len(f.Data)
	}
	return f.Data[:n]
}

func (f CANFrame) String() string {
	return fmt.Sprintf("CANFrame(id=0x%X, len=%d, fd=%d, data=% X)",
		f.ID, f.DataLen, f.IsFD, f.Payload())
}

// Validate checks the transmit-side constraints documented by preset_rs.h.
func (f CANFrame) Validate() error {
	if f.Direction != PresetCANDirectionTX {
		return fmt.Errorf("CAN frame direction must be TX")
	}
	if f.IsFD > 1 || f.BRS > 1 || f.IsExtended > 1 || f.IsRemote > 1 || f.ESI > 1 {
		return fmt.Errorf("CAN frame flags must be 0 or 1")
	}
	if f.Reserved != 0 || f.ESI != 0 {
		return fmt.Errorf("CAN frame ESI and reserved fields must be zero when transmitting")
	}
	maxID := uint32(0x7FF)
	if f.IsExtended != 0 {
		maxID = 0x1FFFFFFF
	}
	if f.ID > maxID {
		return fmt.Errorf("CAN ID 0x%X exceeds the selected frame format", f.ID)
	}
	if f.IsRemote != 0 {
		if f.IsFD != 0 {
			return fmt.Errorf("CAN-FD remote frames are invalid")
		}
		if f.DataLen > 8 {
			return fmt.Errorf("classic CAN remote DLC must be at most 8")
		}
		return nil
	}
	maxLen := byte(8)
	if f.IsFD != 0 {
		maxLen = 64
	}
	if f.DataLen > maxLen {
		return fmt.Errorf("CAN payload length %d exceeds maximum %d", f.DataLen, maxLen)
	}
	return nil
}

type RXStats struct {
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

type CANRuntimeOptions struct {
	CommandCapacity uint32
	RecoveryMs      uint32
}

type CANObserverOptions struct {
	Capacity uint32
	RX       byte
	TX       byte
}

type RequestTimeouts struct {
	P2TimeoutMs     uint32
	P2StarTimeoutMs uint32
	TotalTimeoutMs  uint32
}

type CANDiagnosticConfig struct {
	PhysicalID        uint32
	ResponseID        uint32
	FunctionalID      uint32
	NBsMs             uint32
	NCrMs             uint32
	StMinMs           uint32
	MaxPDULen         uint32
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
	ElapsedUS       uint64
	ResponseLen     uint64
	State           uint32
	Outcome         uint32
	ErrorCode       int32
	TimeoutKind     uint32
	PendingCount    uint32
	ResponseAddress uint32
	LastNRC         byte
}

type CANCaptureEvent struct {
	Sequence      uint64
	RuntimeEpoch  uint64
	TransactionID uint64
	Frame         CANFrame
}

type CANRuntimeInfo struct {
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
func DefaultToomossCANConfig() (ToomossCANConfig, error) {
	return defaultStruct[ToomossCANConfig]("preset_toomoss_default_config")
}

func DefaultToomossCANFDTiming() (CANFDTiming, error) {
	return defaultStruct[CANFDTiming]("preset_toomoss_default_fd_timing")
}

func DefaultRuntimeOptions() CANRuntimeOptions {
	return CANRuntimeOptions{CommandCapacity: 256, RecoveryMs: 100}
}

func DefaultRequestTimeouts() RequestTimeouts {
	return RequestTimeouts{P2TimeoutMs: 1000, P2StarTimeoutMs: 5000, TotalTimeoutMs: 30000}
}

func DefaultCANDiagnosticConfig(physicalID, responseID, functionalID uint32) CANDiagnosticConfig {
	return CANDiagnosticConfig{
		PhysicalID:     physicalID,
		ResponseID:     responseID,
		FunctionalID:   functionalID,
		NBsMs:          1000,
		NCrMs:          1000,
		StMinMs:        0,
		MaxPDULen:      4095,
		Addressing:     PresetAddressPhysical,
		PaddingEnabled: 1,
		PaddingByte:    0xCC,
		BlockSize:      0,
		MaxWaitFrames:  10,
	}
}

// ---- Device CAN setup ----

// InitCAN initializes each supplied logical CAN channel with the library's
// default Toomoss configuration.
func (t *Toomoss) InitCAN(channels ...byte) error {
	for _, v := range channels {
		cfg, err := DefaultToomossCANConfig()
		if err != nil {
			return err
		}
		cfg.Channel = v
		if err := t.InitCANWithConfig(&cfg, nil); err != nil {
			return err
		}
	}
	return nil
}

// InitCANWithConfig initializes one CAN channel. A nil timing pointer lets the
// library derive controller timing from the configured bitrates.
func (t *Toomoss) InitCANWithConfig(cfg *ToomossCANConfig, timing *CANFDTiming) error {
	if cfg == nil {
		return fmt.Errorf("Toomoss config is nil")
	}
	return t.handle.use(func(handle uintptr) error {
		return nativeStatus("preset_toomoss_can_init", handle, uintptr(unsafe.Pointer(cfg)), uintptr(unsafe.Pointer(timing)))
	})
}

// Backend returns PRESET_CAN_BACKEND_* for this device.
func (t *Device) Backend() (byte, error) {
	return useHandle(&t.handle, func(handle uintptr) (byte, error) {
		var backend byte
		if err := nativeStatus("preset_device_get_backend", handle, uintptr(unsafe.Pointer(&backend))); err != nil {
			return 0, err
		}
		return backend, nil
	})
}

// ---- Raw CAN I/O ----

// WriteCAN sends a raw CAN/CAN-FD frame. Standard 11-bit IDs only; use
// WriteCANFrame for extended or remote frames.
func (t *Device) WriteCAN(channel byte, id uint32, data []byte, isFD, brs bool) error {
	if id > 0x7FF {
		return fmt.Errorf("standard CAN ID must be at most 0x7FF")
	}
	if err := validateCanPayload(data, isFD); err != nil {
		return err
	}
	return t.handle.use(func(handle uintptr) error {
		return nativeStatus(
			"preset_can_write_with_brs",
			handle,
			uintptr(channel),
			uintptr(id),
			uintptr(boolByte(isFD)),
			uintptr(boolByte(brs)),
			bytePtr(data),
			uintptr(len(data)))
	})
}

// WriteCANFrame sends a structured frame (extended IDs, classic remote frames).
func (t *Device) WriteCANFrame(channel byte, frame *CANFrame) error {
	if frame == nil {
		return fmt.Errorf("CAN frame is nil")
	}
	if err := frame.Validate(); err != nil {
		return err
	}
	return t.handle.use(func(handle uintptr) error {
		return nativeStatus("preset_can_write_ex", handle, uintptr(channel), uintptr(unsafe.Pointer(frame)))
	})
}

// TryReadCAN drains up to maxFrames frames from the channel queue.
func (t *Device) TryReadCAN(channel byte, maxFrames int) ([]CANFrame, error) {
	if maxFrames < 1 || maxFrames > 65536 {
		return nil, fmt.Errorf("CAN read capacity must be in range 1-65536")
	}
	frames := make([]CANFrame, maxFrames)
	n := uintptr(maxFrames)
	err := t.handle.use(func(handle uintptr) error {
		return nativeStatus(
			"preset_can_try_read_ex",
			handle,
			uintptr(channel),
			uintptr(unsafe.Pointer(&frames[0])),
			uintptr(unsafe.Pointer(&n)))
	})
	if err != nil {
		return nil, err
	}
	if n > uintptr(len(frames)) {
		return nil, fmt.Errorf("preset_can_try_read_ex returned invalid length %d", n)
	}
	return frames[:n], nil
}

func boolByte(v bool) byte {
	if v {
		return 1
	}
	return 0
}

func validateCanPayload(data []byte, isFD bool) error {
	maxLen := 8
	if isFD {
		maxLen = 64
	}
	if len(data) > maxLen {
		return fmt.Errorf("CAN payload length %d exceeds maximum %d", len(data), maxLen)
	}
	return nil
}

// NewCANFrame creates and validates a classic CAN or CAN-FD data frame.
func NewCANFrame(id uint32, data []byte, isFD, brs, isExtended bool) (CANFrame, error) {
	if err := validateCanPayload(data, isFD); err != nil {
		return CANFrame{}, err
	}
	frame := CANFrame{
		ID:         id,
		Direction:  PresetCANDirectionTX,
		IsFD:       boolByte(isFD),
		BRS:        boolByte(brs),
		IsExtended: boolByte(isExtended),
		DataLen:    byte(len(data)),
	}
	copy(frame.Data[:], data)
	if err := frame.Validate(); err != nil {
		return CANFrame{}, err
	}
	return frame, nil
}

// NewCANRemoteFrame creates and validates a classic CAN remote frame.
func NewCANRemoteFrame(id uint32, dlc byte, isExtended bool) (CANFrame, error) {
	frame := CANFrame{
		ID:         id,
		DataLen:    dlc,
		Direction:  PresetCANDirectionTX,
		IsExtended: boolByte(isExtended),
		IsRemote:   1,
	}
	if err := frame.Validate(); err != nil {
		return CANFrame{}, err
	}
	return frame, nil
}

// ---- CAN runtime ----

type CANRuntime struct{ handle nativeHandle }

// NewCANRuntime creates a diagnostic runtime that borrows an initialized
// channel. Pass no options to use 256 commands / 100 ms recovery.
func (t *Device) NewCANRuntime(channel byte, options ...CANRuntimeOptions) (*CANRuntime, error) {
	if len(options) > 1 {
		return nil, fmt.Errorf("at most one CAN runtime options value is allowed")
	}
	opts := DefaultRuntimeOptions()
	if len(options) > 0 {
		opts = options[0]
	}
	return useHandle(&t.handle, func(device uintptr) (*CANRuntime, error) {
		var handle uintptr
		if err := nativeStatus(
			"preset_can_runtime_new",
			device,
			uintptr(channel),
			uintptr(unsafe.Pointer(&opts)),
			uintptr(unsafe.Pointer(&handle))); err != nil {
			return nil, err
		}
		return &CANRuntime{handle: newNativeHandle(handle)}, nil
	})
}

func (r *CANRuntime) GetInfo() (CANRuntimeInfo, error) {
	return useHandle(&r.handle, func(handle uintptr) (CANRuntimeInfo, error) {
		var info CANRuntimeInfo
		if err := nativeStatus("preset_can_runtime_get_info", handle, uintptr(unsafe.Pointer(&info))); err != nil {
			return CANRuntimeInfo{}, err
		}
		return info, nil
	})
}

// CopyError returns the runtime's last CAN/driver error text.
func (r *CANRuntime) CopyError() (string, error) {
	return useHandle(&r.handle, func(handle uintptr) (string, error) {
		data, _, err := callCopy("preset_can_runtime_copy_error", handle)
		return string(data), err
	})
}

func (r *CANRuntime) Write(id uint32, data []byte, isFD, brs bool) error {
	if id > 0x7FF {
		return fmt.Errorf("standard CAN ID must be at most 0x7FF")
	}
	if err := validateCanPayload(data, isFD); err != nil {
		return err
	}
	return r.handle.use(func(handle uintptr) error {
		return nativeStatus(
			"preset_can_runtime_write_with_brs",
			handle,
			uintptr(id),
			uintptr(boolByte(isFD)),
			uintptr(boolByte(brs)),
			bytePtr(data),
			uintptr(len(data)))
	})
}

func (r *CANRuntime) WriteFrame(frame *CANFrame) error {
	if frame == nil {
		return fmt.Errorf("CAN frame is nil")
	}
	if err := frame.Validate(); err != nil {
		return err
	}
	return r.handle.use(func(handle uintptr) error {
		return nativeStatus("preset_can_runtime_write_ex", handle, uintptr(unsafe.Pointer(frame)))
	})
}

func (r *CANRuntime) GetBusLoad() (BusLoad, error) {
	return useHandle(&r.handle, func(handle uintptr) (BusLoad, error) {
		var load BusLoad
		if err := nativeStatus("preset_can_runtime_get_bus_load", handle, uintptr(unsafe.Pointer(&load))); err != nil {
			return BusLoad{}, err
		}
		return load, nil
	})
}

func (r *CANRuntime) Stop() error {
	return r.handle.use(func(handle uintptr) error {
		return nativeStatus("preset_can_runtime_stop", handle)
	})
}

// WaitStopped returns true when stopped, false when only the wait expired
// (PRESET_ERR_WAIT_TIMEOUT).
func (r *CANRuntime) WaitStopped(timeoutMs uint32) (bool, error) {
	return useHandle(&r.handle, func(handle uintptr) (bool, error) {
		return nativeWait("preset_can_runtime_wait_stopped", handle, uintptr(timeoutMs))
	})
}

func (r *CANRuntime) Close() error {
	if r == nil {
		return nil
	}
	return r.handle.close("preset_can_runtime_close")
}

// Shutdown performs stop + wait_stopped + close.
func (r *CANRuntime) Shutdown(waitTimeoutMs ...uint32) error {
	if r == nil {
		return nil
	}
	if len(waitTimeoutMs) > 1 {
		return fmt.Errorf("at most one shutdown timeout is allowed")
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

type CANObserver struct{ handle nativeHandle }

func (r *CANRuntime) NewObserver(options ...CANObserverOptions) (*CANObserver, error) {
	if len(options) > 1 {
		return nil, fmt.Errorf("at most one CAN observer options value is allowed")
	}
	opts := CANObserverOptions{Capacity: 4096, RX: 1, TX: 1}
	if len(options) > 0 {
		opts = options[0]
	}
	return useHandle(&r.handle, func(runtime uintptr) (*CANObserver, error) {
		var handle uintptr
		if err := nativeStatus(
			"preset_can_observer_new",
			runtime,
			uintptr(unsafe.Pointer(&opts)),
			uintptr(unsafe.Pointer(&handle))); err != nil {
			return nil, err
		}
		return &CANObserver{handle: newNativeHandle(handle)}, nil
	})
}

func (o *CANObserver) TryRead(maxEvents int) ([]CANCaptureEvent, error) {
	if maxEvents < 1 || maxEvents > 4096 {
		return nil, fmt.Errorf("CAN observer read capacity must be in range 1-4096")
	}
	events := make([]CANCaptureEvent, maxEvents)
	var n uintptr
	err := o.handle.use(func(handle uintptr) error {
		return nativeStatus(
			"preset_can_observer_try_read",
			handle,
			uintptr(unsafe.Pointer(&events[0])),
			uintptr(maxEvents),
			uintptr(unsafe.Pointer(&n)))
	})
	if err != nil {
		return nil, err
	}
	if n > uintptr(len(events)) {
		return nil, fmt.Errorf("preset_can_observer_try_read returned invalid length %d", n)
	}
	return events[:n], nil
}

func (o *CANObserver) GetStats() (RXStats, error) {
	return useHandle(&o.handle, func(handle uintptr) (RXStats, error) {
		var stats RXStats
		if err := nativeStatus("preset_can_observer_get_stats", handle, uintptr(unsafe.Pointer(&stats))); err != nil {
			return RXStats{}, err
		}
		return stats, nil
	})
}

func (o *CANObserver) Close() error {
	if o == nil {
		return nil
	}
	return o.handle.close("preset_can_observer_close")
}

// ---- Diagnostic requests ----

type DiagnosticRequest struct{ handle nativeHandle }

// StartRequest copies the PDU and configuration, then returns without waiting
// for the ECU.
func (r *CANRuntime) StartRequest(cfg *CANDiagnosticConfig, timeouts *RequestTimeouts, pdu []byte) (*DiagnosticRequest, error) {
	if cfg == nil || timeouts == nil || len(pdu) == 0 {
		return nil, fmt.Errorf("diagnostic config, timeouts, and PDU are required")
	}
	return useHandle(&r.handle, func(runtime uintptr) (*DiagnosticRequest, error) {
		var handle uintptr
		if err := nativeStatus(
			"preset_can_request_start",
			runtime,
			uintptr(unsafe.Pointer(cfg)),
			uintptr(unsafe.Pointer(timeouts)),
			bytePtr(pdu),
			uintptr(len(pdu)),
			uintptr(unsafe.Pointer(&handle))); err != nil {
			return nil, err
		}
		return &DiagnosticRequest{handle: newNativeHandle(handle)}, nil
	})
}

func (q *DiagnosticRequest) GetInfo() (DiagnosticRequestInfo, error) {
	return useHandle(&q.handle, func(handle uintptr) (DiagnosticRequestInfo, error) {
		var info DiagnosticRequestInfo
		if err := nativeStatus("preset_request_get_info", handle, uintptr(unsafe.Pointer(&info))); err != nil {
			return DiagnosticRequestInfo{}, err
		}
		return info, nil
	})
}

func (q *DiagnosticRequest) Cancel() error {
	return q.handle.use(func(handle uintptr) error {
		return nativeStatus("preset_request_cancel", handle)
	})
}

// Wait returns true when the request reached a terminal state, false when only
// this wait expired (PRESET_ERR_WAIT_TIMEOUT).
func (q *DiagnosticRequest) Wait(timeoutMs uint32) (bool, error) {
	return useHandle(&q.handle, func(handle uintptr) (bool, error) {
		return nativeWait("preset_request_wait", handle, uintptr(timeoutMs))
	})
}

// CopyResponse returns the response bytes, or nil when the request is not
// terminal yet (PRESET_ERR_NOT_READY).
func (q *DiagnosticRequest) CopyResponse() ([]byte, error) {
	return useHandle(&q.handle, func(handle uintptr) ([]byte, error) {
		data, rc, err := callCopy("preset_request_copy_response", handle)
		if rc == PresetErrNotReady {
			return nil, nil
		}
		return data, err
	})
}

func (q *DiagnosticRequest) CopyError() (string, error) {
	return useHandle(&q.handle, func(handle uintptr) (string, error) {
		data, _, err := callCopy("preset_request_copy_error", handle)
		return string(data), err
	})
}

func (q *DiagnosticRequest) Close() error {
	if q == nil {
		return nil
	}
	return q.handle.close("preset_request_close")
}

// UDSRequest starts a CAN UDS request, waits for the result, reads the
// response and releases the request handle.
//
// A UDS negative response (7F xx NRC) is returned normally: check info.Outcome
// and info.LastNRC. Transport errors, timeouts and cancellation return an
// error together with the request info and whatever response was collected.
func (r *CANRuntime) UDSRequest(cfg *CANDiagnosticConfig, timeouts *RequestTimeouts, pdu []byte, waitMs uint32) (DiagnosticRequestInfo, []byte, error) {
	req, err := r.StartRequest(cfg, timeouts, pdu)
	if err != nil {
		return DiagnosticRequestInfo{}, nil, err
	}
	return finishRequest(req, waitMs)
}

func finishRequest(req *DiagnosticRequest, waitMs uint32) (info DiagnosticRequestInfo, response []byte, err error) {
	defer func() {
		if closeErr := req.Close(); closeErr != nil {
			if err == nil {
				err = closeErr
			} else {
				err = fmt.Errorf("%w; close diagnostic request: %v", err, closeErr)
			}
		}
	}()

	done, err := req.Wait(waitMs)
	if err != nil {
		return DiagnosticRequestInfo{}, nil, err
	}
	waitExpired := !done
	if waitExpired {
		if err := req.Cancel(); err != nil {
			return DiagnosticRequestInfo{}, nil, err
		}
		done, err = req.Wait(1000)
		if err != nil {
			return DiagnosticRequestInfo{}, nil, err
		}
		if !done {
			return DiagnosticRequestInfo{}, nil, &PresetError{
				Code:    PresetErrWaitTimeout,
				Func:    "uds_request",
				Message: "request cancellation did not finish within 1000 ms",
			}
		}
	}

	info, err = req.GetInfo()
	if err != nil {
		return DiagnosticRequestInfo{}, nil, err
	}
	response, err = req.CopyResponse()
	if err != nil {
		return info, nil, err
	}
	if waitExpired {
		return info, response, &PresetError{
			Code:    PresetErrWaitTimeout,
			Func:    "uds_request",
			Message: fmt.Sprintf("request did not finish within %d ms and was cancelled", waitMs),
		}
	}

	if (info.Outcome == PresetOutcomeError ||
		info.Outcome == PresetOutcomeTimeout ||
		info.Outcome == PresetOutcomeCancelled) && info.ErrorCode < 0 {
		text, _ := req.CopyError()
		return info, response, &PresetError{Code: info.ErrorCode, Func: "uds_request", Message: text}
	}
	return info, response, nil
}
