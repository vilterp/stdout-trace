package format

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompositeLines(t *testing.T) {
	// 4th channel spawning to 2nd channel
	testCases := []struct {
		lines []Line
		out   Line
	}{
		{
			lines: []Line{
				[]rune(" ╷│││ "),
				[]rune(" ╶─╴"),
			},
			out: []rune(" ┌┼┤│ "),
		},
		{
			lines: []Line{
				[]rune(" ││││ "),
				[]rune("   *"),
			},
			out: []rune(" ││*│ "),
		},
	}
	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("case %d", idx), func(t *testing.T) {
			actual := compositeLines(tc.lines)
			require.Equal(t, tc.out, actual)
		})
	}
}

func TestHorizLine(t *testing.T) {
	require.Equal(t, Line(" ╶───╴"), horizLine(1, 5))
}
