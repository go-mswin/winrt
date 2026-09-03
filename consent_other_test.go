//go:build !windows

package winrt

import (
	"errors"
	"testing"
)

func TestAvailableUnsupported(t *testing.T) {
	ok, err := Available()
	if ok {
		t.Errorf("Available() = true off Windows, want false")
	}
	if !errors.Is(err, ErrUnsupported) {
		t.Errorf("Available() err = %v, want ErrUnsupported", err)
	}
}

func TestRequireUserConsentUnsupported(t *testing.T) {
	ok, err := RequireUserConsent("unlock the vault")
	if ok {
		t.Errorf("RequireUserConsent() = true off Windows, want false")
	}
	if !errors.Is(err, ErrUnsupported) {
		t.Errorf("RequireUserConsent() err = %v, want ErrUnsupported", err)
	}
}

// TestAvailabilityUnsupported and TestVerifyUnsupported keep the new entry
// points honest off Windows: there is no Windows Runtime here, which is not
// the same as a person failing to authenticate.
func TestAvailabilityUnsupported(t *testing.T) {
	got, err := Availability()
	if !errors.Is(err, ErrUnsupported) {
		t.Errorf("Availability() error = %v, want ErrUnsupported", err)
	}
	if got != 0 {
		t.Errorf("Availability() = %v, want the zero value", got)
	}
}

func TestVerifyUnsupported(t *testing.T) {
	got, err := Verify("unlock the vault")
	if !errors.Is(err, ErrUnsupported) {
		t.Errorf("Verify() error = %v, want ErrUnsupported", err)
	}
	if got != 0 {
		t.Errorf("Verify() = %v, want the zero value", got)
	}
}
