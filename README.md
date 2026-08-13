# winrt — pure-Go WinRT interop

`github.com/go-mswin/winrt` is a pure-Go **CGO=0** layer for calling Windows
Runtime (WinRT) APIs from Go, built on [`go-mswin/win32`](https://github.com/go-mswin/win32)'s
combase binding. Targets 64-bit Windows (amd64 + arm64). No cgo.

Provides the WinRT plumbing the fleet currently hand-rolls:

- `RoInitialize` / `RoUninitialize` (apartment init)
- `RoGetActivationFactory` / `RoActivateInstance`
- fast-pass `HSTRING` via `WindowsCreateStringReference` (+ `WindowsDeleteString`)
- `IInspectable` / `IUnknown` vtable call helpers + GUID/IID handling

## Consumers
The weft Windows apps' **Windows Hello** (`UserConsentVerifier`) plumbing, and any
future WinRT surface. One owned interop layer instead of duplicated combase hand-rolls.

## License
BSD-3-Clause — copyright the go-mswin authors.
