// Package winrt is a thin, pure-Go (CGO_ENABLED=0) wrapper over the reference
// Windows Runtime projection github.com/saltosystems/winrt-go (built on the
// pure-Go COM/WinRT runtime github.com/go-ole/go-ole).
//
// It does NOT reimplement WinRT: the low-level plumbing — RoInitialize /
// RoGetActivationFactory, HSTRING creation, IInspectable/IUnknown vtable calls
// and the IAsyncOperation machinery — comes entirely from winrt-go + go-ole,
// the maintained reference libraries. This package adds only a small,
// language-idiomatic layer on top: a committed, winrt-go-gen–generated
// projection of Windows.Security.Credentials.UI.UserConsentVerifier (see
// internal/ui) and two ergonomic helpers, [Available] and [RequireUserConsent],
// that run the availability check and the Windows Hello verification prompt to
// completion and decode the WinRT enums into Go values.
//
// Layout mirrors the go-mswin/win32 and go-macos/objc convention: the
// OS-independent core — the WinRT enum types and their String methods, the
// apartment constants and the [ErrUnsupported] sentinel — lives in UNTAGGED
// files exercised to 100% statement coverage on every GOOS; the live glue that
// reaches combase.dll lives in //go:build windows files and is proven on a real
// Windows machine (the Win11 ARM64 QEMU VM). Non-windows GOOS build the core
// and a stub whose [Available]/[RequireUserConsent] return [ErrUnsupported].
package winrt
