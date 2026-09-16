"""ABI v6 ctypes smoke test for preset_rs.dll.

Usage:
    python python/ctypes_smoke_test.py C:\\path\\to\\preset_rs.dll

The test only calls metadata/default-value functions, so no CAN/LIN hardware is
required. It also binds the CAN/LIN symbols used by ABI v6.
"""

from __future__ import annotations

import argparse
import ctypes
import os
from pathlib import Path


PRESET_OK = 0
PRESET_ERR_NOT_INIT = -3
PRESET_ABI_VERSION = 6

PRESET_CAN_BACKEND_NONE = 0
PRESET_CAN_BACKEND_TOOMOSS = 1
PRESET_CAN_BACKEND_TSMASTER = 2
PRESET_CAN_BACKEND_PCAN = 3
PRESET_CAN_BACKEND_VECTOR = 4

PRESET_LIN_PROTOCOL_13 = 0
PRESET_LIN_PROTOCOL_20 = 1
PRESET_LIN_PROTOCOL_21 = 2


class PresetCanFdTiming(ctypes.Structure):
    _fields_ = [
        ("nominal_brp", ctypes.c_uint32),
        ("nominal_tseg1", ctypes.c_uint32),
        ("nominal_tseg2", ctypes.c_uint32),
        ("nominal_sjw", ctypes.c_uint32),
        ("data_brp", ctypes.c_uint32),
        ("data_tseg1", ctypes.c_uint32),
        ("data_tseg2", ctypes.c_uint32),
        ("data_sjw", ctypes.c_uint32),
    ]


class PresetAutoConfig(ctypes.Structure):
    _fields_ = [
        ("channel", ctypes.c_uint8),
        ("is_fd", ctypes.c_uint8),
        ("brs", ctypes.c_uint8),
        ("reserved", ctypes.c_uint8),
        ("nominal_bitrate", ctypes.c_uint32),
        ("data_bitrate", ctypes.c_uint32),
        ("rx_buffer_size", ctypes.c_uint32),
        ("poll_batch_size", ctypes.c_uint32),
        ("min_tx_interval_us", ctypes.c_uint32),
        ("min_rx_poll_interval_us", ctypes.c_uint32),
        ("candidate_order", ctypes.c_uint8 * 4),
        ("reserved_tail", ctypes.c_uint8 * 4),
    ]


class PresetToomossConfig(ctypes.Structure):
    _fields_ = [
        ("channel", ctypes.c_uint8),
        ("brs", ctypes.c_uint8),
        ("reserved", ctypes.c_uint8 * 2),
        ("nominal_bitrate", ctypes.c_uint32),
        ("data_bitrate", ctypes.c_uint32),
        ("rx_buffer_size", ctypes.c_uint32),
        ("poll_batch_size", ctypes.c_uint32),
        ("min_tx_interval_us", ctypes.c_uint32),
        ("min_rx_poll_interval_us", ctypes.c_uint32),
    ]


class PresetPCANConfig(ctypes.Structure):
    _fields_ = [
        ("channel", ctypes.c_uint8),
        ("is_fd", ctypes.c_uint8),
        ("brs", ctypes.c_uint8),
        ("reserved", ctypes.c_uint8),
        ("nominal_bitrate", ctypes.c_uint32),
        ("data_bitrate", ctypes.c_uint32),
        ("rx_buffer_size", ctypes.c_uint32),
        ("poll_batch_size", ctypes.c_uint32),
        ("min_tx_interval_us", ctypes.c_uint32),
        ("min_rx_poll_interval_us", ctypes.c_uint32),
    ]


class PresetPCANLinConfig(ctypes.Structure):
    _fields_ = [
        ("channel", ctypes.c_uint8),
        ("reserved", ctypes.c_uint8),
        ("hardware_handle", ctypes.c_uint16),
        ("baudrate", ctypes.c_uint32),
    ]


