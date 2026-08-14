//go:build windows

// Live WinRT glue for Windows Hello (UserConsentVerifier). RoInitialize,
// activation-factory resolution, HSTRING creation and the IAsyncOperation
// machinery all come from github.com/saltosystems/winrt-go + github.com/go-ole/
// go-ole; this file only orchestrates them and decodes the result enums. The
// one thing go-ole (v1.3.0) does not expose is RoUninitialize, so it is bound
// directly off combase.dll — a single OS proc, not a reinvented runtime.
package winrt

import (
	"fmt"
	"runtime"
	"time"
	"unsafe"

	"github.com/go-mswin/winrt/internal/ui"
	ole "github.com/go-ole/go-ole"
	"github.com/saltosystems/winrt-go/windows/foundation"
	"golang.org/x/sys/windows"
)

// Benign HRESULTs from RoInitialize: S_FALSE means the apartment was already
// initialized on this thread (we still hold a ref → must RoUninitialize);
// RPC_E_CHANGED_MODE means COM was already initialized in a different apartment
// (we did NOT take a ref → must NOT RoUninitialize).
const (
	sFalse          = 0x00000001
	rpcEChangedMode = 0x80010106
)

// asyncWaitTimeout bounds how long the helpers wait for an IAsyncOperation. The
// availability check resolves near-instantly; the verification prompt waits for
// the interactive user, so this is generous.
const asyncWaitTimeout = 90 * time.Second

var (
	modCombase         = windows.NewLazySystemDLL("combase.dll")
	procRoUninitialize = modCombase.NewProc("RoUninitialize")
)

// initRuntime pins the goroutine to its OS thread (so every WinRT call in the
// operation runs on the RoInitialize'd thread) and initializes the
// multithreaded apartment. It returns a cleanup that undoes exactly what it
// did.
func initRuntime() (cleanup func(), err error) {
	runtime.LockOSThread()
	e := ole.RoInitialize(uint32(RoInitMultiThreaded))
	if e == nil {
		return func() { procRoUninitialize.Call(); runtime.UnlockOSThread() }, nil
	}
	code := uintptr(0)
	if oe, ok := e.(*ole.OleError); ok {
		code = oe.Code()
	}
	switch code {
	case sFalse:
		return func() { procRoUninitialize.Call(); runtime.UnlockOSThread() }, nil
	case rpcEChangedMode:
		return func() { runtime.UnlockOSThread() }, nil
	default:
		runtime.UnlockOSThread()
		return nil, fmt.Errorf("winrt: RoInitialize: %w", e)
	}
}

// awaitOperation drives an IAsyncOperation to completion by polling its
// IAsyncInfo status (obtained via QueryInterface), then returns the raw result
// word from GetResults. Polling avoids the parameterized
// AsyncOperationCompletedHandler IID, which WinRT computes per result type.
func awaitOperation(op *foundation.IAsyncOperation) (uintptr, error) {
	unk := (*ole.IUnknown)(unsafe.Pointer(op))
	infoDisp, err := unk.QueryInterface(ole.NewGUID(foundation.GUIDIAsyncInfo))
	if err != nil {
		return 0, fmt.Errorf("winrt: QueryInterface(IAsyncInfo): %w", err)
	}
	info := (*foundation.IAsyncInfo)(unsafe.Pointer(infoDisp))
	defer info.Release()

	deadline := time.Now().Add(asyncWaitTimeout)
	for {
		st, err := info.GetStatus()
		if err != nil {
			return 0, fmt.Errorf("winrt: IAsyncInfo.GetStatus: %w", err)
		}
		switch st {
		case foundation.AsyncStatusCompleted:
			res, err := op.GetResults()
			if err != nil {
				return 0, fmt.Errorf("winrt: GetResults: %w", err)
			}
			return uintptr(res), nil
		case foundation.AsyncStatusError:
			return 0, fmt.Errorf("winrt: async operation ended with status Error")
		case foundation.AsyncStatusCanceled:
			return 0, fmt.Errorf("winrt: async operation was canceled")
		}
		if time.Now().After(deadline) {
			return 0, fmt.Errorf("winrt: async operation timed out after %s", asyncWaitTimeout)
		}
		time.Sleep(15 * time.Millisecond)
	}
}

func release(p unsafe.Pointer) {
	if p != nil {
		(*ole.IUnknown)(p).Release()
	}
}

// checkAvailability runs UserConsentVerifier.CheckAvailabilityAsync to
// completion and decodes the availability enum.
func checkAvailability() (UserConsentVerifierAvailability, error) {
	op, err := ui.UserConsentVerifierCheckAvailabilityAsync()
	if err != nil {
		return 0, fmt.Errorf("winrt: CheckAvailabilityAsync: %w", err)
	}
	defer release(unsafe.Pointer(op))
	res, err := awaitOperation(op)
	if err != nil {
		return 0, err
	}
	return UserConsentVerifierAvailability(uint32(res)), nil
}

// Available reports whether Windows Hello (UserConsentVerifier) is present,
// enrolled and ready to prompt on this machine. It initializes the WinRT
// runtime, runs CheckAvailabilityAsync to completion and returns true only for
// UserConsentVerifierAvailable. It never displays UI.
func Available() (bool, error) {
	cleanup, err := initRuntime()
	if err != nil {
		return false, err
	}
	defer cleanup()
	avail, err := checkAvailability()
	if err != nil {
		return false, err
	}
	return avail == UserConsentVerifierAvailable, nil
}

// RequireUserConsent surfaces the Windows Hello prompt with the given reason and
// reports whether the user proved their identity.
//
// It returns (false, nil) — not an error — when Windows Hello is unavailable
// (not present, not enrolled, disabled by policy, busy): callers that treat
// biometry as a best-effort second factor can proceed on their primary gate.
// It returns (true, nil) when the user is verified, (false, nil) when the user
// is present but declines/cancels/exhausts retries, and a non-nil error only on
// a runtime failure. It displays system UI and must run on an interactive
// session.
func RequireUserConsent(reason string) (bool, error) {
	cleanup, err := initRuntime()
	if err != nil {
		return false, err
	}
	defer cleanup()

	avail, err := checkAvailability()
	if err != nil {
		return false, err
	}
	if avail != UserConsentVerifierAvailable {
		return false, nil
	}

	op, err := ui.UserConsentVerifierRequestVerificationAsync(reason)
	if err != nil {
		return false, fmt.Errorf("winrt: RequestVerificationAsync: %w", err)
	}
	defer release(unsafe.Pointer(op))
	res, err := awaitOperation(op)
	if err != nil {
		return false, err
	}
	return UserConsentVerificationResult(uint32(res)) == UserConsentVerified, nil
}
