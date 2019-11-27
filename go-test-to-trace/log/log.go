package log

import (
	"fmt"
	"testing"
)

func Log(t *testing.T, s string) {
	fmt.Printf("=== LOG   %s: %s\n", t.Name(), s)
}

func Logf(t *testing.T, format string, args ...interface{}) {
	Log(t, fmt.Sprintf(format, args...))
}
