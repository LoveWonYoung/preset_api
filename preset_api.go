// Package preset_api provides a pure-Go Windows binding for preset_rs.dll.
//
// It deliberately uses package syscall instead of cgo. The DLL is loaded on
// the first API call. Applications that keep the DLL outside the normal
// Windows DLL search path should call LoadDLL before using any other function.
package preset_api

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"syscall"
	"unsafe"
)

const defaultDLLName = "preset_rs.dll"

var dllState struct {
	sync.Mutex
	dll     *syscall.DLL
	path    string
	loadErr error
	procs   map[string]*syscall.Proc
}

// LoadDLL loads preset_rs from path. It is optional when preset_rs.dll is on
// the Windows DLL search path. A process may bind to only one DLL path.
func LoadDLL(path string) error {
	if path == "" {
		return errors.New("preset_api: DLL path is empty")
	}

	dllState.Lock()
	if dllState.dll != nil {
		loadedPath := dllState.path
		dllState.Unlock()
		if loadedPath == path {
			return nil
		}
		return fmt.Errorf("preset_api: DLL already loaded from %q", loadedPath)
	}
	dllState.Unlock()

	dll, err := syscall.LoadDLL(path)

	dllState.Lock()
	defer dllState.Unlock()
	if dllState.dll != nil { // Another goroutine won the load race.
		if dll != nil {
			_ = dll.Release()
		}
		if dllState.path == path {
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
	dllState.procs = make(map[string]*syscall.Proc)
	return nil
}

// DLLLoadError reports the most recent DLL loading or symbol lookup error.
func DLLLoadError() error {
	dllState.Lock()
	defer dllState.Unlock()
	return dllState.loadErr
}

func findProc(name string) (*syscall.Proc, error) {
	dllState.Lock()
	if dllState.dll != nil {
		if proc := dllState.procs[name]; proc != nil {
			dllState.Unlock()
			return proc, nil
		}
		dll := dllState.dll
		dllState.Unlock()
		proc, err := dll.FindProc(name)
		if err != nil {
			dllState.Lock()
			dllState.loadErr = fmt.Errorf("preset_api: find symbol %q: %w", name, err)
			dllState.Unlock()
			return nil, err
		}
		dllState.Lock()
		if existing := dllState.procs[name]; existing != nil {
			proc = existing
		} else {
			dllState.procs[name] = proc
		}
		dllState.Unlock()
		return proc, nil
	}
	dllState.Unlock()

	if err := LoadDLL(defaultDLLName); err != nil {
		return nil, err
	}
	return findProc(name)
}

func invoke(name string, args ...uintptr) (uintptr, bool) {
	proc, err := findProc(name)
	if err != nil {
		return 0, false
	}
	r1, _, _ := proc.Call(args...)
	return r1, true
}

func invokeStatus(name string, args ...uintptr) int32 {
	r1, ok := invoke(name, args...)
	if !ok {
		return PRESET_ERR_TRANSPORT
	}
	return int32(uint32(r1))
}

func pointer[T any](value *T) uintptr {
	return uintptr(unsafe.Pointer(value))
}

func slicePointer[T any](values []T) uintptr {
	if len(values) == 0 {
		return 0
	}
	return uintptr(unsafe.Pointer(&values[0]))
}

// Version returns the preset_rs semantic version, or an empty string when the
// DLL could not be loaded. DLLLoadError provides the load failure details.
func Version() string {
	address, ok := invoke("preset_version")
	if !ok || address == 0 {
		return ""
	}
	const maxVersionLength = 4096
	length := 0
	for length < maxVersionLength && *(*byte)(unsafe.Pointer(address + uintptr(length))) != 0 {
		length++
	}
	if length == maxVersionLength {
		return ""
	}
	return unsafe.String((*byte)(unsafe.Pointer(address)), length)
}

// ABIVersion returns the exported C ABI version.
func ABIVersion() uint32 {
	value, _ := invoke("preset_abi_version")
	return uint32(value)
}

// Capabilities returns the PRESET_CAP_* bit set advertised by the DLL.
func Capabilities() uint64 {
	value, _ := invoke("preset_get_capabilities")
	return uint64(value)
}

// Structs larger than eight bytes use the Microsoft x64 hidden return-buffer
// parameter. PresetToomossLinConfig is exactly eight bytes and is returned in
// RAX instead.
func callStruct(name string, destination unsafe.Pointer) bool {
	_, ok := invoke(name, uintptr(destination))
	runtime.KeepAlive(destination)
	return ok
}

func DefaultConfig() (config PresetConfig) {
	callStruct("preset_default_config", unsafe.Pointer(&config))
	return
}

func ToomossDefaultConfig() (config PresetToomossConfig) {
	callStruct("preset_toomoss_default_config", unsafe.Pointer(&config))
	return
}

func ToomossLinDefaultConfig() (config PresetToomossLinConfig) {
	value, ok := invoke("preset_toomoss_lin_default_config")
	if ok {
		*(*uint64)(unsafe.Pointer(&config)) = uint64(value)
	}
	return
}

func ToomossElinsDefaultConfig() (config PresetToomossElinsConfig) {
	callStruct("preset_toomoss_elins_default_config", unsafe.Pointer(&config))
	return
}

func PcanDefaultConfig() (config PresetPCANConfig) {
	callStruct("preset_pcan_default_config", unsafe.Pointer(&config))
	return
}

func TsmasterDefaultConfig() (config PresetTSMasterConfig) {
	callStruct("preset_tsmaster_default_config", unsafe.Pointer(&config))
	return
}

func VectorDefaultConfig() (config PresetVectorConfig) {
	callStruct("preset_vector_default_config", unsafe.Pointer(&config))
	return
}
