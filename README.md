# preset_api

Windows amd64 Go bindings for `preset_rs.dll` ABI v12. The package supports
Toomoss, PCAN, TSMaster, and Vector CAN/LIN backends exposed by the native DLL.

```go
package main

import (
	"fmt"

	"github.com/LoveWonYoung/preset_api"
)

func main() {
	device, err := preset.NewToomoss()
	if err != nil {
		panic(err)
	}
	defer device.Close()

	if err := device.InitCAN(0); err != nil {
		panic(err)
	}

	runtime, err := device.NewCANRuntime(0)
	if err != nil {
		panic(err)
	}
	defer runtime.Shutdown()

	config := preset.DefaultCANDiagnosticConfig(0x73A, 0x7BA, 0x7DF)
	timeouts := preset.DefaultRequestTimeouts()
	info, response, err := runtime.UDSRequest(
		&config,
		&timeouts,
		[]byte{0x19, 0x02, 0x09},
		35_000,
	)
	fmt.Printf("outcome=%d response=% X err=%v\n", info.Outcome, response, err)
}
```

The default DLL path is `bin\preset_rs.dll`, relative to the process working
directory. Call `preset.SetDLLPath` before any other package function to use a
different location.

Close borrowed resources before their owner:

1. diagnostic requests, observers, and periodic senders;
2. CAN runtimes and LIN UDS clients;
3. the physical device.

The binding serializes each handle's close operation against calls using that
handle. A UDS negative response is a normal diagnostic outcome; inspect
`DiagnosticRequestInfo.Outcome` and `LastNRC` instead of treating it as a local
call failure.

Run validation on Windows amd64:

```powershell
gofmt -w .
go test ./...
go vet ./...
```

The example under `cmd/preset-example` requires a real adapter and the
corresponding vendor driver; ordinary unit tests do not.
