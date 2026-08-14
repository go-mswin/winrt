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
