// Package ui holds the winrt-go-gen–generated Go projection of the WinRT class
// Windows.Security.Credentials.UI.UserConsentVerifier (Windows Hello). The
// generated file (userconsentverifier.go) is committed so the build is
// reproducible offline; regenerate it with:
//
// The generated code is //go:build windows only; this doc file carries no build
// constraint so the package is a valid (empty) package on other GOOS, letting
// `go build ./...` and `go vet ./...` succeed everywhere.
//
//go:generate go run github.com/saltosystems/winrt-go/cmd/winrt-go-gen -class Windows.Security.Credentials.UI.UserConsentVerifier -method-filter CheckAvailabilityAsync -method-filter RequestVerificationAsync -method-filter !*
package ui
