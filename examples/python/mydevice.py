"""MyDevice 二次封装用法示例。"""

from __future__ import annotations

import _path  # noqa: F401

from my_device import Backend, MyDevice

PHYSICAL_ID = 0x700
RESPONSE_ID = 0x708
FUNCTIONAL_ID = 0x7DF


def main() -> None:
    with MyDevice(
        Backend.Toomoss,
        PHYSICAL_ID,
        RESPONSE_ID,
        FUNCTIONAL_ID,
        channels=[0],
    ) as device:
        try:
            vin = device.request(bytes.fromhex("22 F1 80"), timeout_ms=1000)
            print("UDS 22 F1 80 ->", vin.hex(" "))
        except Exception as exc:
            print("UDS request:", exc)

        for frame in device.rxfn(timeout_ms=2000):
            print(
                f"{frame.direction} time_us={frame.timestamp_us} "
                f"id=0x{frame.id:X} dlc={frame.dlc} "
                f"fd={int(frame.is_fd)} brs={int(frame.brs)} data={frame.data.hex(' ')}"
            )


if __name__ == "__main__":
    main()
