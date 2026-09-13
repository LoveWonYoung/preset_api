# preset_api

`preset_api` 是 `preset_rs.dll` 的纯 Go Windows 绑定。实现只使用标准库
`syscall`，不使用 cgo。

```go
package main

import (
	"fmt"
	"log"

	preset "preset_api"
)

func main() {
	// DLL 位于 Windows DLL 搜索路径时可以省略 LoadDLL。
	if err := preset.LoadDLL(`C:\path\to\preset_rs.dll`); err != nil {
		log.Fatal(err)
	}

	fmt.Println(preset.Version(), preset.ABIVersion())
	config := preset.PcanDefaultConfig()
	device, status := preset.PcanOpen(&config)
	if status != preset.PRESET_OK {
		log.Fatalf("open failed: %d", status)
	}
	defer preset.DeviceClose(&device)
}
```

所有返回 `status` 的函数沿用 `preset_rs.h` 中的 `PRESET_OK` 和
`PRESET_ERR_*`。请求与批量读取函数返回实际写入数量；若输出切片太小，
会返回 `PRESET_ERR_BUFFER_TOO_SMALL`，同时数量为所需容量。

`preset_rs` 的同步错误信息是 Windows 线程局部数据。若需可靠调用
`LastError`/`LastErrorString`，请用 `runtime.LockOSThread` 包住失败的 API
调用和错误读取；CAN worker 的异步错误可直接用 `CanUdsLastErrorString`。

运行测试：

```powershell
$env:CGO_ENABLED = '0'
go test ./...

# 可选的 DLL ABI 冒烟测试
$env:PRESET_RS_DLL = 'C:\path\to\preset_rs.dll'
go test -v ./...
```