class PresetTSMasterLinConfig(ctypes.Structure):
    _fields_ = [
        ("application_channel", ctypes.c_uint8),
        ("hardware_channel", ctypes.c_uint8),
        ("protocol", ctypes.c_uint8),
        ("reserved", ctypes.c_uint8),
        ("hardware_index", ctypes.c_int32),
        ("device_type", ctypes.c_int32),
        ("baudrate", ctypes.c_uint32),
    ]


class PresetVectorLinConfig(ctypes.Structure):
    _fields_ = [
        ("application_channel", ctypes.c_uint8),
        ("version", ctypes.c_uint8),
        ("reserved", ctypes.c_uint8 * 2),
        ("hardware_type", ctypes.c_int32),
        ("hardware_index", ctypes.c_uint32),
        ("hardware_channel", ctypes.c_uint32),
        ("baudrate", ctypes.c_uint32),
        ("rx_queue_size", ctypes.c_uint32),
    ]


class PresetLinFrame(ctypes.Structure):
    _fields_ = [
        ("frame_id", ctypes.c_uint8),
        ("data_len", ctypes.c_uint8),
        ("reserved", ctypes.c_uint8 * 2),
        ("data", ctypes.c_uint8 * 8),
    ]


class PresetCanFrameEx(ctypes.Structure):
    _fields_ = [
        ("id", ctypes.c_uint32),
        ("dlc", ctypes.c_uint8),
        ("data_len", ctypes.c_uint8),
        ("direction", ctypes.c_uint8),
        ("is_fd", ctypes.c_uint8),
        ("brs", ctypes.c_uint8),
        ("reserved", ctypes.c_uint8 * 7),
        ("timestamp_us", ctypes.c_uint64),
        ("data", ctypes.c_uint8 * 64),
    ]


DeviceHandle = ctypes.c_void_p
CanClientHandle = ctypes.c_void_p
LinClientHandle = ctypes.c_void_p


def bind(function: ctypes._CFuncPtr, argtypes: list[type], restype: type) -> None:
    function.argtypes = argtypes
    function.restype = restype


