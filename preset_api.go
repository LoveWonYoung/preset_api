//go:build windows && amd64

// Package preset_api provides a pure-Go Windows AMD64 binding for preset_rs.dll.
//
// It deliberately uses package syscall instead of cgo. Applications must call
// LoadDLL with an explicit path before using the API.
package preset_api

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"unsafe"
)

const maxDLLArguments = 8

const (
	loadLibrarySearchDLLLoadDir  = 0x00000100
	loadLibrarySearchDefaultDirs = 0x00001000
)

var loadLibraryExW = syscall.NewLazyDLL("kernel32.dll").NewProc("LoadLibraryExW")

// callArgument keeps Go pointers as pointers until the final syscall boundary.
// This lets invoke pin them for the duration of Proc.Call instead of relying on
// uintptr values surviving across helper calls.
type callArgument struct {
	value   uintptr
	pointer unsafe.Pointer
}

type callResult struct {
	r1 uintptr
	r2 uintptr
}

type procLookup struct {
	proc *syscall.Proc
	err  error
}

var dllState struct {
	sync.Mutex
	dll     *syscall.DLL
	path    string
	loadErr error
}

var procCache sync.Map // map[string]procLookup; DLL is immutable after load.

// LoadDLL loads preset_rs from an explicit path. Relative paths are resolved
// against the current working directory before LoadLibrary is called, avoiding
// the ambient Windows DLL search path. A process may bind to only one DLL.
func LoadDLL(path string) error {
	if path == "" {
		return errors.New("preset_api: DLL path is empty")
	}
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("preset_api: resolve DLL path %q: %w", path, err)
	}
	return loadDLL(filepath.Clean(absolutePath))
}

func loadDLL(path string) error {
	dllState.Lock()
	if dllState.dll != nil {
		loadedPath := dllState.path
		dllState.Unlock()
		if strings.EqualFold(loadedPath, path) {
			return nil
		}
		return fmt.Errorf("preset_api: DLL already loaded from %q", loadedPath)
	}
	dllState.Unlock()

	dll, err := secureLoadDLL(path)
	if err == nil {
		err = validateDLL(dll)
		if err != nil {
			_ = dll.Release()
			dll = nil
		}
	}

	dllState.Lock()
	defer dllState.Unlock()
	if dllState.dll != nil { // Another goroutine won the load race.
		if dll != nil {
			_ = dll.Release()
		}
		if strings.EqualFold(dllState.path, path) {
			return nil
		}
		return fmt.Errorf("preset_api: DLL already loaded from %q", dllState.path)
	}
	if err != nil {
		dllState.loadErr = fmt.Errorf("preset_api: load %q: %w", path, err)
		return dllState.loadErr
	}
	dllState.dll = dll
	dllState.path = path
	dllState.loadErr = nil
	return nil
}

func secureLoadDLL(path string) (*syscall.DLL, error) {
	pathPointer, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return nil, err
	}
	var pinner runtime.Pinner
	pinner.Pin(pathPointer)
	defer pinner.Unpin()
	handle, _, callErr := loadLibraryExW.Call(
		uintptr(unsafe.Pointer(pathPointer)),
		0,
		loadLibrarySearchDLLLoadDir|loadLibrarySearchDefaultDirs,
	)
	runtime.KeepAlive(pathPointer)
	if handle == 0 {
		return nil, fmt.Errorf("LoadLibraryExW: %w", callErr)
	}
	return &syscall.DLL{Name: path, Handle: syscall.Handle(handle)}, nil
}

func validateDLL(dll *syscall.DLL) error {
	proc, err := dll.FindProc("preset_abi_version")
	if err != nil {
		return fmt.Errorf("find preset_abi_version: %w", err)
	}
	value, _, _ := proc.Call()
	version := uint32(value)
	if version != SupportedABIVersion {
		return fmt.Errorf("unsupported ABI version %d (want %d)", version, SupportedABIVersion)
	}
	return nil
}

// DLLLoadError reports the most recent DLL loading or symbol lookup error.
func DLLLoadError() error {
	dllState.Lock()
	defer dllState.Unlock()
	return dllState.loadErr
}

func findProc(name string) (*syscall.Proc, error) {
	if cached, ok := procCache.Load(name); ok {
		lookup := cached.(procLookup)
		return lookup.proc, lookup.err
	}

	dllState.Lock()
	dll := dllState.dll
	dllState.Unlock()
	if dll == nil {
		err := errors.New("preset_api: DLL is not loaded; call LoadDLL first")
		dllState.Lock()
		dllState.loadErr = err
		dllState.Unlock()
		return nil, err
	}

	proc, err := dll.FindProc(name)
	lookup := procLookup{proc: proc}
	if err != nil {
		lookup.proc = nil
		lookup.err = fmt.Errorf("preset_api: find symbol %q: %w", name, err)
		dllState.Lock()
		dllState.loadErr = lookup.err
		dllState.Unlock()
	}
	actual, _ := procCache.LoadOrStore(name, lookup)
	lookup = actual.(procLookup)
	return lookup.proc, lookup.err
}

