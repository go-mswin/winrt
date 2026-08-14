package winrt

import (
	"errors"
	"testing"
)

func TestErrUnsupported(t *testing.T) {
	if ErrUnsupported == nil {
		t.Fatal("ErrUnsupported must not be nil")
	}
	if got := ErrUnsupported.Error(); got != "winrt: unsupported on this platform (windows only)" {
		t.Fatalf("ErrUnsupported.Error() = %q", got)
	}
	if !errors.Is(ErrUnsupported, ErrUnsupported) {
		t.Fatal("errors.Is(ErrUnsupported, ErrUnsupported) must hold")
	}
}

func TestApartmentConstants(t *testing.T) {
	if RoInitSingleThreaded != 0 || RoInitMultiThreaded != 1 {
		t.Fatalf("apartment constants drifted: STA=%d MTA=%d", RoInitSingleThreaded, RoInitMultiThreaded)
	}
}

func TestUserConsentVerifierAvailabilityString(t *testing.T) {
	cases := map[UserConsentVerifierAvailability]string{
		UserConsentVerifierAvailable:            "Available",
		UserConsentVerifierDeviceNotPresent:     "DeviceNotPresent",
		UserConsentVerifierNotConfiguredForUser: "NotConfiguredForUser",
		UserConsentVerifierDisabledByPolicy:     "DisabledByPolicy",
		UserConsentVerifierDeviceBusy:           "DeviceBusy",
		UserConsentVerifierAvailability(7):      "UserConsentVerifierAvailability(7)",
		UserConsentVerifierAvailability(-1):     "UserConsentVerifierAvailability(-1)",
	}
	for in, want := range cases {
		if got := in.String(); got != want {
			t.Errorf("UserConsentVerifierAvailability(%d).String() = %q, want %q", int32(in), got, want)
		}
	}
}

func TestUserConsentVerificationResultString(t *testing.T) {
	cases := map[UserConsentVerificationResult]string{
		UserConsentVerified:              "Verified",
		UserConsentDeviceNotPresent:      "DeviceNotPresent",
		UserConsentNotConfiguredForUser:  "NotConfiguredForUser",
		UserConsentDisabledByPolicy:      "DisabledByPolicy",
		UserConsentDeviceBusy:            "DeviceBusy",
		UserConsentRetriesExhausted:      "RetriesExhausted",
		UserConsentCanceled:              "Canceled",
		UserConsentVerificationResult(9): "UserConsentVerificationResult(9)",
	}
	for in, want := range cases {
		if got := in.String(); got != want {
			t.Errorf("UserConsentVerificationResult(%d).String() = %q, want %q", int32(in), got, want)
		}
	}
}

func TestItoa(t *testing.T) {
	cases := map[int32]string{
		0:           "0",
		7:           "7",
		-1:          "-1",
		42:          "42",
		-2147483648: "-2147483648",
		2147483647:  "2147483647",
	}
	for in, want := range cases {
		if got := itoa(in); got != want {
			t.Errorf("itoa(%d) = %q, want %q", in, got, want)
		}
	}
}
