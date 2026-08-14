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
