"""Python ctypes binding for the ``preset_rs`` ABI v6 Windows DLL.

The public classes in this module own their native handles and can be used as
context managers.  The DLL itself remains available as ``sdk.dll`` for callers
that need the status-code-oriented C API directly.

Typical usage::

    from preset_rs_native_sdk import PresetRSNativeSDK

    sdk = PresetRSNativeSDK(r"C:\\path\\to\\preset_rs.dll")
    config = sdk.pcan_default_config()
    with sdk.pcan_open(config) as device:
        uds_config = sdk.default_config()
        with device.new_can_uds_client(config.channel, uds_config) as client:
            response = client.request(bytes.fromhex("22 F1 90"), timeout_ms=1000)
            print(response.hex(" "))

All high-level operations raise :class:`PresetError` on a non-zero native
status.  Native return codes and every ABI structure remain exported for code
that prefers to call ``sdk.dll.preset_*`` directly.
"""

from __future__ import annotations

import ctypes
import os
import threading
from pathlib import Path
from typing import Iterable, Sequence


SUPPORTED_ABI_VERSION = 6

PRESET_OK = 0
PRESET_ERR_NULL_PTR = -1
PRESET_ERR_BUFFER_TOO_SMALL = -2
PRESET_ERR_NOT_INIT = -3
PRESET_ERR_TIMEOUT = -4
PRESET_ERR_UDS_NEGATIVE = -5
PRESET_ERR_TRANSPORT = -6
PRESET_ERR_UNSUPPORTED = -7
PRESET_ERR_INVALID_ARG = -8
PRESET_ERR_BUSY = -9
PRESET_ERR_ISOTP = -10
PRESET_ERR_PANIC = -11

PRESET_CAP_CLASSIC_CAN = 1 << 0
PRESET_CAP_CAN_FD = 1 << 1
PRESET_CAP_FUNCTIONAL = 1 << 2
PRESET_CAP_RAW_CAN = 1 << 3
PRESET_CAP_STANDARD_ID_ONLY = 1 << 4
PRESET_CAP_RUNTIME_FLOW_CONTROL = 1 << 5
PRESET_CAP_TOOMOSS_SHARED_DEVICE = 1 << 6
PRESET_CAP_TOOMOSS_LIN = 1 << 7
PRESET_CAP_TOOMOSS = 1 << 8
PRESET_CAP_TSMASTER = 1 << 9
PRESET_CAP_PCAN = 1 << 10
PRESET_CAP_VECTOR = 1 << 11
PRESET_CAP_TOOMOSS_ELINS = 1 << 12
PRESET_CAP_LIN_UDS = 1 << 13
PRESET_CAP_FILE_LOGGING = 1 << 14
PRESET_CAP_CAN_CHANNEL_SELECT = 1 << 15
PRESET_CAP_BUS_LOAD = 1 << 16
PRESET_CAP_RAW_CAN_METADATA = 1 << 17
PRESET_CAP_RUNTIME_BRS = 1 << 18
PRESET_CAP_EXPLICIT_FD_TIMING = 1 << 19
PRESET_CAP_AUTO_DRIVER = 1 << 20
PRESET_CAP_PCAN_LIN = 1 << 21
PRESET_CAP_TSMASTER_LIN = 1 << 22
PRESET_CAP_VECTOR_LIN = 1 << 23
PRESET_CAP_RAW_CAN_TIMESTAMP = 1 << 24
PRESET_CAP_MANUAL_ISOTP_FRAMING = 1 << 25
PRESET_CAP_DEVICE_RAW_CAN = 1 << 26
PRESET_CAP_MANUAL_TP_MODE = 1 << 27
PRESET_CAP_TP_FRAME_BUILDER = 1 << 28

PRESET_CAN_DIRECTION_TX = 0
PRESET_CAN_DIRECTION_RX = 1

PRESET_TP_FLOW_CONTINUE_TO_SEND = 0
PRESET_TP_FLOW_WAIT = 1
PRESET_TP_FLOW_OVERFLOW = 2

PRESET_CAN_BACKEND_NONE = 0
PRESET_CAN_BACKEND_TOOMOSS = 1
PRESET_CAN_BACKEND_TSMASTER = 2
PRESET_CAN_BACKEND_PCAN = 3
PRESET_CAN_BACKEND_VECTOR = 4

PRESET_LOG_FILTER_OFF = 0
PRESET_LOG_FILTER_LIST = 1

PRESET_TOOMOSS_CHANNEL_1 = 1 << 0
PRESET_TOOMOSS_CHANNEL_2 = 1 << 1
PRESET_TOOMOSS_CHANNEL_3 = 1 << 2
PRESET_TOOMOSS_CHANNEL_4 = 1 << 3

PRESET_TOOMOSS_LIN_POWER_0V = 0
PRESET_TOOMOSS_LIN_POWER_12V = 1
PRESET_TOOMOSS_LIN_POWER_5V = 2

PRESET_TOOMOSS_ELINS_VER_IND83080 = 0
PRESET_TOOMOSS_ELINS_VER_IND83220 = 1
PRESET_TOOMOSS_ELINS_VER_IND83010 = 2

PRESET_LIN_PROTOCOL_13 = 0
PRESET_LIN_PROTOCOL_20 = 1
PRESET_LIN_PROTOCOL_21 = 2
PRESET_LIN_FUNCTIONAL_NAD = 0x7E
PRESET_LIN_BROADCAST_NAD = 0x7F

# Vector XL hardware type values from vxlapi.h.
PRESET_VECTOR_HWTYPE_ANY = -1
PRESET_VECTOR_HWTYPE_NONE = 0
PRESET_VECTOR_HWTYPE_VIRTUAL = 1
PRESET_VECTOR_HWTYPE_CANCARDX = 2
PRESET_VECTOR_HWTYPE_CANAC2PCI = 6
PRESET_VECTOR_HWTYPE_CANCARDY = 12
PRESET_VECTOR_HWTYPE_CANCARDXL = 15
PRESET_VECTOR_HWTYPE_CANCASEXL = 21
PRESET_VECTOR_HWTYPE_CANCASEXL_LOG_OBSOLETE = 23
PRESET_VECTOR_HWTYPE_CANBOARDXL = 25
PRESET_VECTOR_HWTYPE_CANBOARDXL_PXI = 27
PRESET_VECTOR_HWTYPE_VN2600 = 29
PRESET_VECTOR_HWTYPE_VN2610 = 29
PRESET_VECTOR_HWTYPE_VN3300 = 37
PRESET_VECTOR_HWTYPE_VN3600 = 39
PRESET_VECTOR_HWTYPE_VN7600 = 41
PRESET_VECTOR_HWTYPE_CANCARDXLE = 43
PRESET_VECTOR_HWTYPE_VN8900 = 45
PRESET_VECTOR_HWTYPE_VN8950 = 47
PRESET_VECTOR_HWTYPE_VN2640 = 53
PRESET_VECTOR_HWTYPE_VN1610 = 55
PRESET_VECTOR_HWTYPE_VN1630 = 57
PRESET_VECTOR_HWTYPE_VN1640 = 59
PRESET_VECTOR_HWTYPE_VN8970 = 61
PRESET_VECTOR_HWTYPE_VN1611 = 63
PRESET_VECTOR_HWTYPE_VN5240 = 64
PRESET_VECTOR_HWTYPE_VN5610 = 65
PRESET_VECTOR_HWTYPE_VN5620 = 66
PRESET_VECTOR_HWTYPE_VN7570 = 67
PRESET_VECTOR_HWTYPE_VN5650 = 68
PRESET_VECTOR_HWTYPE_IPCLIENT = 69
PRESET_VECTOR_HWTYPE_IPSERVER = 71
PRESET_VECTOR_HWTYPE_VX1121 = 73
PRESET_VECTOR_HWTYPE_VX1131 = 75
PRESET_VECTOR_HWTYPE_VT6204 = 77
PRESET_VECTOR_HWTYPE_VN1630_LOG = 79
PRESET_VECTOR_HWTYPE_VN7610 = 81
PRESET_VECTOR_HWTYPE_VN7572 = 83
PRESET_VECTOR_HWTYPE_VN8972 = 85
PRESET_VECTOR_HWTYPE_VN0601 = 87
PRESET_VECTOR_HWTYPE_VN5640 = 89
PRESET_VECTOR_HWTYPE_VX0312 = 91
PRESET_VECTOR_HWTYPE_VH6501 = 94
PRESET_VECTOR_HWTYPE_VN8800 = 95
PRESET_VECTOR_HWTYPE_IPCL8800 = 96
PRESET_VECTOR_HWTYPE_IPSRV8800 = 97
PRESET_VECTOR_HWTYPE_CSMCAN = 98
PRESET_VECTOR_HWTYPE_VN5610A = 101
PRESET_VECTOR_HWTYPE_VN7640 = 102
PRESET_VECTOR_HWTYPE_VX1135 = 104
PRESET_VECTOR_HWTYPE_VN4610 = 105
PRESET_VECTOR_HWTYPE_VT6306 = 107
PRESET_VECTOR_HWTYPE_VT6104A = 108
PRESET_VECTOR_HWTYPE_VN5430 = 109
PRESET_VECTOR_HWTYPE_VTSSERVICE = 110
PRESET_VECTOR_HWTYPE_VN1530 = 112
PRESET_VECTOR_HWTYPE_VN1531 = 113
PRESET_VECTOR_HWTYPE_VX1161A = 114
PRESET_VECTOR_HWTYPE_VX1161B = 115

