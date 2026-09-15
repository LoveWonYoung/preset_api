"""统一 CAN 设备二次封装。

调用方只构造 MyDevice 并 open：DLL、后端配置、通道初始化和 UDS client
都在类内部完成。
"""

from __future__ import annotations

import os
import time
from dataclasses import dataclass
from enum import IntEnum
from typing import Sequence

from preset_rs_native_sdk import (
    PRESET_CAN_BACKEND_PCAN,
    PRESET_CAN_BACKEND_TOOMOSS,
    PRESET_CAN_BACKEND_TSMASTER,
    PRESET_CAN_BACKEND_VECTOR,
    PRESET_TSMASTER_TC1016,
    PRESET_VECTOR_HWTYPE_ANY,
    PresetCanFrameEx,
    PresetRSNativeSDK,
)


_DEFAULT_DLL = os.getcwd() + r"\preset_rs.dll"


class Backend(IntEnum):
    """与 preset_rs 后端编号一致。"""

    Toomoss = PRESET_CAN_BACKEND_TOOMOSS
    TSMaster = PRESET_CAN_BACKEND_TSMASTER
    Pcan = PRESET_CAN_BACKEND_PCAN
    Vector = PRESET_CAN_BACKEND_VECTOR


@dataclass(frozen=True)
class CanFrame:
    id: int
    dlc: int
    data: bytes
    direction: str
    is_fd: bool
    brs: bool

    @classmethod
    def from_native(cls, frame: PresetCanFrameEx) -> "CanFrame":
        return cls(
            id=int(frame.id),
            dlc=int(frame.dlc),
            data=frame.payload,
            direction="TX" if frame.direction == 0 else "RX",
            is_fd=bool(frame.is_fd),
            brs=bool(frame.brs),
        )


