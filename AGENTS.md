# 项目维护指南

## 项目定位

- 这是 Windows amd64 下调用 `bin/preset_rs.dll` 的 Go 示例/绑定（目前为 `package main`，不是可直接 import 的库）。
- `preset_rs.h` 是原生 ABI 的权威定义，当前版本为 `PRESET_RS_ABI_VERSION 12`。不要根据猜测修改函数参数、返回值或结构体布局。
- `main.go` 是依赖真实设备的演示；不要在 CI 或普通单元测试中运行它。

## 代码组织

- `device.go`：共同的 `Device` 句柄、自动驱动以及 PCAN/TSMaster/Vector CAN 初始化。
- `toomoss_api.go`：Toomoss LIN 和设备关闭。
- `can.go`：CAN/CAN-FD、诊断运行时、观察者、UDS 请求。
- `lin.go`：各厂商 LIN 配置、LIN 主从读写和 LIN UDS。
- `extras.go`：能力查询、日志和周期 CAN 发送。
- `abi_windows_amd64.go`：Windows x64 C 结构体返回规则；`can_layout_test.go`：结构体大小和偏移回归测试。

## ABI 与资源安全

- 新增绑定时先核对 `preset_rs.h` 的 C 签名、类型宽度、对齐、零基通道含义与返回码。C `size_t` 和指针在 amd64 上用 `uintptr`；C `uint8_t` 用 `byte`。
- C 的结构体按字段顺序声明；增加结构体时在 `can_layout_test.go` 中加入大小和关键偏移测试。返回结构体大于 8 字节时 Windows x64 使用隐藏的结果指针。
- 原生状态码为负数表示错误，调用 `presetStatus`/`nativeStatus` 并保留 `LastError()`；不要把 UDS 阴性响应误当成本地调用失败。
- `Device` 拥有物理设备；runtime、observer、periodic、request 和 LIN UDS client 借用其句柄。按 request/observer/periodic → runtime、LIN client → device 的顺序关闭。不要在同一句柄执行原生调用的同时关闭它。
- 新厂商优先复用 `Device` 上的通用 CAN/LIN 方法；设备初始化方法应保持厂商专属。仅在 `preset_rs.h` 中有声明且 DLL 实际导出时添加函数。
- 不要修改或删除已提交的 DLL；它们是本项目运行时依赖，硬件测试需安装对应厂商驱动。

## 验证

在 Windows amd64 运行 `gofmt -w <修改的 Go 文件>`、`go test ./...` 和 `go vet ./...`。单元测试应可在无硬件环境运行；设备行为需要单独用硬件验证。