# TSMaster device subtype values.  Values are deliberately written out so the
# Python constants cannot shift if a name is inserted later.
_TSMASTER_DEVICE_NAMES = (
    "TSCAN_PRO", "TSCAN_LITE1", "TC1001", "TL1001", "TC1011", "TM5011",
    "TC1002", "TC1014", "TSCANFD2517", "TC1026", "TC1016", "TC1012",
    "TC1013", "TLOG1002", "TC1034", "TC1018", "GW2116", "TC2115",
    "MP1013", "TC1113", "TC1114", "TP1013", "TC1017", "TP1018",
    "TF10XX", "TL1004_FD_4_LIN_2", "TE1051", "TP1051", "TP1034",
    "TTS9015", "TP1026", "TTS1026", "TTS1034", "TTS1018", "TL1011",
    "TTS1015_LIAUTO", "TTS1013_LIAUTO", "TTS1016PRO", "TC1054PRO",
    "TC1054", "TLOG1038", "TO1013", "TC1034PRO", "TC1018PRO",
    "TC1038PRO", "TC1014PRO", "TC1034PROPLUS", "TA1038", "TC1055PRO",
    "TC1056PRO", "TC1057PRO", "TC4016", "GW2208", "TLOG1039", "GW1040",
    "TC3014", "TP1014", "TA825_4", "TC1013HV", "TC1052", "TTS1017PRO",
    "TLOG1057", "TC1017PRO", "GW2202", "GW2204", "GW2212", "TA821",
    "TX1000", "TC1055PROPLUS", "TC1043",
)
globals().update(
    {f"PRESET_TSMASTER_{name}": value for value, name in enumerate(_TSMASTER_DEVICE_NAMES, 1)}
)
del _TSMASTER_DEVICE_NAMES


class _PresetStructure(ctypes.Structure):
    """Base structure with a compact diagnostic representation."""

    def __repr__(self) -> str:
        values = ", ".join(f"{name}={getattr(self, name)!r}" for name, *_ in self._fields_)
        return f"{type(self).__name__}({values})"


class PresetConfig(_PresetStructure):
    _fields_ = [
        ("physical_id", ctypes.c_uint32),
        ("response_id", ctypes.c_uint32),
        ("functional_id", ctypes.c_uint32),
        ("n_bs_ms", ctypes.c_uint32),
        ("n_cr_ms", ctypes.c_uint32),
        ("st_min_us", ctypes.c_uint32),
        ("max_pdu_len", ctypes.c_uint32),
        ("rx_capacity", ctypes.c_uint32),
        ("is_fd", ctypes.c_uint8),
        ("padding_enabled", ctypes.c_uint8),
        ("padding_byte", ctypes.c_uint8),
        ("block_size", ctypes.c_uint8),
        ("max_wait_frames", ctypes.c_uint8),
        ("raw_rx_enabled", ctypes.c_uint8),
        ("reserved", ctypes.c_uint8 * 2),
    ]


class PresetTpFrameConfig(_PresetStructure):
    _fields_ = [
        ("is_fd", ctypes.c_uint8),
        ("padding_enabled", ctypes.c_uint8),
        ("padding_byte", ctypes.c_uint8),
        ("reserved", ctypes.c_uint8),
    ]


class PresetTpEncodedFrame(_PresetStructure):
    _fields_ = [
        ("data_len", ctypes.c_uint8),
        ("reserved", ctypes.c_uint8 * 7),
        ("data", ctypes.c_uint8 * 64),
    ]

    @property
    def payload(self) -> bytes:
        return bytes(self.data[: self.data_len])


class PresetCanFdTiming(_PresetStructure):
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


class PresetAutoConfig(_PresetStructure):
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


class PresetToomossConfig(_PresetStructure):
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


class PresetToomossLinConfig(_PresetStructure):
    _fields_ = [
        ("channels_mask", ctypes.c_uint8),
        ("master_mode", ctypes.c_uint8),
        ("reserved", ctypes.c_uint8 * 2),
        ("baudrate", ctypes.c_uint32),
    ]


class PresetToomossElinsConfig(_PresetStructure):
    _fields_ = [
        ("channels_mask", ctypes.c_uint8),
        ("resistor_enabled", ctypes.c_uint8),
        ("version", ctypes.c_uint8),
        ("reserved", ctypes.c_uint8),
        ("baudrate", ctypes.c_uint32),
        ("receive_timeout_us", ctypes.c_uint32),
    ]


class PresetToomossLinFrame(_PresetStructure):
    _fields_ = [
        ("timestamp", ctypes.c_uint32),
        ("frame_id", ctypes.c_uint8),
        ("data_len", ctypes.c_uint8),
        ("classic_checksum", ctypes.c_uint8),
        ("reserved", ctypes.c_uint8),
        ("data", ctypes.c_uint8 * 8),
    ]

    @property
    def payload(self) -> bytes:
        return bytes(self.data[: self.data_len])


class PresetPCANLinConfig(_PresetStructure):
    _fields_ = [
        ("channel", ctypes.c_uint8),
        ("reserved", ctypes.c_uint8),
        ("hardware_handle", ctypes.c_uint16),
        ("baudrate", ctypes.c_uint32),
    ]


class PresetTSMasterLinConfig(_PresetStructure):
    _fields_ = [
        ("application_channel", ctypes.c_uint8),
        ("hardware_channel", ctypes.c_uint8),
        ("protocol", ctypes.c_uint8),
        ("reserved", ctypes.c_uint8),
        ("hardware_index", ctypes.c_int32),
        ("device_type", ctypes.c_int32),
        ("baudrate", ctypes.c_uint32),
    ]


class PresetVectorLinConfig(_PresetStructure):
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


class PresetLinFrame(_PresetStructure):
    _fields_ = [
        ("frame_id", ctypes.c_uint8),
        ("data_len", ctypes.c_uint8),
        ("reserved", ctypes.c_uint8 * 2),
        ("data", ctypes.c_uint8 * 8),
    ]

    @property
    def payload(self) -> bytes:
        return bytes(self.data[: self.data_len])


