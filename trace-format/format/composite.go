package format

import (
	"fmt"
	"strconv"
	"strings"
)

type Line []rune

const (
	Log    = '*'
	Finish = 'X'
	Empty  = ' '

	VertLine     = '│'
	VertDownHalf = '╷'

	HorizFull  = '─'
	HorizRight = '╶'
	HorizLeft  = '╴'

	Cross    = '┼'
	OutLeft  = '┤'
	OutRight = '├'

	DownLeft  = '┐'
	DownRight = '┌'
)

func compositeChars(a rune, b rune) rune {
	if a == Empty {
		return b
	}
	if b == Empty {
		return a
	}
	if b == Log || b == Finish {
		return b
	}
	if a == VertLine {
		switch b {
		case HorizFull:
			return Cross
		case HorizLeft:
			return OutLeft
		case HorizRight:
			return OutRight
		case Log:
			return Log // TODO: compose?
		case Finish:
			return Finish
		}
	}
	if a == VertDownHalf {
		switch b {
		case HorizLeft:
			return DownLeft
		case HorizRight:
			return DownRight
		}
	}
	panic(fmt.Sprintf("unknown combo %v %v", strconv.QuoteRune(a), strconv.QuoteRune(b))) // I'm lazy
}

// left to right => bottom to top
func compositeLines(lines []Line) Line {
	maxLen := maxLength(lines)
	out := Line(strings.Repeat(" ", maxLen))
	for idx := range out {
		for _, line := range lines {
			if idx >= len(line) {
				continue
			}
			out[idx] = compositeChars(out[idx], line[idx])
		}
	}
	return out
}

func (l Line) String() string {
	return string(l)
}

func (l Line) spaceOut() Line {
	var out Line
	for _, r := range l {
		a, b := spaceOut(r)
		out = append(out, a, b)
	}
	return out
}

func spaceOut(r rune) (rune, rune) {
	switch r {
	case HorizFull:
		fallthrough
	case HorizRight:
		fallthrough
	case Cross:
		fallthrough
	case OutRight:
		fallthrough
	case DownRight:
		return r, HorizFull
	default:
		return r, Empty
	}
}

func horizLine(leftIdx int, rightIdx int) Line {
	var out Line
	for i := 0; i < leftIdx; i++ {
		out = append(out, Empty)
	}
	out = append(out, HorizRight)
	for i := 0; i < rightIdx-leftIdx-1; i++ {
		out = append(out, HorizFull)
	}
	out = append(out, HorizLeft)
	return out
}

func maxLength(lines []Line) int {
	maxLen := 0
	for _, s := range lines {
		if len(s) > maxLen {
			maxLen = len(s)
		}
	}
	return maxLen
}