def load_api(path: Path) -> ctypes.CDLL:
    if os.name != "nt":
        raise RuntimeError("preset_rs.dll ctypes smoke test requires Windows")
    if not path.is_file():
        raise FileNotFoundError(path)

    # LOAD_LIBRARY_SEARCH_DLL_LOAD_DIR | LOAD_LIBRARY_SEARCH_DEFAULT_DIRS.
    api = ctypes.CDLL(str(path.resolve()), winmode=0x00000100 | 0x00001000)

    bind(api.preset_version, [], ctypes.c_char_p)
    bind(api.preset_abi_version, [], ctypes.c_uint32)
    bind(api.preset_get_capabilities, [], ctypes.c_uint64)

    bind(api.preset_auto_default_config, [], PresetAutoConfig)
    bind(api.preset_pcan_default_fd_timing, [], PresetCanFdTiming)
    bind(api.preset_toomoss_default_fd_timing, [], PresetCanFdTiming)
    bind(api.preset_pcan_lin_default_config, [], PresetPCANLinConfig)
    bind(api.preset_tsmaster_lin_default_config, [], PresetTSMasterLinConfig)
    bind(api.preset_vector_lin_default_config, [], PresetVectorLinConfig)

    bind(api.preset_auto_open, [ctypes.POINTER(PresetAutoConfig), ctypes.POINTER(DeviceHandle)], ctypes.c_int32)
    bind(api.preset_device_get_backend, [DeviceHandle, ctypes.POINTER(ctypes.c_uint8)], ctypes.c_int32)
    bind(
        api.preset_toomoss_can_init_with_timing,
        [DeviceHandle, ctypes.POINTER(PresetToomossConfig), ctypes.POINTER(PresetCanFdTiming)],
        ctypes.c_int32,
    )
    bind(
        api.preset_pcan_open_with_timing,
        [ctypes.POINTER(PresetPCANConfig), ctypes.POINTER(PresetCanFdTiming), ctypes.POINTER(DeviceHandle)],
        ctypes.c_int32,
    )
    bind(
        api.preset_pcan_can_init_with_timing,
        [DeviceHandle, ctypes.POINTER(PresetPCANConfig), ctypes.POINTER(PresetCanFdTiming)],
        ctypes.c_int32,
    )

    bind(api.preset_can_uds_set_brs, [CanClientHandle, ctypes.c_uint8], ctypes.c_int32)
    bind(api.preset_can_uds_get_brs, [CanClientHandle, ctypes.POINTER(ctypes.c_uint8)], ctypes.c_int32)
    bind(api.preset_can_uds_set_raw_tx_echo, [CanClientHandle, ctypes.c_uint8], ctypes.c_int32)
    bind(
        api.preset_can_uds_try_read_ex,
        [CanClientHandle, ctypes.POINTER(PresetCanFrameEx), ctypes.POINTER(ctypes.c_size_t)],
        ctypes.c_int32,
    )

    bind(api.preset_pcan_lin_open, [ctypes.POINTER(PresetPCANLinConfig), ctypes.POINTER(DeviceHandle)], ctypes.c_int32)
    bind(api.preset_pcan_lin_init, [DeviceHandle, ctypes.POINTER(PresetPCANLinConfig)], ctypes.c_int32)
    bind(
        api.preset_tsmaster_lin_open,
        [ctypes.POINTER(PresetTSMasterLinConfig), ctypes.POINTER(DeviceHandle)],
        ctypes.c_int32,
    )
    bind(api.preset_tsmaster_lin_init, [DeviceHandle, ctypes.POINTER(PresetTSMasterLinConfig)], ctypes.c_int32)
    bind(
        api.preset_vector_lin_open,
        [ctypes.POINTER(PresetVectorLinConfig), ctypes.POINTER(DeviceHandle)],
        ctypes.c_int32,
    )
    bind(api.preset_vector_lin_init, [DeviceHandle, ctypes.POINTER(PresetVectorLinConfig)], ctypes.c_int32)
    bind(
        api.preset_lin_master_write,
        [DeviceHandle, ctypes.c_uint8, ctypes.c_uint8, ctypes.POINTER(ctypes.c_uint8), ctypes.c_size_t],
        ctypes.c_int32,
    )
    bind(
        api.preset_lin_master_read,
        [DeviceHandle, ctypes.c_uint8, ctypes.c_uint8, ctypes.c_uint32, ctypes.POINTER(PresetLinFrame)],
        ctypes.c_int32,
    )
    bind(
        api.preset_lin_uds_client_new,
        [DeviceHandle, ctypes.c_uint8, ctypes.c_uint8, ctypes.POINTER(LinClientHandle)],
        ctypes.c_int32,
    )
    return api