class PresetToomossElinsMessage(_PresetStructure):
    _fields_ = [
        ("data_len", ctypes.c_uint8),
        ("break_bits", ctypes.c_uint8),
        ("status", ctypes.c_uint8),
        ("flags", ctypes.c_uint8),
        ("sync", ctypes.c_uint8),
        ("timestamp_high", ctypes.c_uint8),
        ("msg_send_times", ctypes.c_uint16),
        ("timestamp", ctypes.c_uint32),
        ("cmd_code", ctypes.c_uint8),
        ("device_id", ctypes.c_uint8),
        ("register_address", ctypes.c_uint16),
        ("crc16", ctypes.c_uint16),
        ("data", ctypes.c_uint8 * 64),
        ("ack_value", ctypes.c_uint8 * 4),
    ]

    @property
    def payload(self) -> bytes:
        return bytes(self.data[: min(self.data_len, len(self.data))])


class PresetPCANConfig(_PresetStructure):
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


class PresetTSMasterConfig(_PresetStructure):
    _fields_ = [
        ("application_channel", ctypes.c_uint8),
        ("hardware_channel", ctypes.c_uint8),
        ("is_fd", ctypes.c_uint8),
        ("brs", ctypes.c_uint8),
        ("hardware_index", ctypes.c_int32),
        ("device_type", ctypes.c_int32),
        ("nominal_bitrate", ctypes.c_uint32),
        ("data_bitrate", ctypes.c_uint32),
        ("rx_buffer_size", ctypes.c_uint32),
        ("poll_batch_size", ctypes.c_uint32),
        ("min_tx_interval_us", ctypes.c_uint32),
        ("min_rx_poll_interval_us", ctypes.c_uint32),
    ]


class PresetVectorConfig(_PresetStructure):
    _fields_ = [
        ("application_channel", ctypes.c_uint8),
        ("is_fd", ctypes.c_uint8),
        ("brs", ctypes.c_uint8),
        ("reserved", ctypes.c_uint8),
        ("hardware_type", ctypes.c_int32),
        ("hardware_index", ctypes.c_uint32),
        ("hardware_channel", ctypes.c_uint32),
        ("nominal_bitrate", ctypes.c_uint32),
        ("data_bitrate", ctypes.c_uint32),
        ("nominal_sjw", ctypes.c_uint32),
        ("nominal_tseg1", ctypes.c_uint32),
        ("nominal_tseg2", ctypes.c_uint32),
        ("data_sjw", ctypes.c_uint32),
        ("data_tseg1", ctypes.c_uint32),
        ("data_tseg2", ctypes.c_uint32),
        ("driver_rx_queue_size", ctypes.c_uint32),
        ("rx_buffer_size", ctypes.c_uint32),
        ("poll_batch_size", ctypes.c_uint32),
        ("min_tx_interval_us", ctypes.c_uint32),
        ("min_rx_poll_interval_us", ctypes.c_uint32),
    ]


class PresetCanFrame(_PresetStructure):
    _fields_ = [
        ("id", ctypes.c_uint32),
        ("data_len", ctypes.c_uint8),
        ("is_fd", ctypes.c_uint8),
        ("reserved", ctypes.c_uint8 * 2),
        ("data", ctypes.c_uint8 * 64),
    ]

    @property
    def payload(self) -> bytes:
        return bytes(self.data[: self.data_len])


class PresetCanFrameEx(_PresetStructure):
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

    @property
    def payload(self) -> bytes:
        return bytes(self.data[: self.data_len])


class PresetRxStats(_PresetStructure):
    _fields_ = [
        ("received", ctypes.c_uint64),
        ("dropped", ctypes.c_uint64),
        ("queued", ctypes.c_uint32),
        ("capacity", ctypes.c_uint32),
    ]


class PresetBusLoad(_PresetStructure):
    _fields_ = [
        ("load", ctypes.c_double),
        ("window_ns", ctypes.c_uint64),
        ("nominal_bitrate", ctypes.c_uint32),
        ("data_bitrate", ctypes.c_uint32),
        ("frame_count", ctypes.c_uint64),
    ]


DeviceHandle = ctypes.c_void_p
CanClientHandle = ctypes.c_void_p
LinClientHandle = ctypes.c_void_p
_U8P = ctypes.POINTER(ctypes.c_uint8)
_U32P = ctypes.POINTER(ctypes.c_uint32)
_SIZEP = ctypes.POINTER(ctypes.c_size_t)


_STATUS_NAMES = {
    PRESET_OK: "PRESET_OK",
    PRESET_ERR_NULL_PTR: "PRESET_ERR_NULL_PTR",
    PRESET_ERR_BUFFER_TOO_SMALL: "PRESET_ERR_BUFFER_TOO_SMALL",
    PRESET_ERR_NOT_INIT: "PRESET_ERR_NOT_INIT",
    PRESET_ERR_TIMEOUT: "PRESET_ERR_TIMEOUT",
    PRESET_ERR_UDS_NEGATIVE: "PRESET_ERR_UDS_NEGATIVE",
    PRESET_ERR_TRANSPORT: "PRESET_ERR_TRANSPORT",
    PRESET_ERR_UNSUPPORTED: "PRESET_ERR_UNSUPPORTED",
    PRESET_ERR_INVALID_ARG: "PRESET_ERR_INVALID_ARG",
    PRESET_ERR_BUSY: "PRESET_ERR_BUSY",
    PRESET_ERR_ISOTP: "PRESET_ERR_ISOTP",
    PRESET_ERR_PANIC: "PRESET_ERR_PANIC",
}


class PresetError(RuntimeError):
    """A non-zero status returned by ``preset_rs.dll``."""

    def __init__(self, status: int, operation: str, detail: str = "") -> None:
        self.status = int(status)
        self.operation = operation
        self.detail = detail
        name = _STATUS_NAMES.get(self.status, "PRESET_ERR_UNKNOWN")
        message = f"{operation} failed: {name} ({self.status})"
        if detail:
            message += f": {detail}"
        super().__init__(message)


def _bind(function: ctypes._CFuncPtr, argtypes: Sequence[type], restype: type | None) -> None:
    function.argtypes = list(argtypes)
    function.restype = restype


def _byte_input(data: bytes | bytearray | memoryview | Iterable[int]):
    raw = bytes(data)
    if not raw:
        return raw, None
    array = (ctypes.c_uint8 * len(raw)).from_buffer_copy(raw)
    return raw, array


