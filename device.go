//go:build windows && amd64

package preset

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"syscall"
	"unsafe"
)

const defaultDLLPath = `bin\preset_rs.dll`

var nativeLibrary struct {
	once   sync.Once
	mu     sync.Mutex
	path   string
	loaded bool
	handle syscall.Handle
	err    error
	procs  sync.Map
}

// SetDLLPath changes the DLL path. It must be called before any other package
// function that accesses the native library.
func SetDLLPath(path string) error {
	if path == "" {
		return fmt.Errorf("DLL path is empty")
	}
	nativeLibrary.mu.Lock()
	defer nativeLibrary.mu.Unlock()
	if nativeLibrary.loaded {
		return fmt.Errorf("native library loading has already been attempted")
	}
	nativeLibrary.path = path
	return nil
}

// ErrClosed is returned when an operation uses an already closed native handle.
var ErrClosed = errors.New("preset: native handle is closed")

type nativeHandle struct {
	state *nativeHandleState
}

type nativeHandleState struct {
	mu    sync.RWMutex
	value uintptr
}

func newNativeHandle(value uintptr) nativeHandle {
	return nativeHandle{state: &nativeHandleState{value: value}}
}

func (h *nativeHandle) use(fn func(uintptr) error) error {
	if h == nil || h.state == nil {
		return ErrClosed
	}
	h.state.mu.RLock()
	defer h.state.mu.RUnlock()
	if h.state.value == 0 {
		return ErrClosed
	}
	return fn(h.state.value)
}

func useHandle[T any](h *nativeHandle, fn func(uintptr) (T, error)) (T, error) {
	var zero T
	if h == nil || h.state == nil {
		return zero, ErrClosed
	}
	h.state.mu.RLock()
	defer h.state.mu.RUnlock()
	if h.state.value == 0 {
		return zero, ErrClosed
	}
	return fn(h.state.value)
}

func (h *nativeHandle) close(name string) error {
	if h == nil || h.state == nil {
		return nil
	}
	h.state.mu.Lock()
	defer h.state.mu.Unlock()
	if h.state.value == 0 {
		return nil
	}
	return nativeStatus(name, uintptr(unsafe.Pointer(&h.state.value)))
}

type Device struct{ handle nativeHandle }

type Toomoss struct{ Device }
type PCAN struct{ Device }
type TSMaster struct{ Device }
type Vector struct{ Device }
type AutoDevice struct{ Device }

type AutoConfig struct {
	Channel        byte
	CAN            CANChannelConfig
	CandidateOrder [4]byte
}

type PCANConfig struct {
	Channel byte
	CAN     CANChannelConfig
}

type TSMasterConfig struct {
	ApplicationChannel byte
	HardwareChannel    byte
	HardwareIndex      int32
	DeviceType         int32
	CAN                CANChannelConfig
}

type VectorConfig struct {
	ApplicationChannel byte
	HardwareType       int32
	HardwareIndex      uint32
	HardwareChannel    uint32
	NominalSjw         uint32
	NominalTseg1       uint32
	NominalTseg2       uint32
	DataSjw            uint32
	DataTseg1          uint32
	DataTseg2          uint32
	DriverRXQueueSize  uint32
	CAN                CANChannelConfig
}

func proc(name string) (uintptr, error) {
	if cached, ok := nativeLibrary.procs.Load(name); ok {
		return cached.(uintptr), nil
	}

	nativeLibrary.once.Do(func() {
		nativeLibrary.mu.Lock()
		nativeLibrary.loaded = true
		path := nativeLibrary.path
		if path == "" {
			path = defaultDLLPath
		}
		nativeLibrary.mu.Unlock()
		nativeLibrary.handle, nativeLibrary.err = syscall.LoadLibrary(path)
	})
	if nativeLibrary.err != nil {
		nativeLibrary.mu.Lock()
		path := nativeLibrary.path
		if path == "" {
			path = defaultDLLPath
		}
		nativeLibrary.mu.Unlock()
		return 0, fmt.Errorf("load %s: %w", path, nativeLibrary.err)
	}

	p, err := syscall.GetProcAddress(nativeLibrary.handle, name)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", name, err)
	}
	nativeLibrary.procs.Store(name, p)
	return p, nil
}

func nativeResult(name string, args ...uintptr) (uintptr, error) {
	p, err := proc(name)
	if err != nil {
		return 0, err
	}
	return syscallR1(p, args...), nil
}

func nativeStatus(name string, args ...uintptr) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	rc, err := nativeResult(name, args...)
	if err != nil {
		return err
	}
	return presetStatus(name, rc)
}

func openNamedDevice(name string) (uintptr, error) {
	if err := checkABI(); err != nil {
		return 0, err
	}
	var handle uintptr
	if err := nativeStatus(name, uintptr(unsafe.Pointer(&handle))); err != nil {
		return 0, err
	}
	return handle, nil
}