func invoke(name string, args ...callArgument) (callResult, error) {
	proc, err := findProc(name)
	if err != nil {
		return callResult{}, err
	}
	if len(args) > maxDLLArguments {
		err := fmt.Errorf("preset_api: call %q has %d arguments; maximum is %d", name, len(args), maxDLLArguments)
		dllState.Lock()
		dllState.loadErr = err
		dllState.Unlock()
		return callResult{}, err
	}

	var values [maxDLLArguments]uintptr
	var pinner runtime.Pinner
	pinned := false
	for index, argument := range args {
		if argument.pointer != nil {
			pinner.Pin(argument.pointer)
			if !pinned {
				pinned = true
				defer pinner.Unpin()
			}
			values[index] = uintptr(argument.pointer)
		} else {
			values[index] = argument.value
		}
	}

	r1, r2, _ := proc.Call(values[:len(args)]...)
	runtime.KeepAlive(args)
	return callResult{r1: r1, r2: r2}, nil
}

func invokeStatus(name string, args ...callArgument) int32 {
	result, err := invoke(name, args...)
	if err != nil {
		return PRESET_ERR_TRANSPORT
	}
	return int32(uint32(result.r1))
}

func word(value uintptr) callArgument {
	return callArgument{value: value}
}

func pointer[T any](value *T) callArgument {
	return rawPointer(unsafe.Pointer(value))
}

func rawPointer(value unsafe.Pointer) callArgument {
	return callArgument{pointer: value}
}

func slicePointer[T any](values []T) callArgument {
	if len(values) == 0 {
		return callArgument{}
	}
	return pointer(&values[0])
}

// Version returns the preset_rs semantic version.
//
//go:nocheckptr
func Version() (string, error) {
	result, err := invoke("preset_version")
	address := result.r1
	if err != nil {
		return "", err
	}
	if address == 0 {
		return "", errors.New("preset_api: preset_version returned null")
	}
	const maxVersionLength = 4096
	bytes := unsafe.Slice((*byte)(unsafe.Pointer(address)), maxVersionLength)
	length := 0
	for length < maxVersionLength && bytes[length] != 0 {
		length++
	}
	if length == maxVersionLength {
		return "", errors.New("preset_api: preset_version is not null-terminated")
	}
	return string(bytes[:length]), nil
}

// ABIVersion returns the exported C ABI version.
func ABIVersion() (uint32, error) {
	result, err := invoke("preset_abi_version")
	return uint32(result.r1), err
}

// Capabilities returns the PRESET_CAP_* bit set advertised by the DLL.
func Capabilities() (uint64, error) {
	result, err := invoke("preset_get_capabilities")
	return uint64(result.r1), err
}

func DefaultConfig() (config PresetConfig, err error) {
	err = callStruct("preset_default_config", unsafe.Pointer(&config), unsafe.Sizeof(config))
	return config, err
}

func TpFrameDefaultConfig() (config PresetTpFrameConfig, err error) {
	err = callStruct("preset_tp_frame_default_config", unsafe.Pointer(&config), unsafe.Sizeof(config))
	return config, err
}

func AutoDefaultConfig() (config PresetAutoConfig, err error) {
	err = callStruct("preset_auto_default_config", unsafe.Pointer(&config), unsafe.Sizeof(config))
	return config, err
}

func PcanDefaultFDTiming() (timing PresetCanFdTiming, err error) {
	err = callStruct("preset_pcan_default_fd_timing", unsafe.Pointer(&timing), unsafe.Sizeof(timing))
	return timing, err
}

func ToomossDefaultFDTiming() (timing PresetCanFdTiming, err error) {
	err = callStruct("preset_toomoss_default_fd_timing", unsafe.Pointer(&timing), unsafe.Sizeof(timing))
	return timing, err
}

func ToomossDefaultConfig() (config PresetToomossConfig, err error) {
	err = callStruct("preset_toomoss_default_config", unsafe.Pointer(&config), unsafe.Sizeof(config))
	return config, err
}

func ToomossLinDefaultConfig() (config PresetToomossLinConfig, err error) {
	err = callStruct("preset_toomoss_lin_default_config", unsafe.Pointer(&config), unsafe.Sizeof(config))
	return config, err
}

func ToomossElinsDefaultConfig() (config PresetToomossElinsConfig, err error) {
	err = callStruct("preset_toomoss_elins_default_config", unsafe.Pointer(&config), unsafe.Sizeof(config))
	return config, err
}

func PcanDefaultConfig() (config PresetPCANConfig, err error) {
	err = callStruct("preset_pcan_default_config", unsafe.Pointer(&config), unsafe.Sizeof(config))
	return config, err
}

func TsmasterDefaultConfig() (config PresetTSMasterConfig, err error) {
	err = callStruct("preset_tsmaster_default_config", unsafe.Pointer(&config), unsafe.Sizeof(config))
	return config, err
}

func VectorDefaultConfig() (config PresetVectorConfig, err error) {
	err = callStruct("preset_vector_default_config", unsafe.Pointer(&config), unsafe.Sizeof(config))
	return config, err
}

func PcanLinDefaultConfig() (config PresetPCANLinConfig, err error) {
	err = callStruct("preset_pcan_lin_default_config", unsafe.Pointer(&config), unsafe.Sizeof(config))
	return config, err
}

func TsmasterLinDefaultConfig() (config PresetTSMasterLinConfig, err error) {
	err = callStruct("preset_tsmaster_lin_default_config", unsafe.Pointer(&config), unsafe.Sizeof(config))
	return config, err
}

func VectorLinDefaultConfig() (config PresetVectorLinConfig, err error) {
	err = callStruct("preset_vector_lin_default_config", unsafe.Pointer(&config), unsafe.Sizeof(config))
	return config, err
}