class PresetRSNativeSDK:
    """Loaded ``preset_rs.dll`` plus process-wide API operations."""

    def __init__(self, dll_path: str | os.PathLike[str], *, validate_abi: bool = True) -> None:
        if os.name != "nt":
            raise OSError("preset_rs.dll is supported only on Windows")
        path = Path(dll_path).expanduser().resolve()
        if not path.is_file():
            raise FileNotFoundError(path)
        if ctypes.sizeof(ctypes.c_void_p) != 8:
            raise OSError("preset_rs ABI v6 requires 64-bit Python")

        # Search the target DLL directory for vendor dependencies while keeping
        # the normal Windows safe-search directories enabled.
        self.path = path
        self.dll = ctypes.CDLL(str(path), winmode=0x00000100 | 0x00001000)
        self._log_lock = threading.RLock()
        self._bind_api()

        if validate_abi:
            abi = self.abi_version
            if abi != SUPPORTED_ABI_VERSION:
                raise OSError(
                    f"unsupported preset_rs ABI version {abi}; expected {SUPPORTED_ABI_VERSION}"
                )

    def _bind_api(self) -> None:
        d = self.dll
        status = ctypes.c_int32
        u8 = ctypes.c_uint8
        u32 = ctypes.c_uint32
        size = ctypes.c_size_t

        _bind(d.preset_version, [], ctypes.c_char_p)
        _bind(d.preset_abi_version, [], u32)
        _bind(d.preset_get_capabilities, [], ctypes.c_uint64)
        for name, structure in (
            ("preset_default_config", PresetConfig),
            ("preset_tp_frame_default_config", PresetTpFrameConfig),
            ("preset_auto_default_config", PresetAutoConfig),
            ("preset_pcan_default_fd_timing", PresetCanFdTiming),
            ("preset_toomoss_default_fd_timing", PresetCanFdTiming),
            ("preset_toomoss_default_config", PresetToomossConfig),
            ("preset_toomoss_lin_default_config", PresetToomossLinConfig),
            ("preset_toomoss_elins_default_config", PresetToomossElinsConfig),
            ("preset_pcan_default_config", PresetPCANConfig),
            ("preset_tsmaster_default_config", PresetTSMasterConfig),
            ("preset_vector_default_config", PresetVectorConfig),
            ("preset_pcan_lin_default_config", PresetPCANLinConfig),
            ("preset_tsmaster_lin_default_config", PresetTSMasterLinConfig),
            ("preset_vector_lin_default_config", PresetVectorLinConfig),
        ):
            _bind(getattr(d, name), [], structure)

        _bind(d.preset_log_init, [ctypes.c_char_p], status)
        _bind(d.preset_log_shutdown, [], None)
        _bind(d.preset_set_print_log, [u8], None)
        _bind(d.preset_log_dropped_count, [], ctypes.c_uint64)
        _bind(d.preset_set_log_filter, [u8, _U32P, size], status)

        _bind(d.preset_auto_open, [ctypes.POINTER(PresetAutoConfig), ctypes.POINTER(DeviceHandle)], status)
        _bind(d.preset_device_get_backend, [DeviceHandle, _U8P], status)
        _bind(d.preset_toomoss_open, [ctypes.POINTER(DeviceHandle)], status)
        _bind(d.preset_toomoss_can_init, [DeviceHandle, ctypes.POINTER(PresetToomossConfig)], status)
        _bind(
            d.preset_toomoss_can_init_with_timing,
            [DeviceHandle, ctypes.POINTER(PresetToomossConfig), ctypes.POINTER(PresetCanFdTiming)],
            status,
        )
        _bind(d.preset_toomoss_lin_init, [DeviceHandle, ctypes.POINTER(PresetToomossLinConfig)], status)
        _bind(d.preset_toomoss_lin_write, [DeviceHandle, u8, u8, _U8P, size], status)
        _bind(
            d.preset_toomoss_lin_read,
            [DeviceHandle, u8, u8, ctypes.POINTER(PresetToomossLinFrame)],
            status,
        )
        _bind(
            d.preset_toomoss_lin_try_read,
            [DeviceHandle, u8, ctypes.POINTER(PresetToomossLinFrame), _SIZEP],
            status,
        )
        _bind(d.preset_toomoss_lin_break, [DeviceHandle, u8], status)
        _bind(d.preset_toomoss_lin_set_power, [DeviceHandle, u8, u8], status)
        _bind(d.preset_toomoss_elins_init, [DeviceHandle, ctypes.POINTER(PresetToomossElinsConfig)], status)
        _bind(
            d.preset_toomoss_elins_read,
            [DeviceHandle, u8, ctypes.POINTER(PresetToomossElinsMessage), _SIZEP],
            status,
        )

        _bind(d.preset_pcan_open, [ctypes.POINTER(PresetPCANConfig), ctypes.POINTER(DeviceHandle)], status)
        _bind(
            d.preset_pcan_open_with_timing,
            [ctypes.POINTER(PresetPCANConfig), ctypes.POINTER(PresetCanFdTiming), ctypes.POINTER(DeviceHandle)],
            status,
        )
        _bind(d.preset_pcan_can_init, [DeviceHandle, ctypes.POINTER(PresetPCANConfig)], status)
        _bind(
            d.preset_pcan_can_init_with_timing,
            [DeviceHandle, ctypes.POINTER(PresetPCANConfig), ctypes.POINTER(PresetCanFdTiming)],
            status,
        )
        _bind(d.preset_tsmaster_open, [ctypes.POINTER(PresetTSMasterConfig), ctypes.POINTER(DeviceHandle)], status)
        _bind(d.preset_tsmaster_can_init, [DeviceHandle, ctypes.POINTER(PresetTSMasterConfig)], status)
        _bind(d.preset_vector_open, [ctypes.POINTER(PresetVectorConfig), ctypes.POINTER(DeviceHandle)], status)
        _bind(d.preset_vector_can_init, [DeviceHandle, ctypes.POINTER(PresetVectorConfig)], status)

        tp_config = ctypes.POINTER(PresetTpFrameConfig)
        _bind(d.preset_can_write, [DeviceHandle, u8, u32, u8, _U8P, size], status)
        _bind(
            d.preset_can_try_read,
            [DeviceHandle, u8, ctypes.POINTER(PresetCanFrame), _SIZEP],
            status,
        )
        _bind(
            d.preset_can_try_read_ex,
            [DeviceHandle, u8, ctypes.POINTER(PresetCanFrameEx), _SIZEP],
            status,
        )
        _bind(
            d.preset_can_tp_encode_single_frame,
            [tp_config, _U8P, size, _U8P, size, _SIZEP],
            status,
        )
        _bind(
            d.preset_can_tp_encode_first_frame,
            [tp_config, _U8P, size, u32, _U8P, size, _SIZEP],
            status,
        )
        _bind(
            d.preset_can_tp_encode_flow_control_frame,
            [tp_config, u8, u8, u8, _U8P, size, _SIZEP],
            status,
        )
        _bind(
            d.preset_can_tp_encode_consecutive_frame,
            [tp_config, _U8P, size, u8, _U8P, size, _SIZEP],
            status,
        )
        _bind(
            d.preset_can_tp_build_frames,
            [_U8P, size, u8, u8, ctypes.POINTER(PresetTpEncodedFrame), _SIZEP],
            status,
        )
        _bind(
            d.preset_can_tp_write_single_frame,
            [DeviceHandle, u8, u32, tp_config, _U8P, size],
            status,
        )
        _bind(
            d.preset_can_tp_write_first_frame,
            [DeviceHandle, u8, u32, tp_config, _U8P, size, u32],
            status,
        )
        _bind(
            d.preset_can_tp_write_flow_control_frame,
            [DeviceHandle, u8, u32, tp_config, u8, u8, u8],
            status,
        )
        _bind(
            d.preset_can_tp_write_consecutive_frame,
            [DeviceHandle, u8, u32, tp_config, _U8P, size, u8],
            status,
        )

        _bind(d.preset_pcan_lin_open, [ctypes.POINTER(PresetPCANLinConfig), ctypes.POINTER(DeviceHandle)], status)
        _bind(d.preset_pcan_lin_init, [DeviceHandle, ctypes.POINTER(PresetPCANLinConfig)], status)
        _bind(
            d.preset_tsmaster_lin_open,
            [ctypes.POINTER(PresetTSMasterLinConfig), ctypes.POINTER(DeviceHandle)],
            status,
        )
        _bind(d.preset_tsmaster_lin_init, [DeviceHandle, ctypes.POINTER(PresetTSMasterLinConfig)], status)
        _bind(
            d.preset_vector_lin_open,
            [ctypes.POINTER(PresetVectorLinConfig), ctypes.POINTER(DeviceHandle)],
            status,
        )
        _bind(d.preset_vector_lin_init, [DeviceHandle, ctypes.POINTER(PresetVectorLinConfig)], status)
        _bind(d.preset_lin_master_write, [DeviceHandle, u8, u8, _U8P, size], status)
        _bind(
            d.preset_lin_master_read,
            [DeviceHandle, u8, u8, u32, ctypes.POINTER(PresetLinFrame)],
            status,
        )

        _bind(
            d.preset_can_uds_client_new,
            [DeviceHandle, u8, ctypes.POINTER(PresetConfig), ctypes.POINTER(CanClientHandle)],
            status,
        )
        request_args = [CanClientHandle, _U8P, size, u32, _U8P, size, _SIZEP]
        _bind(d.preset_can_uds_request, request_args, status)
        _bind(d.preset_can_uds_functional_request, request_args, status)
        _bind(d.preset_can_uds_set_default_st_min, [CanClientHandle, u32], status)
        _bind(d.preset_can_uds_set_default_block_size, [CanClientHandle, u32], status)
        _bind(d.preset_can_uds_set_manual_flow_control, [CanClientHandle, u8], status)
        _bind(d.preset_can_uds_set_manual_tp_mode, [CanClientHandle, u8], status)
        _bind(d.preset_can_uds_set_brs, [CanClientHandle, u8], status)
        _bind(d.preset_can_uds_get_brs, [CanClientHandle, _U8P], status)
        _bind(d.preset_can_uds_set_raw_tx_echo, [CanClientHandle, u8], status)
        _bind(d.preset_can_uds_write, [CanClientHandle, u32, u8, _U8P, size], status)
        _bind(
            d.preset_can_uds_tp_write_single_frame,
            [CanClientHandle, u32, tp_config, _U8P, size],
            status,
        )
        _bind(
            d.preset_can_uds_tp_write_first_frame,
            [CanClientHandle, u32, tp_config, _U8P, size, u32],
            status,
        )
        _bind(
            d.preset_can_uds_tp_write_flow_control_frame,
            [CanClientHandle, u32, tp_config, u8, u8, u8],
            status,
        )
        _bind(
            d.preset_can_uds_tp_write_consecutive_frame,
            [CanClientHandle, u32, tp_config, _U8P, size, u8],
            status,
        )
        _bind(
            d.preset_can_uds_try_read,
            [CanClientHandle, ctypes.POINTER(PresetCanFrame), _SIZEP],
            status,
        )
        _bind(
            d.preset_can_uds_try_read_ex,
            [CanClientHandle, ctypes.POINTER(PresetCanFrameEx), _SIZEP],
            status,
        )
        _bind(d.preset_can_uds_rx_get_stats, [CanClientHandle, ctypes.POINTER(PresetRxStats)], status)
        _bind(d.preset_can_uds_set_bus_load_enabled, [CanClientHandle, u8], status)
        _bind(d.preset_can_uds_get_bus_load, [CanClientHandle, ctypes.POINTER(PresetBusLoad)], status)
        _bind(d.preset_can_uds_last_error, [CanClientHandle, _U8P, size, _SIZEP], status)
        _bind(d.preset_can_uds_client_close, [ctypes.POINTER(CanClientHandle)], status)

        _bind(
            d.preset_lin_uds_client_new,
            [DeviceHandle, u8, u8, ctypes.POINTER(LinClientHandle)],
            status,
        )
        _bind(
            d.preset_lin_uds_request,
            [LinClientHandle, _U8P, size, u32, _U8P, _U8P, size, _SIZEP],
            status,
        )
        _bind(d.preset_lin_uds_set_poll_interval, [LinClientHandle, u32], status)
        _bind(d.preset_lin_uds_client_close, [ctypes.POINTER(LinClientHandle)], status)

        _bind(d.preset_device_close, [ctypes.POINTER(DeviceHandle)], status)
        _bind(d.preset_last_error, [_U8P, size, _SIZEP], status)

    @property
    def version(self) -> str:
        value = self.dll.preset_version()
        if value is None:
            raise RuntimeError("preset_version returned NULL")
        return value.decode("utf-8", errors="replace")

    @property
    def abi_version(self) -> int:
        return int(self.dll.preset_abi_version())

    @property
    def capabilities(self) -> int:
        return int(self.dll.preset_get_capabilities())

    def has_capability(self, capability: int) -> bool:
        return bool(self.capabilities & capability)

    def default_config(self) -> PresetConfig:
        return self.dll.preset_default_config()

    def tp_frame_default_config(self) -> PresetTpFrameConfig:
        return self.dll.preset_tp_frame_default_config()

    def _encode_tp_data_frame(
        self,
        name: str,
        config: PresetTpFrameConfig | None,
        data,
        *middle: int,
    ) -> bytes:
        raw, array = _byte_input(data)
        output = (ctypes.c_uint8 * 64)()
        count = ctypes.c_size_t()
        config_pointer = ctypes.byref(config) if config is not None else None
        status = getattr(self.dll, name)(
            config_pointer, array, len(raw), *middle,
            output, len(output), ctypes.byref(count),
        )
        if status == PRESET_ERR_BUFFER_TOO_SMALL:
            raise PresetError(status, name, f"frame requires {count.value} bytes")
        self._check(status, name)
        if count.value > len(output):
            raise RuntimeError(f"{name} returned invalid length {count.value}")
        return bytes(output[: count.value])

    def tp_encode_single_frame(
        self, data, config: PresetTpFrameConfig | None = None
    ) -> bytes:
        return self._encode_tp_data_frame(
            "preset_can_tp_encode_single_frame", config, data
        )

    def tp_encode_first_frame(
        self,
        first_chunk,
        total_message_size: int,
        config: PresetTpFrameConfig | None = None,
    ) -> bytes:
        return self._encode_tp_data_frame(
            "preset_can_tp_encode_first_frame",
            config,
            first_chunk,
            total_message_size,
        )

    def tp_encode_flow_control_frame(
        self,
        flow_status: int,
        block_size: int,
        st_min: int,
        config: PresetTpFrameConfig | None = None,
    ) -> bytes:
        output = (ctypes.c_uint8 * 64)()
        count = ctypes.c_size_t()
        config_pointer = ctypes.byref(config) if config is not None else None
        name = "preset_can_tp_encode_flow_control_frame"
        status = getattr(self.dll, name)(
            config_pointer, flow_status, block_size, st_min,
            output, len(output), ctypes.byref(count),
        )
        if status == PRESET_ERR_BUFFER_TOO_SMALL:
            raise PresetError(status, name, f"frame requires {count.value} bytes")
        self._check(status, name)
        if count.value > len(output):
            raise RuntimeError(f"{name} returned invalid length {count.value}")
        return bytes(output[: count.value])

    def tp_encode_consecutive_frame(
        self,
        data_chunk,
        sequence_number: int,
        config: PresetTpFrameConfig | None = None,
    ) -> bytes:
        return self._encode_tp_data_frame(
            "preset_can_tp_encode_consecutive_frame",
            config,
            data_chunk,
            sequence_number,
        )

    def tp_build_frames(
        self, payload, *, is_fd: bool = False, padding_byte: int = 0xAA
    ) -> list[PresetTpEncodedFrame]:
        raw, array = _byte_input(payload)
        count = ctypes.c_size_t()
        name = "preset_can_tp_build_frames"
        status = getattr(self.dll, name)(
            array, len(raw), bool(is_fd), padding_byte, None, ctypes.byref(count)
        )
        if status != PRESET_ERR_BUFFER_TOO_SMALL:
            self._check(status, name)
        if count.value == 0:
            return []
        frames = (PresetTpEncodedFrame * count.value)()
        capacity = count.value
        status = getattr(self.dll, name)(
            array, len(raw), bool(is_fd), padding_byte, frames, ctypes.byref(count)
        )
        self._check(status, name)
        if count.value > capacity:
            raise RuntimeError(
                f"{name} returned invalid frame count {count.value} for capacity {capacity}"
            )
        return list(frames[: count.value])

    def auto_default_config(self) -> PresetAutoConfig:
        return self.dll.preset_auto_default_config()

    def pcan_default_fd_timing(self) -> PresetCanFdTiming:
        return self.dll.preset_pcan_default_fd_timing()

    def toomoss_default_fd_timing(self) -> PresetCanFdTiming:
        return self.dll.preset_toomoss_default_fd_timing()

    def toomoss_default_config(self) -> PresetToomossConfig:
        return self.dll.preset_toomoss_default_config()

    def toomoss_lin_default_config(self) -> PresetToomossLinConfig:
        return self.dll.preset_toomoss_lin_default_config()

    def toomoss_elins_default_config(self) -> PresetToomossElinsConfig:
        return self.dll.preset_toomoss_elins_default_config()

    def pcan_default_config(self) -> PresetPCANConfig:
        return self.dll.preset_pcan_default_config()

    def tsmaster_default_config(self) -> PresetTSMasterConfig:
        return self.dll.preset_tsmaster_default_config()

    def vector_default_config(self) -> PresetVectorConfig:
        return self.dll.preset_vector_default_config()

    def pcan_lin_default_config(self) -> PresetPCANLinConfig:
        return self.dll.preset_pcan_lin_default_config()

    def tsmaster_lin_default_config(self) -> PresetTSMasterLinConfig:
        return self.dll.preset_tsmaster_lin_default_config()

    def vector_lin_default_config(self) -> PresetVectorLinConfig:
        return self.dll.preset_vector_lin_default_config()

    def last_error(self) -> str:
        return self._read_error(self.dll.preset_last_error)

    def _read_error(self, function, *prefix) -> str:
        required = ctypes.c_size_t()
        status = int(function(*prefix, None, 0, ctypes.byref(required)))
        if status == PRESET_OK or required.value == 0:
            return ""
        if status != PRESET_ERR_BUFFER_TOO_SMALL:
            return ""
        output = (ctypes.c_uint8 * required.value)()
        status = int(function(*prefix, output, len(output), ctypes.byref(required)))
        if status != PRESET_OK or required.value > len(output):
            return ""
        return bytes(output[: required.value]).decode("utf-8", errors="replace")

    def _check(self, status: int, operation: str) -> None:
        if status != PRESET_OK:
            raise PresetError(status, operation, self.last_error())

    def log_init(self, name: str) -> None:
        if "\x00" in name:
            raise ValueError("log name cannot contain NUL")
        with self._log_lock:
            self._check(self.dll.preset_log_init(name.encode("utf-8")), "preset_log_init")

    def log_shutdown(self) -> None:
        with self._log_lock:
            self.dll.preset_log_shutdown()

    def set_print_log(self, enabled: bool) -> None:
        with self._log_lock:
            self.dll.preset_set_print_log(bool(enabled))

    @property
    def log_dropped_count(self) -> int:
        return int(self.dll.preset_log_dropped_count())

    def set_log_filter(self, mode: int, ids: Iterable[int] = ()) -> None:
        values = tuple(ids)
        array = (ctypes.c_uint32 * len(values))(*values) if values else None
        self._check(
            self.dll.preset_set_log_filter(mode, array, len(values)),
            "preset_set_log_filter",
        )

    def _open(self, function_name: str, *configs: ctypes.Structure) -> "PresetDevice":
        handle = DeviceHandle()
        arguments = [ctypes.byref(config) for config in configs]
        arguments.append(ctypes.byref(handle))
        status = int(getattr(self.dll, function_name)(*arguments))
        self._check(status, function_name)
        if not handle.value:
            raise RuntimeError(f"{function_name} succeeded but returned NULL")
        return PresetDevice(self, handle)

    def auto_open(self, config: PresetAutoConfig | None = None) -> "PresetDevice":
        if config is None:
            config = self.auto_default_config()
        return self._open("preset_auto_open", config)

    def toomoss_open(self) -> "PresetDevice":
        return self._open("preset_toomoss_open")

    def pcan_open(
        self,
        config: PresetPCANConfig | None = None,
        timing: PresetCanFdTiming | None = None,
    ) -> "PresetDevice":
        if config is None:
            config = self.pcan_default_config()
        if timing is None:
            return self._open("preset_pcan_open", config)
        return self._open("preset_pcan_open_with_timing", config, timing)

    def tsmaster_open(self, config: PresetTSMasterConfig | None = None) -> "PresetDevice":
        if config is None:
            config = self.tsmaster_default_config()
        return self._open("preset_tsmaster_open", config)

    def vector_open(self, config: PresetVectorConfig | None = None) -> "PresetDevice":
        if config is None:
            config = self.vector_default_config()
        return self._open("preset_vector_open", config)

    def pcan_lin_open(self, config: PresetPCANLinConfig | None = None) -> "PresetDevice":
        if config is None:
            config = self.pcan_lin_default_config()
        return self._open("preset_pcan_lin_open", config)

    def tsmaster_lin_open(self, config: PresetTSMasterLinConfig | None = None) -> "PresetDevice":
        if config is None:
            config = self.tsmaster_lin_default_config()
        return self._open("preset_tsmaster_lin_open", config)

    def vector_lin_open(self, config: PresetVectorLinConfig | None = None) -> "PresetDevice":
        if config is None:
            config = self.vector_lin_default_config()
        return self._open("preset_vector_lin_open", config)


