# preset_api

`preset_api` 是 `preset_rs.dll` 的纯 Go 绑定，仅支持 Windows AMD64。实现只使用
标准库 `syscall`，不使用 cgo。

仓库根目录是 Go 包；`python/` 是 ctypes 绑定和 `MyDevice` 二次封装；可运行示例在
`examples/go/` 与 `examples/python/`。

```powershell
go run ./examples/go/mydevice
python examples/python/mydevice.py
python examples/python/tsmaster_rx.py
```

## 使用

程序必须先用明确路径加载 DLL。相对路径会在调用 Windows `LoadLibrary` 之前转换为
绝对路径，并用安全搜索标志加载依赖 DLL（目标 DLL 目录及 Windows 默认安全目录），
不依赖当前目录或 `PATH`。加载时会校验 ABI 版本，当前仅接受 ABI v6。

```go
package main

import (
	"fmt"
	"log"

	preset "preset_api"
)

func main() {
	if err := preset.LoadDLL(`C:\path\to\preset_rs.dll`); err != nil {
		log.Fatal(err)
	}

	version, err := preset.Version()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(version)

	config, err := preset.PcanDefaultConfig()
	if err != nil {
		log.Fatal(err)
	}
	device, status := preset.PcanOpen(&config)
	if status != preset.PRESET_OK {
		log.Fatalf("open failed: %d (%s)", status, preset.LastErrorString())
	}
	defer preset.DeviceClose(&device)
}
```

ABI v6 提供自动 CAN 后端选择、显式 CAN-FD 时序、运行时 BRS/发送回显、带方向和
DLC 元数据的原始帧读取，以及 PCAN、TSMaster、Vector 的 LIN 主站接口。例如：

```go
autoConfig, _ := preset.AutoDefaultConfig()
device, status := preset.AutoOpen(&autoConfig)
if status != preset.PRESET_OK {
	log.Fatal(preset.LastErrorString())
}
backend, _ := preset.DeviceGetBackend(device)
fmt.Println("CAN backend:", backend)
preset.DeviceClose(&device)

linConfig, _ := preset.PcanLinDefaultConfig()
linDevice, status := preset.PcanLinOpen(&linConfig)
if status == preset.PRESET_OK {
	frame, readStatus := preset.LinMasterRead(linDevice, linConfig.Channel, 0x3d, 100)
	fmt.Println(frame, readStatus)
	preset.DeviceClose(&linDevice)
}
```

`LinUdsClientNew` 可用于所有支持的 LIN 后端。旧的 `ToomossLinUdsClientNew` Go 名称
仍保留为兼容别名，但内部调用 ABI v6 的通用构造函数。`CanUdsTryReadEx`
返回统一为微秒的厂商硬件/驱动单调时间 `TimestampUS`；不同后端的时间原点可能不同，
该值适合排序和计算间隔，不是系统时间。值为 0 表示硬件时间不可用，例如软件生成的
TX echo。所有原始帧读取函数会消费同一个队列，不要在同一个 client 上混用。

所有返回 `status` 的函数沿用 `preset_rs.h` 中的 `PRESET_OK` 和
`PRESET_ERR_*`。请求与批量读取函数返回实际写入数量；若输出切片太小，会返回
`PRESET_ERR_BUFFER_TOO_SMALL`，同时数量为所需容量。DLL 加载或符号解析失败会返回
`PRESET_ERR_TRANSPORT`，具体原因由 `DLLLoadError()` 提供。

## 生命周期与并发

- `PresetDevice`、`PresetCanUdsClient` 和 `PresetLinUdsClient` 是共享状态的句柄值；
  复制后仍指向同一个原生句柄，任意副本成功关闭后，其他副本的 `IsClosed()` 也会
  返回 `true`。
- 对同一个句柄，普通调用和 Close 已做读写锁同步，Close 不会与正在执行的 DLL
  调用并发释放 Rust 对象。仍建议按“先关闭 UDS client，再关闭 device”的顺序显式
  释放资源。
- Rust 原生句柄只是 DLL 地址值，不需要调用方使用 `runtime.KeepAlive`。配置、切片和
  输出结构等 Go 内存会在 DLL 调用的同步期间由绑定层统一 pin；DLL 不得在函数返回后
  保存这些 Go 指针。

`preset_rs` 的同步错误信息使用 Windows 线程局部存储。需要让错误信息和失败调用严格
对应时，使用 `WithLastError`：

```go
status, message := preset.WithLastError(func() int32 {
	return preset.CanUdsSetDefaultSTMin(client, 1)
})
```

CAN worker 的异步错误使用 `CanUdsLastErrorString`。

## 测试

```powershell
$env:CGO_ENABLED = '0'
go test ./...

# 可选的 DLL ABI 冒烟测试
$env:PRESET_RS_DLL = 'C:\path\to\preset_rs.dll'
go test -v ./...

# 使用 Python ctypes 检查 ABI v6 结构体、默认值及 CAN/LIN 调用签名
python python/ctypes_smoke_test.py $env:PRESET_RS_DLL
```

从非 Windows 主机只做编译检查：

```sh
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go test -c
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go vet -unsafeptr=false ./...
```

`vet` 仅关闭 `unsafeptr` 检查，因为 `preset_version` 按 C ABI 返回 DLL 静态字符串
地址；从 `uintptr` 恢复该非 Go 指针是这里有意且必要的操作。
