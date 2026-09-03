//go:build !windows

package winrt

// Available reports whether Windows Hello is ready. Off Windows there is no
// Windows Runtime, so it always returns ErrUnsupported.
func Available() (bool, error) { return false, ErrUnsupported }

// RequireUserConsent surfaces the Windows Hello prompt. Off Windows there is no
// Windows Runtime, so it always returns ErrUnsupported.
func RequireUserConsent(reason string) (bool, error) {
	_ = reason
	return false, ErrUnsupported
}

// Availability reports what UserConsentVerifier says. Off Windows there is no
// Windows Runtime, so it always returns ErrUnsupported.
func Availability() (UserConsentVerifierAvailability, error) { return 0, ErrUnsupported }

// Verify surfaces the Windows Hello prompt and returns what Windows said. Off
// Windows there is no Windows Runtime, so it always returns ErrUnsupported.
func Verify(reason string) (UserConsentVerificationResult, error) {
	_ = reason
	return 0, ErrUnsupported
}