class MyDevice:
    def __init__(
        self,
        backend: Backend | int,
        phys_id: int,
        resp_id: int,
        func_id: int,
        channels: Sequence[int],
        device_type: int | None = None,
        *,
        dll_path: str | None = None,
        is_fd: bool = True,
        brs: bool = True,
        nominal_bitrate: int = 500_000,
        data_bitrate: int = 2_000_000,
        hardware_index: int = 0,
        raw_rx_enabled: bool = True,
    ) -> None:
        if not channels:
            raise ValueError("channels 不能为空")
        self.backend = Backend(backend)
        self.phys_id = phys_id
        self.resp_id = resp_id
        self.func_id = func_id
        self.channels = [int(ch) for ch in channels]
        self.device_type = device_type
        self.dll_path = dll_path or _DEFAULT_DLL
        self.is_fd = is_fd
        self.brs = brs
        self.nominal_bitrate = nominal_bitrate
        self.data_bitrate = data_bitrate
        self.hardware_index = hardware_index
        self.raw_rx_enabled = raw_rx_enabled
        self._sdk: PresetRSNativeSDK | None = None
        self._device = None
        self._clients: dict[int, object] = {}

    def open(self) -> "MyDevice":
        if self._device is not None:
            return self
        self._sdk = PresetRSNativeSDK(self.dll_path)
        try:
            self._device = self._open_backend()
            self._init_extra_channels()
            for channel in self.channels:
                self._clients[channel] = self._device.new_can_uds_client(
                    channel, self._uds_config()
                )
        except Exception:
            self.close()
            raise
        return self

    def txfn(
        self,
        can_id: int,
        data,
        *,
        channel: int | None = None,
        is_fd: bool | None = None,
    ) -> None:
        client = self._client(channel)
        client.write(can_id, data, is_fd=self.is_fd if is_fd is None else is_fd)

    def rxfn(
        self,
        timeout_ms: int = 0,
        *,
        channel: int | None = None,
        capacity: int = 256,
    ) -> list[CanFrame]:
        client = self._client(channel)
        if timeout_ms <= 0:
            return [CanFrame.from_native(frame) for frame in client.try_read_ex(capacity)]

        frames: list[CanFrame] = []
        deadline = time.monotonic() + timeout_ms / 1000.0
        while True:
            frames.extend(CanFrame.from_native(frame) for frame in client.try_read_ex(capacity))
            if frames or time.monotonic() >= deadline:
                return frames
            time.sleep(0.02)

    def request(
        self,
        payload,
        timeout_ms: int = 1000,
        *,
        channel: int | None = None,
        functional: bool = False,
    ) -> bytes:
        client = self._client(channel)
        if functional:
            return client.functional_request(payload, timeout_ms=timeout_ms)
        return client.request(payload, timeout_ms=timeout_ms)

    def close(self) -> None:
        first_error: Exception | None = None
        clients = list(self._clients.values())
        self._clients.clear()
        for client in clients:
            try:
                client.close()
            except Exception as exc:
                if first_error is None:
                    first_error = exc
        device = self._device
        self._device = None
        if device is not None:
            try:
                device.close()
            except Exception as exc:
                if first_error is None:
                    first_error = exc
        if first_error is not None:
            raise first_error

    def __enter__(self) -> "MyDevice":
        return self.open()

    def __exit__(self, exc_type, exc_value, traceback) -> None:
        self.close()

    def _client(self, channel: int | None):
        if self._device is None:
            raise RuntimeError("设备未打开，先调用 open()")
        target = self.channels[0] if channel is None else int(channel)
        client = self._clients.get(target)
        if client is None:
            raise ValueError(f"通道 {target} 未初始化，可用通道: {self.channels}")
        return client

    def _uds_config(self):
        assert self._sdk is not None
        cfg = self._sdk.default_config()
        cfg.physical_id = self.phys_id
        cfg.response_id = self.resp_id
        cfg.functional_id = self.func_id
        cfg.is_fd = 1 if self.is_fd else 0
        cfg.raw_rx_enabled = 1 if self.raw_rx_enabled else 0
        return cfg

    def _open_backend(self):
        assert self._sdk is not None
        first = self.channels[0]
        if self.backend is Backend.Toomoss:
            device = self._sdk.toomoss_open()
            for channel in self.channels:
                device.toomoss_can_init(self._toomoss_config(channel))
            return device
        if self.backend is Backend.TSMaster:
            return self._sdk.tsmaster_open(self._tsmaster_config(first))
        if self.backend is Backend.Pcan:
            return self._sdk.pcan_open(self._pcan_config(first))
        if self.backend is Backend.Vector:
            return self._sdk.vector_open(self._vector_config(first))
        raise ValueError(f"不支持的后端: {self.backend}")

    def _init_extra_channels(self) -> None:
        extra = self.channels[1:]
        if not extra:
            return
        assert self._device is not None
        if self.backend is Backend.Toomoss:
            return
        if self.backend is Backend.TSMaster:
            for channel in extra:
                self._device.tsmaster_can_init(self._tsmaster_config(channel))
            return
        if self.backend is Backend.Pcan:
            for channel in extra:
                self._device.pcan_can_init(self._pcan_config(channel))
            return
        if self.backend is Backend.Vector:
            for channel in extra:
                self._device.vector_can_init(self._vector_config(channel))

    def _toomoss_config(self, channel: int):
        assert self._sdk is not None
        cfg = self._sdk.toomoss_default_config()
        cfg.channel = channel
        cfg.brs = 1 if self.brs else 0
        cfg.nominal_bitrate = self.nominal_bitrate
        cfg.data_bitrate = self.data_bitrate
        return cfg

    def _tsmaster_config(self, channel: int):
        assert self._sdk is not None
        cfg = self._sdk.tsmaster_default_config()
        cfg.application_channel = channel
        cfg.hardware_channel = channel
        cfg.hardware_index = self.hardware_index
        cfg.device_type = (
            PRESET_TSMASTER_TC1016 if self.device_type is None else self.device_type
        )
        cfg.is_fd = 1 if self.is_fd else 0
        cfg.brs = 1 if self.brs else 0
        cfg.nominal_bitrate = self.nominal_bitrate
        cfg.data_bitrate = self.data_bitrate
        return cfg

    def _pcan_config(self, channel: int):
        assert self._sdk is not None
        cfg = self._sdk.pcan_default_config()
        cfg.channel = channel
        cfg.is_fd = 1 if self.is_fd else 0
        cfg.brs = 1 if self.brs else 0
        cfg.nominal_bitrate = self.nominal_bitrate
        cfg.data_bitrate = self.data_bitrate
        return cfg

    def _vector_config(self, channel: int):
        assert self._sdk is not None
        cfg = self._sdk.vector_default_config()
        cfg.application_channel = channel
        cfg.hardware_channel = channel
        cfg.hardware_index = self.hardware_index
        cfg.hardware_type = (
            PRESET_VECTOR_HWTYPE_ANY if self.device_type is None else self.device_type
        )
        cfg.is_fd = 1 if self.is_fd else 0
        cfg.brs = 1 if self.brs else 0
        cfg.nominal_bitrate = self.nominal_bitrate
        cfg.data_bitrate = self.data_bitrate
        return cfg