// Native struct-return helpers use a hidden output pointer on Windows amd64.
func defaultStruct[T any](name string) (T, error) {
	var result T
	p, err := proc(name)
	if err != nil {
		return result, err
	}
	if unsafe.Sizeof(result) > 8 {
		syscallR1(p, uintptr(unsafe.Pointer(&result)))
	} else {
		copyRegisterResult(unsafe.Pointer(&result), unsafe.Sizeof(result), syscallR1(p))
	}
	return result, nil
}

// syscallR1 exposes only the first return register of a native call.
func syscallR1(fn uintptr, args ...uintptr) uintptr {
	r1, _, _ := syscall.SyscallN(fn, args...)
	return r1
}

func NewPCAN() (*PCAN, error) {
	h, err := openNamedDevice("preset_pcan_open")
	if err != nil {
		return nil, err
	}
	return &PCAN{Device{handle: newNativeHandle(h)}}, nil
}

func NewToomoss() (*Toomoss, error) {
	h, err := openNamedDevice("preset_toomoss_open")
	if err != nil {
		return nil, err
	}
	return &Toomoss{Device{handle: newNativeHandle(h)}}, nil
}

func NewTSMaster() (*TSMaster, error) {
	h, err := openNamedDevice("preset_tsmaster_open")
	if err != nil {
		return nil, err
	}
	return &TSMaster{Device{handle: newNativeHandle(h)}}, nil
}

func NewVector() (*Vector, error) {
	h, err := openNamedDevice("preset_vector_open")
	if err != nil {
		return nil, err
	}
	return &Vector{Device{handle: newNativeHandle(h)}}, nil
}

// Close releases the physical device. Close all borrowed handles first.
func (t *Device) Close() error {
	if t == nil {
		return nil
	}
	return t.handle.close("preset_device_close")
}

func DefaultAutoConfig() (AutoConfig, error) {
	return defaultStruct[AutoConfig]("preset_auto_default_config")
}
func DefaultPCANConfig() (PCANConfig, error) {
	return defaultStruct[PCANConfig]("preset_pcan_default_config")
}
func DefaultTSMasterConfig() (TSMasterConfig, error) {
	return defaultStruct[TSMasterConfig]("preset_tsmaster_default_config")
}
func DefaultVectorConfig() (VectorConfig, error) {
	return defaultStruct[VectorConfig]("preset_vector_default_config")
}
func DefaultPCANFDTiming() (CANFDTiming, error) {
	return defaultStruct[CANFDTiming]("preset_pcan_default_fd_timing")
}

// NewAuto selects an installed backend in the configured priority order.
func NewAuto(cfg *AutoConfig) (*AutoDevice, error) {
	if cfg == nil {
		return nil, fmt.Errorf("auto config is nil")
	}
	if err := checkABI(); err != nil {
		return nil, err
	}
	var h uintptr
	if err := nativeStatus("preset_auto_open", uintptr(unsafe.Pointer(cfg)), uintptr(unsafe.Pointer(&h))); err != nil {
		return nil, err
	}
	return &AutoDevice{Device{handle: newNativeHandle(h)}}, nil
}

func (t *PCAN) InitCANWithConfig(cfg *PCANConfig, timing *CANFDTiming) error {
	if cfg == nil {
		return fmt.Errorf("PCAN config is nil")
	}
	return t.handle.use(func(handle uintptr) error {
		return nativeStatus("preset_pcan_can_init", handle, uintptr(unsafe.Pointer(cfg)), uintptr(unsafe.Pointer(timing)))
	})
}

func (t *TSMaster) InitCANWithConfig(cfg *TSMasterConfig) error {
	if cfg == nil {
		return fmt.Errorf("TSMaster config is nil")
	}
	return t.handle.use(func(handle uintptr) error {
		return nativeStatus("preset_tsmaster_can_init", handle, uintptr(unsafe.Pointer(cfg)))
	})
}

func (t *Vector) InitCANWithConfig(cfg *VectorConfig) error {
	if cfg == nil {
		return fmt.Errorf("Vector config is nil")
	}
	return t.handle.use(func(handle uintptr) error {
		return nativeStatus("preset_vector_can_init", handle, uintptr(unsafe.Pointer(cfg)))
	})
}

func (t *PCAN) InitCAN(channels ...byte) error {
	for _, ch := range channels {
		cfg, err := DefaultPCANConfig()
		if err != nil {
			return err
		}
		cfg.Channel = ch
		if err := t.InitCANWithConfig(&cfg, nil); err != nil {
			return err
		}
	}
	return nil
}

func (t *TSMaster) InitCAN(channels ...byte) error {
	for _, ch := range channels {
		cfg, err := DefaultTSMasterConfig()
		if err != nil {
			return err
		}
		cfg.ApplicationChannel, cfg.HardwareChannel = ch, ch
		if err := t.InitCANWithConfig(&cfg); err != nil {
			return err
		}
	}
	return nil
}

func (t *Vector) InitCAN(channels ...byte) error {
	for _, ch := range channels {
		cfg, err := DefaultVectorConfig()
		if err != nil {
			return err
		}
		cfg.ApplicationChannel, cfg.HardwareChannel = ch, uint32(ch)
		if err := t.InitCANWithConfig(&cfg); err != nil {
			return err
		}
	}
	return nil
}