class _NativeHandle:
    _close_symbol = ""

    def __init__(self, sdk: PresetRSNativeSDK, handle: ctypes.c_void_p) -> None:
        self._sdk = sdk
        self._handle = handle
        self._lock = threading.RLock()

    @property
    def closed(self) -> bool:
        with self._lock:
            return not bool(self._handle.value)

    def _require_handle(self) -> ctypes.c_void_p:
        if not self._handle.value:
            raise RuntimeError(f"{type(self).__name__} is closed")
        return self._handle

    def close(self) -> None:
        with self._lock:
            if not self._handle.value:
                return
            status = int(getattr(self._sdk.dll, self._close_symbol)(ctypes.byref(self._handle)))
            self._sdk._check(status, self._close_symbol)

    def __enter__(self):
        self._require_handle()
        return self

    def __exit__(self, exc_type, exc_value, traceback) -> None:
        self.close()

    def __del__(self) -> None:
        try:
            self.close()
        except Exception:
            # Destructors cannot report native shutdown errors.  Explicit close
            # or a context manager should be used when the error matters.
            pass


class PresetDevice(_NativeHandle):
    """Owning native device handle."""

    _close_symbol = "preset_device_close"

    @property
    def backend(self) -> int:
        with self._lock:
            handle = self._require_handle()
            value = ctypes.c_uint8()
            self._sdk._check(
                self._sdk.dll.preset_device_get_backend(handle, ctypes.byref(value)),
                "preset_device_get_backend",
            )
            return int(value.value)

    def _init(self, name: str, *configs: ctypes.Structure) -> None:
        with self._lock:
            arguments = [self._require_handle(), *(ctypes.byref(value) for value in configs)]
            self._sdk._check(getattr(self._sdk.dll, name)(*arguments), name)

    def toomoss_can_init(
        self,
        config: PresetToomossConfig,
        timing: PresetCanFdTiming | None = None,
    ) -> None:
        if timing is None:
            self._init("preset_toomoss_can_init", config)
        else:
            self._init("preset_toomoss_can_init_with_timing", config, timing)

    def pcan_can_init(
        self,
        config: PresetPCANConfig,
        timing: PresetCanFdTiming | None = None,
    ) -> None:
        if timing is None:
            self._init("preset_pcan_can_init", config)
        else:
            self._init("preset_pcan_can_init_with_timing", config, timing)

    def tsmaster_can_init(self, config: PresetTSMasterConfig) -> None:
        self._init("preset_tsmaster_can_init", config)

    def vector_can_init(self, config: PresetVectorConfig) -> None:
        self._init("preset_vector_can_init", config)

    def toomoss_lin_init(self, config: PresetToomossLinConfig) -> None:
        self._init("preset_toomoss_lin_init", config)

    def toomoss_elins_init(self, config: PresetToomossElinsConfig) -> None:
        self._init("preset_toomoss_elins_init", config)

    def pcan_lin_init(self, config: PresetPCANLinConfig) -> None:
        self._init("preset_pcan_lin_init", config)

    def tsmaster_lin_init(self, config: PresetTSMasterLinConfig) -> None:
        self._init("preset_tsmaster_lin_init", config)

    def vector_lin_init(self, config: PresetVectorLinConfig) -> None:
        self._init("preset_vector_lin_init", config)

    def new_can_uds_client(self, channel: int, config: PresetConfig) -> "PresetCanUdsClient":
        with self._lock:
            handle = CanClientHandle()
            status = self._sdk.dll.preset_can_uds_client_new(
                self._require_handle(), channel, ctypes.byref(config), ctypes.byref(handle)
            )
            self._sdk._check(status, "preset_can_uds_client_new")
            if not handle.value:
                raise RuntimeError("preset_can_uds_client_new succeeded but returned NULL")
            return PresetCanUdsClient(self._sdk, handle, self, max(1, int(config.max_pdu_len)))

    def new_lin_uds_client(self, channel: int, nad: int) -> "PresetLinUdsClient":
        with self._lock:
            handle = LinClientHandle()
            status = self._sdk.dll.preset_lin_uds_client_new(
                self._require_handle(), channel, nad, ctypes.byref(handle)
            )
            self._sdk._check(status, "preset_lin_uds_client_new")
            if not handle.value:
                raise RuntimeError("preset_lin_uds_client_new succeeded but returned NULL")
            return PresetLinUdsClient(self._sdk, handle, self)

    def can_write(self, channel: int, can_id: int, data, *, is_fd: bool = False) -> None:
        raw, array = _byte_input(data)
        name = "preset_can_write"
        with self._lock:
            self._sdk._check(
                getattr(self._sdk.dll, name)(
                    self._require_handle(), channel, can_id, bool(is_fd), array, len(raw)
                ),
                name,
            )

    def can_try_read(self, channel: int, capacity: int = 256) -> list[PresetCanFrame]:
        return self._read_many("preset_can_try_read", PresetCanFrame, channel, capacity)

    def can_try_read_ex(self, channel: int, capacity: int = 256) -> list[PresetCanFrameEx]:
        return self._read_many("preset_can_try_read_ex", PresetCanFrameEx, channel, capacity)

    def _tp_write_data(
        self,
        name: str,
        channel: int,
        can_id: int,
        config: PresetTpFrameConfig | None,
        data,
        *tail: int,
    ) -> None:
        raw, array = _byte_input(data)
        config_pointer = ctypes.byref(config) if config is not None else None
        with self._lock:
            self._sdk._check(
                getattr(self._sdk.dll, name)(
                    self._require_handle(), channel, can_id, config_pointer,
                    array, len(raw), *tail,
                ),
                name,
            )

    def tp_write_single_frame(
        self,
        channel: int,
        can_id: int,
        data,
        config: PresetTpFrameConfig | None = None,
    ) -> None:
        self._tp_write_data(
            "preset_can_tp_write_single_frame", channel, can_id, config, data
        )

    def tp_write_first_frame(
        self,
        channel: int,
        can_id: int,
        first_chunk,
        total_message_size: int,
        config: PresetTpFrameConfig | None = None,
    ) -> None:
        self._tp_write_data(
            "preset_can_tp_write_first_frame",
            channel,
            can_id,
            config,
            first_chunk,
            total_message_size,
        )

    def tp_write_flow_control_frame(
        self,
        channel: int,
        can_id: int,
        flow_status: int,
        block_size: int,
        st_min: int,
        config: PresetTpFrameConfig | None = None,
    ) -> None:
        name = "preset_can_tp_write_flow_control_frame"
        config_pointer = ctypes.byref(config) if config is not None else None
        with self._lock:
            self._sdk._check(
                getattr(self._sdk.dll, name)(
                    self._require_handle(), channel, can_id, config_pointer,
                    flow_status, block_size, st_min,
                ),
                name,
            )

    def tp_write_consecutive_frame(
        self,
        channel: int,
        can_id: int,
        data_chunk,
        sequence_number: int,
        config: PresetTpFrameConfig | None = None,
    ) -> None:
        self._tp_write_data(
            "preset_can_tp_write_consecutive_frame",
            channel,
            can_id,
            config,
            data_chunk,
            sequence_number,
        )

    def _write(self, name: str, channel: int, frame_id: int, data) -> None:
        raw, array = _byte_input(data)
        with self._lock:
            self._sdk._check(
                getattr(self._sdk.dll, name)(self._require_handle(), channel, frame_id, array, len(raw)),
                name,
            )

    def toomoss_lin_write(self, channel: int, frame_id: int, data) -> None:
        self._write("preset_toomoss_lin_write", channel, frame_id, data)

    def lin_master_write(self, channel: int, frame_id: int, data) -> None:
        self._write("preset_lin_master_write", channel, frame_id, data)

    def toomoss_lin_read(self, channel: int, frame_id: int) -> PresetToomossLinFrame:
        with self._lock:
            frame = PresetToomossLinFrame()
            self._sdk._check(
                self._sdk.dll.preset_toomoss_lin_read(
                    self._require_handle(), channel, frame_id, ctypes.byref(frame)
                ),
                "preset_toomoss_lin_read",
            )
            return frame

    def lin_master_read(self, channel: int, frame_id: int, timeout_ms: int) -> PresetLinFrame:
        with self._lock:
            frame = PresetLinFrame()
            self._sdk._check(
                self._sdk.dll.preset_lin_master_read(
                    self._require_handle(), channel, frame_id, timeout_ms, ctypes.byref(frame)
                ),
                "preset_lin_master_read",
            )
            return frame

    def _read_many(self, name: str, structure: type[_PresetStructure], channel: int, capacity: int):
        if capacity <= 0:
            raise ValueError("capacity must be positive")
        values = (structure * capacity)()
        count = ctypes.c_size_t(capacity)
        with self._lock:
            self._sdk._check(
                getattr(self._sdk.dll, name)(
                    self._require_handle(), channel, values, ctypes.byref(count)
                ),
                name,
            )
        if count.value > capacity:
            raise RuntimeError(
                f"{name} returned invalid length {count.value} for capacity {capacity}"
            )
        return list(values[: count.value])

    def toomoss_lin_try_read(self, channel: int, capacity: int = 256) -> list[PresetToomossLinFrame]:
        return self._read_many("preset_toomoss_lin_try_read", PresetToomossLinFrame, channel, capacity)

    def toomoss_elins_read(self, channel: int, capacity: int = 256) -> list[PresetToomossElinsMessage]:
        return self._read_many("preset_toomoss_elins_read", PresetToomossElinsMessage, channel, capacity)

    def toomoss_lin_break(self, channel: int) -> None:
        with self._lock:
            self._sdk._check(
                self._sdk.dll.preset_toomoss_lin_break(self._require_handle(), channel),
                "preset_toomoss_lin_break",
            )

    def toomoss_lin_set_power(self, channel: int, voltage: int) -> None:
        with self._lock:
            self._sdk._check(
                self._sdk.dll.preset_toomoss_lin_set_power(
                    self._require_handle(), channel, voltage
                ),
                "preset_toomoss_lin_set_power",
            )


