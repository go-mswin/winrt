//go:build windows

package winrt

// Live, on-a-real-Windows-machine proof of the WinRT Hello glue. Gated behind
// an environment variable so headless CI skips it (per
// feedback-a-skip-is-not-a-pass, the skip fires ONLY on the outer env gate,
// never on an error from the code under test). Cross-compiled here
// (GOOS=windows go test -c), copied to the Win11 ARM64 QEMU VM and run there.
//
//	TestLiveAvailability — RoInitialize, resolve the UserConsentVerifier
//	activation factory (with the correct IID baked into the generated
//	projection), run CheckAvailabilityAsync to completion and decode the
//	availability enum. This proves the whole reference plumbing end to end. It
//	does NOT assert a specific availability value: a VM with no enrolled Hello
//	reports DeviceNotPresent/NotConfiguredForUser, which is a PASS for the
//	plumbing — what must hold is that the call returns without error.
//
// The interactive biometric prompt (RequestVerificationAsync) cannot be
// auto-verified without an enrolled Hello device, exactly like the macOS
// LAContext -34018 pattern; it is exercised only when a human is enrolled.

import (
	"os"
	"testing"
)

func TestLiveAvailability(t *testing.T) {
	if os.Getenv("WINRT_LIVE") != "1" {
		t.Skip("set WINRT_LIVE=1 to run the live UserConsentVerifier availability probe on Windows")
	}

	cleanup, err := initRuntime()
	if err != nil {
		t.Fatalf("initRuntime (RoInitialize): %v", err)
	}
	defer cleanup()

	avail, err := checkAvailability()
	if err != nil {
		t.Fatalf("checkAvailability (factory resolve + CheckAvailabilityAsync + decode): %v", err)
	}
	t.Logf("AVAILABILITY_OK factory resolved, async completed, decoded enum = %s (%d)", avail, int32(avail))

	ok, err := Available()
	if err != nil {
		t.Fatalf("Available: %v", err)
	}
	t.Logf("AVAILABLE_OK Available()=%v (true only when a Hello device is enrolled)", ok)
}

func TestLiveRequireUserConsent(t *testing.T) {
	if os.Getenv("WINRT_LIVE_PROMPT") != "1" {
		t.Skip("set WINRT_LIVE_PROMPT=1 to attempt the interactive Windows Hello prompt (needs an enrolled device)")
	}
	// When no Hello device is enrolled this returns (false, nil) WITHOUT
	// prompting; with an enrolled device it surfaces the biometric UI and
	// returns whether the user verified. Either way it must not error.
	ok, err := RequireUserConsent("Unlock your SSH key passphrase to connect to the Weft cluster")
	if err != nil {
		t.Fatalf("RequireUserConsent: %v", err)
	}
	t.Logf("REQUIRE_CONSENT_OK verified=%v (false when Hello is unavailable OR the user declined)", ok)
}
