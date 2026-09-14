"""Toomoss CAN / UDS 用法示例。

先用 `sdk.toomoss_open()` 打开设备，再 `toomoss_can_init()` 初始化通道，
然后用该通道创建 UDS client。不要直接构造 `PresetDevice` / `PresetCanUdsClient`。

运行前把 `DLL_PATH` 改成本机 `preset_rs.dll`，并接好 Toomoss 与 ECU。
"""

from __future__ import annotations

import os
import time

from preset_rs_native_sdk import (
    PRESET_CAN_BACKEND_TOOMOSS,
    PresetRSNativeSDK,
)


# DLL_PATH = os.environ.get("PRESET_RS_DLL", r"C:\path\to\preset_rs.dll")
DLL_PATH = os.getcwd() + r'\preset_rs.dll'

# 按实车改：Tester -> ECU / ECU -> Tester / 功能寻址。
PHYSICAL_ID = 0x73A
RESPONSE_ID = 0x7BA
FUNCTIONAL_ID = 0x7DF


def main() -> None:
    sdk = PresetRSNativeSDK(DLL_PATH)
    # sdk.set_print_log(True)
    # sdk.log_init("test log")
    print(f"preset_rs {sdk.version}  ABI={sdk.abi_version}")
    can_cfg = sdk.toomoss_default_config()
    # Toomoss 通道 0 起算：面板 CAN1/2/3/4 = channel 0/1/2/3。接 1 通道就写 0。
    can_cfg.channel = 0
    can_cfg.nominal_bitrate = 500_000
    can_cfg.data_bitrate = 2_000_000
    # can_cfg.brs = 1

    timing = sdk.toomoss_default_fd_timing()
    timing.nominal_brp = 1
    timing.nominal_tseg1 = 59
    timing.nominal_tseg2 = 20
    timing.nominal_sjw = 2
    timing.data_brp = 1
    timing.data_tseg1 = 14
    timing.data_tseg2 = 5
    timing.data_sjw = 2

    uds_cfg = sdk.default_config()
    uds_cfg.physical_id = PHYSICAL_ID
    uds_cfg.response_id = RESPONSE_ID
    uds_cfg.functional_id = FUNCTIONAL_ID
    uds_cfg.is_fd = 1
    uds_cfg.raw_rx_enabled = 1

    with sdk.toomoss_open() as device:
        print("backend:", device.backend, "(1 = Toomoss)" if device.backend == PRESET_CAN_BACKEND_TOOMOSS else "")
        device.toomoss_can_init(can_cfg)
        time.sleep(0.1)
        with device.new_can_uds_client(can_cfg.channel, uds_cfg) as client:
            # 读 VIN（UDS 0x22 F1 90）。无 ECU 时会超时，属正常现象。
            try:
                vin = client.request(bytes.fromhex("22 F1 80"), timeout_ms=1000)
                print("UDS 22 F1 80 ->", vin.hex(" "))
            except Exception as exc:
                print("UDS request:", exc)
            # 原始 CAN 发送 + 轮询接收（try_read 与 try_read_ex 不要混用）。
            client.write(0x123, bytes.fromhex("11 22 33 44"), is_fd=True)
            deadline = time.monotonic() + 2.0
            while time.monotonic() < deadline:
                
                for frame in client.try_read_ex():
                    direction = "TX" if frame.direction == 0 else "RX"
                    print(
                        f"{direction} id=0x{frame.id:X} dlc={frame.dlc} "
                        f"fd={frame.is_fd} brs={frame.brs} data={frame.payload.hex(' ')}"
                    )
                print(f"bus load={client.bus_load.load*100:.2f}%")
                time.sleep(0.02)
            
            stats = client.rx_stats
            print(f"rx received={stats.received} dropped={stats.dropped} queued={stats.queued}")
            

def lin_example(sdk: PresetRSNativeSDK) -> None:
    """同一把 Toomoss 设备上的 LIN 主站示例（需要 LIN 通道硬件）。"""
    lin_cfg = sdk.toomoss_lin_default_config()
    with sdk.toomoss_open() as device:
        device.toomoss_lin_init(lin_cfg)
        device.toomoss_lin_write(0, 0x3C, bytes.fromhex("01 02 03"))
        try:
            frame = device.toomoss_lin_read(0, 0x3D)
            print("LIN 0x3D ->", frame.payload.hex(" "))
        except Exception as exc:
            print("LIN read:", exc)


if __name__ == "__main__":
    main()
