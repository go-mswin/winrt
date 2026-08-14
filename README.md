# winrt — pure-Go Windows Hello (WinRT) for the fleet

`github.com/go-mswin/winrt` is a thin, pure-Go **CGO=0** wrapper over the
reference WinRT projection [`github.com/saltosystems/winrt-go`](https://github.com/saltosystems/winrt-go)
(itself built on the pure-Go COM/WinRT runtime
[`github.com/go-ole/go-ole`](https://github.com/go-ole/go-ole)). It targets
64-bit Windows (amd64 + arm64). No cgo.

It **does not reimplement WinRT**. `RoInitialize` / `RoGetActivationFactory`,
`HSTRING` creation, the `IInspectable`/`IUnknown` vtable calls and the
`IAsyncOperation` machinery all come from winrt-go + go-ole, the maintained
reference libraries (per the fleet's "use the reference library" rule). This
module adds only:

- a **committed, `winrt-go-gen`–generated projection** of
  `Windows.Security.Credentials.UI.UserConsentVerifier` (Windows Hello), so the
  build is reproducible offline — see `internal/ui` (regenerate with the
  `//go:generate` directive there);
- two ergonomic helpers that run the async operations to completion and decode
  the WinRT enums into Go values:

```go
// Available reports whether Windows Hello can prompt (never shows UI).
ok, err := winrt.Available()

// RequireUserConsent surfaces the Hello prompt and reports whether the user
// verified. Returns (false, nil) — not an error — when Hello is unavailable,
// so callers treating biometry as a best-effort factor can proceed.
verified, err := winrt.RequireUserConsent("Unlock your SSH key passphrase")
```

The OS-independent core (the WinRT enum types + `String` methods, apartment
constants, `ErrUnsupported`) builds on every GOOS and is covered to 100%; the
live glue is `//go:build windows`. Off Windows, `Available` /
`RequireUserConsent` return `ErrUnsupported`.

## Verified
On the Win11 ARM64 QEMU VM: `RoInitialize` OK, the `UserConsentVerifier`
activation factory resolves (with the correct metadata IID
`af4f3f91-564c-4ddc-b8b5-973447627c65`), `CheckAvailabilityAsync` runs to
completion and decodes its enum. The interactive biometric prompt needs an
enrolled Hello device and is not headless-testable (the LAContext `-34018`
pattern).

## Consumers
`openweft/weft-app-windows` and `openweft/weft-loom-app-windows` — one owned,
correct Hello gate instead of duplicated hand-rolled combase calls.

## License
BSD-3-Clause — copyright the go-mswin authors.
