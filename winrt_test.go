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

// TestAvailabilityBecomesTheRightReason. The point of the translation is that
// somebody is told what to DO: no reader, no enrolment and a policy against it
// send a person to three different places.
func TestAvailabilityBecomesTheRightReason(t *testing.T) {
	for _, c := range []struct {
		in   UserConsentVerifierAvailability
		want UserConsentVerificationResult
	}{
		{UserConsentVerifierDeviceNotPresent, UserConsentDeviceNotPresent},
		{UserConsentVerifierNotConfiguredForUser, UserConsentNotConfiguredForUser},
		{UserConsentVerifierDisabledByPolicy, UserConsentDisabledByPolicy},
		{UserConsentVerifierDeviceBusy, UserConsentDeviceBusy},
		// Available means no prompt outcome exists yet, and something this
		// package has not been taught means the same: neither is an answer
		// from a person.
		{UserConsentVerifierAvailable, UserConsentDeviceNotPresent},
		{UserConsentVerifierAvailability(99), UserConsentDeviceNotPresent},
	} {
		if got := availabilityAsResult(c.in); got != c.want {
			t.Errorf("availabilityAsResult(%v) = %v, want %v", c.in, got, c.want)
		}
	}
}

// TestAskedSeparatesAPersonFromAMachine. Counting "no verifier on this
// computer" as a failure accuses somebody who was never shown a prompt.
func TestAskedSeparatesAPersonFromAMachine(t *testing.T) {
	for _, c := range []struct {
		in   UserConsentVerificationResult
		want bool
	}{
		{UserConsentVerified, true},
		{UserConsentCanceled, true},
		{UserConsentRetriesExhausted, true},
		{UserConsentDeviceNotPresent, false},
		{UserConsentNotConfiguredForUser, false},
		{UserConsentDisabledByPolicy, false},
		{UserConsentDeviceBusy, false},
		{UserConsentVerificationResult(99), false},
	} {
		if got := c.in.Asked(); got != c.want {
			t.Errorf("%v.Asked() = %v, want %v", c.in, got, c.want)
		}
	}
}
