package log

import (
	"fmt"
	"testing"
)

func Log(t *testing.T, s interface{}) {
	fmt.Printf("=== LOG   %s: %v\n", t.Name(), s)
}

func Logf(t *testing.T, format string, args ...interface{}) {
	Log(t, fmt.Sprintf(format, args...))
}