def smoke_test(api: ctypes.CDLL) -> None:
    assert ctypes.sizeof(PresetCanFdTiming) == 32
    assert ctypes.sizeof(PresetAutoConfig) == 36
    assert ctypes.sizeof(PresetPCANLinConfig) == 8
    assert ctypes.sizeof(PresetTSMasterLinConfig) == 16
    assert ctypes.sizeof(PresetVectorLinConfig) == 24
    assert ctypes.sizeof(PresetLinFrame) == 12
    assert ctypes.sizeof(PresetCanFrameEx) == 88

    version = api.preset_version().decode("utf-8")
    abi_version = api.preset_abi_version()
    capabilities = api.preset_get_capabilities()
    assert abi_version == PRESET_ABI_VERSION, abi_version
    required_v6_capabilities = sum(1 << bit for bit in range(17, 25))
    assert capabilities & required_v6_capabilities == required_v6_capabilities

    auto = api.preset_auto_default_config()
    assert auto.nominal_bitrate == 500_000
    assert list(auto.candidate_order) == [
        PRESET_CAN_BACKEND_TOOMOSS,
        PRESET_CAN_BACKEND_TSMASTER,
        PRESET_CAN_BACKEND_PCAN,
        PRESET_CAN_BACKEND_VECTOR,
    ]

    pcan_timing = api.preset_pcan_default_fd_timing()
    toomoss_timing = api.preset_toomoss_default_fd_timing()
    assert (pcan_timing.nominal_brp, pcan_timing.data_brp) == (20, 4)
    assert (toomoss_timing.nominal_brp, toomoss_timing.data_tseg1) == (1, 14)

    pcan_lin = api.preset_pcan_lin_default_config()
    tsmaster_lin = api.preset_tsmaster_lin_default_config()
    vector_lin = api.preset_vector_lin_default_config()
    assert pcan_lin.baudrate == 19_200
    assert tsmaster_lin.protocol == PRESET_LIN_PROTOCOL_21
    assert vector_lin.version == 3 and vector_lin.rx_queue_size == 16_384

    # Exercise pointer arguments and the ABI v6 call signatures without
    # opening hardware. A null opaque handle must be rejected deterministically.
    backend = ctypes.c_uint8(PRESET_CAN_BACKEND_NONE)
    assert api.preset_device_get_backend(None, ctypes.byref(backend)) == PRESET_ERR_NOT_INIT
    assert api.preset_can_uds_set_brs(None, 1) == PRESET_ERR_NOT_INIT
    brs = ctypes.c_uint8()
    assert api.preset_can_uds_get_brs(None, ctypes.byref(brs)) == PRESET_ERR_NOT_INIT
    assert api.preset_can_uds_set_raw_tx_echo(None, 1) == PRESET_ERR_NOT_INIT
    frames = (PresetCanFrameEx * 1)()
    frame_count = ctypes.c_size_t(len(frames))
    assert api.preset_can_uds_try_read_ex(None, frames, ctypes.byref(frame_count)) == PRESET_ERR_NOT_INIT

    assert api.preset_toomoss_can_init_with_timing(
        None, ctypes.byref(PresetToomossConfig()), ctypes.byref(toomoss_timing)
    ) == PRESET_ERR_NOT_INIT
    assert api.preset_pcan_can_init_with_timing(
        None, ctypes.byref(PresetPCANConfig()), ctypes.byref(pcan_timing)
    ) == PRESET_ERR_NOT_INIT
    assert api.preset_pcan_lin_init(None, ctypes.byref(pcan_lin)) == PRESET_ERR_NOT_INIT
    assert api.preset_tsmaster_lin_init(None, ctypes.byref(tsmaster_lin)) == PRESET_ERR_NOT_INIT
    assert api.preset_vector_lin_init(None, ctypes.byref(vector_lin)) == PRESET_ERR_NOT_INIT

    payload = (ctypes.c_uint8 * 1)(0)
    assert api.preset_lin_master_write(None, 0, 0x3C, payload, len(payload)) == PRESET_ERR_NOT_INIT
    lin_frame = PresetLinFrame()
    assert api.preset_lin_master_read(None, 0, 0x3D, 10, ctypes.byref(lin_frame)) == PRESET_ERR_NOT_INIT
    lin_client = LinClientHandle()
    assert api.preset_lin_uds_client_new(None, 0, 0x22, ctypes.byref(lin_client)) == PRESET_ERR_NOT_INIT

    print(f"preset_rs version: {version}")
    print(f"ABI version: {abi_version}")
    print(f"capabilities: 0x{capabilities:016x}")
    print(f"AutoDriver order: {list(auto.candidate_order)}")
    print("ctypes ABI v6 smoke test passed")


def default_dll_path() -> Path:
    configured = os.environ.get("PRESET_RS_DLL")
    if configured:
        return Path(configured)
    return Path(__file__).resolve().parent.parent / "preset_rs" / "target" / "release" / "preset_rs.dll"


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("dll", nargs="?", type=Path, default=default_dll_path())
    args = parser.parse_args()
    smoke_test(load_api(args.dll))


if __name__ == "__main__":
    main()
