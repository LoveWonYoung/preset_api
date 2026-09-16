"""TSMaster CAN / UDS 用法示例。

TSMaster 没有 Toomoss 那套 PresetCanFdTiming。
`tsmaster_open(cfg)` 已经初始化第一路通道；同一路不要再 `tsmaster_can_init`。
UDS client 的 channel 用 `application_channel`，不要写 `cfg.channel`（该字段不存在）。
加第二路必须在创建任何 CAN client 之前调用 `tsmaster_can_init`。
"""

from __future__ import annotations

import os
import time

import _path  # noqa: F401

from preset_rs_native_sdk import (
    PRESET_CAN_BACKEND_TSMASTER,
    PRESET_TSMASTER_TC1016,
    PresetRSNativeSDK,
)


DLL_PATH = os.getcwd() + r"\preset_rs.dll"

# 按实车改：Tester -> ECU / ECU -> Tester / 功能寻址。
PHYSICAL_ID = 0x700
RESPONSE_ID = 0x708
FUNCTIONAL_ID = 0x7DF


def main() -> None:
    sdk = PresetRSNativeSDK(DLL_PATH)
    print(f"preset_rs {sdk.version}  ABI={sdk.abi_version}")

    can_cfg = sdk.tsmaster_default_config()
    # 逻辑通道 / 物理通道都是 0 起算：面板 CH1 = 0。
    can_cfg.application_channel = 0
    can_cfg.hardware_channel = 0
    can_cfg.hardware_index = 0
    can_cfg.device_type = PRESET_TSMASTER_TC1016  # 默认 TC1016，按实卡改
    can_cfg.is_fd = 1
    can_cfg.brs = 0
    can_cfg.nominal_bitrate = 500_000
    can_cfg.data_bitrate = 2_000_000

    uds_cfg = sdk.default_config()
    uds_cfg.physical_id = PHYSICAL_ID
    uds_cfg.response_id = RESPONSE_ID
    uds_cfg.functional_id = FUNCTIONAL_ID
    uds_cfg.is_fd = 1
    uds_cfg.raw_rx_enabled = 1

    # open 时传入 config，第一路就已经初始化完成。
    with sdk.tsmaster_open(can_cfg) as device:
        print(
            "backend:",
            device.backend,
            "(2 = TSMaster)" if device.backend == PRESET_CAN_BACKEND_TSMASTER else "",
        )

        # 若还要第二路，必须在 new_can_uds_client 之前：
        # ch2 = sdk.tsmaster_default_config()
        # ch2.application_channel = 1
        # ch2.hardware_channel = 1
        # ch2.device_type = can_cfg.device_type
        # ch2.is_fd = 1
        # device.tsmaster_can_init(ch2)

        with device.new_can_uds_client(can_cfg.application_channel, uds_cfg) as client:
            try:
                vin = client.request(bytes.fromhex("22 F1 80"), timeout_ms=1000)
                print("UDS 22 F1 80 ->", vin.hex(" "))
            except Exception as exc:
                print("UDS request:", exc)

            deadline = time.monotonic() + 2.0
            while time.monotonic() < deadline:
                for frame in client.try_read_ex():
                    direction = "TX" if frame.direction == 0 else "RX"
                    print(
                        f"{direction} time_us={frame.timestamp_us} "
                        f"id=0x{frame.id:X} dlc={frame.dlc} "
                        f"fd={frame.is_fd} brs={frame.brs} data={frame.payload.hex(' ')}"
                    )
                time.sleep(0.02)

            stats = client.rx_stats
            print(f"rx received={stats.received} dropped={stats.dropped} queued={stats.queued}")


def lin_example(sdk: PresetRSNativeSDK) -> None:
    """TSMaster LIN 主站。可与 CAN 共用同一把 device。"""
    lin_cfg = sdk.tsmaster_lin_default_config()
    lin_cfg.application_channel = 0
    lin_cfg.hardware_channel = 0
    lin_cfg.device_type = PRESET_TSMASTER_TC1016
    with sdk.tsmaster_lin_open(lin_cfg) as device:
        with device.new_lin_uds_client(lin_cfg.application_channel, nad=0x7E) as client:
            nad, payload = client.request(bytes.fromhex("22 F1 80"), timeout_ms=1000)
            print(f"LIN NAD=0x{nad:02X} ->", payload.hex(" "))


if __name__ == "__main__":
    main()
