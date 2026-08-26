//go:build !unix

package spin

import "os"

// hushStdin is a no-op on platforms without termios support.
func hushStdin(*os.File) (restore func(), err error) {
	return func() {}, nil
}