class PresetCanUdsClient(_NativeHandle):
    """CAN/CAN-FD UDS client borrowing a :class:`PresetDevice`."""

    _close_symbol = "preset_can_uds_client_close"

    def __init__(
        self,
        sdk: PresetRSNativeSDK,
        handle: ctypes.c_void_p,
        device: PresetDevice,
        response_capacity: int,
    ) -> None:
        super().__init__(sdk, handle)
        self._device = device
        self.response_capacity = response_capacity

    def _request(self, name: str, payload, timeout_ms: int, response_capacity: int | None) -> bytes:
        raw, request = _byte_input(payload)
        capacity = self.response_capacity if response_capacity is None else response_capacity
        if capacity <= 0:
            raise ValueError("response_capacity must be positive")
        output = (ctypes.c_uint8 * capacity)()
        count = ctypes.c_size_t()
        with self._lock:
            status = getattr(self._sdk.dll, name)(
                self._require_handle(), request, len(raw), timeout_ms,
                output, capacity, ctypes.byref(count),
            )
            if status == PRESET_ERR_BUFFER_TOO_SMALL:
                raise PresetError(status, name, f"response requires {count.value} bytes")
            self._sdk._check(status, name)
        if count.value > capacity:
            raise RuntimeError(f"{name} returned invalid length {count.value} for capacity {capacity}")
        return bytes(output[: count.value])

    def request(self, payload, timeout_ms: int, response_capacity: int | None = None) -> bytes:
        return self._request("preset_can_uds_request", payload, timeout_ms, response_capacity)

    def functional_request(self, payload, timeout_ms: int, response_capacity: int | None = None) -> bytes:
        return self._request(
            "preset_can_uds_functional_request", payload, timeout_ms, response_capacity
        )

    def set_default_st_min(self, st_min_ms: int) -> None:
        self._simple("preset_can_uds_set_default_st_min", st_min_ms)

    def set_default_block_size(self, block_size: int) -> None:
        self._simple("preset_can_uds_set_default_block_size", block_size)

    def set_manual_flow_control(self, enabled: bool) -> None:
        self._simple("preset_can_uds_set_manual_flow_control", bool(enabled))

    def set_manual_tp_mode(self, enabled: bool) -> None:
        self._simple("preset_can_uds_set_manual_tp_mode", bool(enabled))

    def set_brs(self, enabled: bool) -> None:
        self._simple("preset_can_uds_set_brs", bool(enabled))

    def set_raw_tx_echo(self, enabled: bool) -> None:
        self._simple("preset_can_uds_set_raw_tx_echo", bool(enabled))

    def set_bus_load_enabled(self, enabled: bool) -> None:
        self._simple("preset_can_uds_set_bus_load_enabled", bool(enabled))

    def _simple(self, name: str, value: int | bool) -> None:
        with self._lock:
            self._sdk._check(getattr(self._sdk.dll, name)(self._require_handle(), value), name)

    @property
    def brs(self) -> bool:
        with self._lock:
            value = ctypes.c_uint8()
            self._sdk._check(
                self._sdk.dll.preset_can_uds_get_brs(self._require_handle(), ctypes.byref(value)),
                "preset_can_uds_get_brs",
            )
            return bool(value.value)

    def write(self, can_id: int, data, *, is_fd: bool = False) -> None:
        raw, array = _byte_input(data)
        with self._lock:
            self._sdk._check(
                self._sdk.dll.preset_can_uds_write(
                    self._require_handle(), can_id, bool(is_fd), array, len(raw)
                ),
                "preset_can_uds_write",
            )

    def _tp_write_data(
        self,
        name: str,
        can_id: int,
        config: PresetTpFrameConfig | None,
        data,
        *tail: int,
    ) -> None:
        raw, array = _byte_input(data)
        config_pointer = ctypes.byref(config) if config is not None else None
        with self._lock:
            self._sdk._check(
                getattr(self._sdk.dll, name)(
                    self._require_handle(), can_id, config_pointer,
                    array, len(raw), *tail,
                ),
                name,
            )

    def tp_write_single_frame(
        self,
        can_id: int,
        data,
        config: PresetTpFrameConfig | None = None,
    ) -> None:
        self._tp_write_data(
            "preset_can_uds_tp_write_single_frame", can_id, config, data
        )

    def tp_write_first_frame(
        self,
        can_id: int,
        first_chunk,
        total_message_size: int,
        config: PresetTpFrameConfig | None = None,
    ) -> None:
        self._tp_write_data(
            "preset_can_uds_tp_write_first_frame",
            can_id,
            config,
            first_chunk,
            total_message_size,
        )

    def tp_write_flow_control_frame(
        self,
        can_id: int,
        flow_status: int,
        block_size: int,
        st_min: int,
        config: PresetTpFrameConfig | None = None,
    ) -> None:
        name = "preset_can_uds_tp_write_flow_control_frame"
        config_pointer = ctypes.byref(config) if config is not None else None
        with self._lock:
            self._sdk._check(
                getattr(self._sdk.dll, name)(
                    self._require_handle(), can_id, config_pointer,
                    flow_status, block_size, st_min,
                ),
                name,
            )

    def tp_write_consecutive_frame(
        self,
        can_id: int,
        data_chunk,
        sequence_number: int,
        config: PresetTpFrameConfig | None = None,
    ) -> None:
        self._tp_write_data(
            "preset_can_uds_tp_write_consecutive_frame",
            can_id,
            config,
            data_chunk,
            sequence_number,
        )

    def _try_read(self, name: str, structure: type[_PresetStructure], capacity: int):
        if capacity <= 0:
            raise ValueError("capacity must be positive")
        frames = (structure * capacity)()
        count = ctypes.c_size_t(capacity)
        with self._lock:
            self._sdk._check(
                getattr(self._sdk.dll, name)(
                    self._require_handle(), frames, ctypes.byref(count)
                ),
                name,
            )
        if count.value > capacity:
            raise RuntimeError(
                f"{name} returned invalid length {count.value} for capacity {capacity}"
            )
        return list(frames[: count.value])

    def try_read(self, capacity: int = 256) -> list[PresetCanFrame]:
        return self._try_read("preset_can_uds_try_read", PresetCanFrame, capacity)

    def try_read_ex(self, capacity: int = 256) -> list[PresetCanFrameEx]:
        return self._try_read("preset_can_uds_try_read_ex", PresetCanFrameEx, capacity)

    @property
    def rx_stats(self) -> PresetRxStats:
        with self._lock:
            stats = PresetRxStats()
            self._sdk._check(
                self._sdk.dll.preset_can_uds_rx_get_stats(
                    self._require_handle(), ctypes.byref(stats)
                ),
                "preset_can_uds_rx_get_stats",
            )
            return stats

    @property
    def bus_load(self) -> PresetBusLoad:
        with self._lock:
            load = PresetBusLoad()
            self._sdk._check(
                self._sdk.dll.preset_can_uds_get_bus_load(
                    self._require_handle(), ctypes.byref(load)
                ),
                "preset_can_uds_get_bus_load",
            )
            return load

    def last_error(self) -> str:
        with self._lock:
            return self._sdk._read_error(
                self._sdk.dll.preset_can_uds_last_error, self._require_handle()
            )


