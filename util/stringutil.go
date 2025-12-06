package util

import "fmt"

// Quote wraps a string in double quotes.
// Corresponds to C++ util::quote().
func Quote(s string) string {
	return fmt.Sprintf("'%s'", s)
}
