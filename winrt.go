package winrt

import "errors"

// ErrUnsupported is returned by the entry points on non-Windows platforms,
// where there is no Windows Runtime. It is stable and may be tested with
// errors.Is.
var ErrUnsupported = errors.New("winrt: unsupported on this platform (windows only)")

// Apartment selects the COM/WinRT threading model passed to RoInitialize. The
// helpers in this package use the multithreaded apartment, which needs no
// message pump.
type Apartment uint32

const (
	// RoInitSingleThreaded is RO_INIT_SINGLETHREADED (a single-threaded
	// apartment, requiring a running message pump).
	RoInitSingleThreaded Apartment = 0
	// RoInitMultiThreaded is RO_INIT_MULTITHREADED (the multithreaded
	// apartment). Async completions arrive on a thread-pool thread, so a
	// polling waiter needs no message pump.
	RoInitMultiThreaded Apartment = 1
)

// UserConsentVerifierAvailability mirrors the WinRT
// Windows.Security.Credentials.UI.UserConsentVerifierAvailability enum reported
// by UserConsentVerifier.CheckAvailabilityAsync. Only [UserConsentVerifierAvailable]
// means Windows Hello can prompt the user.
type UserConsentVerifierAvailability int32

const (
	// UserConsentVerifierAvailable — a verifier device is present, enrolled and
	// ready.
	UserConsentVerifierAvailable UserConsentVerifierAvailability = 0
	// UserConsentVerifierDeviceNotPresent — no verifier device is present.
	UserConsentVerifierDeviceNotPresent UserConsentVerifierAvailability = 1
	// UserConsentVerifierNotConfiguredForUser — the user has not enrolled
	// Windows Hello.
	UserConsentVerifierNotConfiguredForUser UserConsentVerifierAvailability = 2
	// UserConsentVerifierDisabledByPolicy — administrative policy disabled the
	// verifier.
	UserConsentVerifierDisabledByPolicy UserConsentVerifierAvailability = 3
	// UserConsentVerifierDeviceBusy — the verifier device is busy.
	UserConsentVerifierDeviceBusy UserConsentVerifierAvailability = 4
)

// String returns the enumerator name.
func (a UserConsentVerifierAvailability) String() string {
	switch a {
	case UserConsentVerifierAvailable:
		return "Available"
	case UserConsentVerifierDeviceNotPresent:
		return "DeviceNotPresent"
	case UserConsentVerifierNotConfiguredForUser:
		return "NotConfiguredForUser"
	case UserConsentVerifierDisabledByPolicy:
		return "DisabledByPolicy"
	case UserConsentVerifierDeviceBusy:
		return "DeviceBusy"
	default:
		return "UserConsentVerifierAvailability(" + itoa(int32(a)) + ")"
	}
}

// UserConsentVerificationResult mirrors the WinRT
// Windows.Security.Credentials.UI.UserConsentVerificationResult enum returned
// by UserConsentVerifier.RequestVerificationAsync. Only [UserConsentVerified]
// means the user proved their identity.
type UserConsentVerificationResult int32

const (
	// UserConsentVerified — the user was successfully verified.
	UserConsentVerified UserConsentVerificationResult = 0
	// UserConsentDeviceNotPresent — no verifier device is present.
	UserConsentDeviceNotPresent UserConsentVerificationResult = 1
	// UserConsentNotConfiguredForUser — the user has not enrolled Windows Hello.
	UserConsentNotConfiguredForUser UserConsentVerificationResult = 2
	// UserConsentDisabledByPolicy — administrative policy disabled the verifier.
	UserConsentDisabledByPolicy UserConsentVerificationResult = 3
	// UserConsentDeviceBusy — the verifier device is busy.
	UserConsentDeviceBusy UserConsentVerificationResult = 4
	// UserConsentRetriesExhausted — the user exhausted the allowed attempts.
	UserConsentRetriesExhausted UserConsentVerificationResult = 5
	// UserConsentCanceled — the verification was canceled.
	UserConsentCanceled UserConsentVerificationResult = 6
)

// String returns the enumerator name.
func (r UserConsentVerificationResult) String() string {
	switch r {
	case UserConsentVerified:
		return "Verified"
	case UserConsentDeviceNotPresent:
		return "DeviceNotPresent"
	case UserConsentNotConfiguredForUser:
		return "NotConfiguredForUser"
	case UserConsentDisabledByPolicy:
		return "DisabledByPolicy"
	case UserConsentDeviceBusy:
		return "DeviceBusy"
	case UserConsentRetriesExhausted:
		return "RetriesExhausted"
	case UserConsentCanceled:
		return "Canceled"
	default:
		return "UserConsentVerificationResult(" + itoa(int32(r)) + ")"
	}
}

// itoa formats a signed 32-bit integer without importing strconv, keeping the
// OS-independent core dependency-free.
func itoa(v int32) string {
	if v == 0 {
		return "0"
	}
	neg := v < 0
	u := uint32(v)
	if neg {
		u = uint32(-v)
	}
	var buf [12]byte
	i := len(buf)
	for u > 0 {
		i--
		buf[i] = byte('0' + u%10)
		u /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