class PresetLinUdsClient(_NativeHandle):
    """LIN transport-protocol/UDS client borrowing a :class:`PresetDevice`."""

    _close_symbol = "preset_lin_uds_client_close"

    def __init__(
        self, sdk: PresetRSNativeSDK, handle: ctypes.c_void_p, device: PresetDevice
    ) -> None:
        super().__init__(sdk, handle)
        self._device = device

    def request(
        self,
        payload,
        timeout_ms: int,
        response_capacity: int = 4095,
    ) -> tuple[int, bytes]:
        if response_capacity <= 0:
            raise ValueError("response_capacity must be positive")
        raw, request = _byte_input(payload)
        response_nad = ctypes.c_uint8()
        output = (ctypes.c_uint8 * response_capacity)()
        count = ctypes.c_size_t()
        with self._lock:
            status = self._sdk.dll.preset_lin_uds_request(
                self._require_handle(), request, len(raw), timeout_ms,
                ctypes.byref(response_nad), output, response_capacity, ctypes.byref(count),
            )
            if status == PRESET_ERR_BUFFER_TOO_SMALL:
                raise PresetError(
                    status,
                    "preset_lin_uds_request",
                    f"response requires {count.value} bytes",
                )
            self._sdk._check(status, "preset_lin_uds_request")
        if count.value > response_capacity:
            raise RuntimeError(
                "preset_lin_uds_request returned invalid length "
                f"{count.value} for capacity {response_capacity}"
            )
        return int(response_nad.value), bytes(output[: count.value])

    def set_poll_interval(self, interval_ms: int) -> None:
        with self._lock:
            self._sdk._check(
                self._sdk.dll.preset_lin_uds_set_poll_interval(
                    self._require_handle(), interval_ms
                ),
                "preset_lin_uds_set_poll_interval",
            )


# Short aliases for application code that does not need the ABI-specific name.
PresetRS = PresetRSNativeSDK
Device = PresetDevice
CanUdsClient = PresetCanUdsClient
LinUdsClient = PresetLinUdsClient
