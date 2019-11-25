package tracer

import (
	"fmt"
	"strings"
)

type Formatter struct {
	spanChannels []int
	openSpans    map[int]*Span
}

func NewFormatter() *Formatter {
	return &Formatter{
		openSpans: map[int]*Span{},
	}
}

func (f *Formatter) Handle(evt *TraceEvent) {
	switch evt.TraceEvent {
	case startSpanEvt:
		f.openSpans[evt.SpanID] = &Span{
			startedAt: evt.Timestamp,
			operation: evt.Operation,
			id:        evt.SpanID,
			parentID:  evt.ParentID,
		}
		f.addSpanToChannels(evt.SpanID)
		f.logLeftTrack(evt.SpanID, evt.ParentID, startSpanEvt)
		fmt.Print("\t")
		fmt.Printf("start: %s\n", evt.Operation)
	case logEvt:
		f.logLeftTrack(evt.SpanID, -1, logEvt)
		fmt.Print("\t")
		fmt.Println(evt.LogLine)
	case finishSpanEvt:
		span := f.openSpans[evt.SpanID]
		span.finishedAt = &evt.Timestamp
		f.logLeftTrack(evt.SpanID, -1, finishSpanEvt)
		fmt.Print("\t")
		duration := span.finishedAt.Sub(span.startedAt)
		fmt.Printf("finish: %s (%v)\n", span.operation, duration)
		f.removeFromTrack(evt.SpanID)
	}
	//fmt.Println(f.spanChannels)
	//fmt.Println(f.openSpans)
}

func (f *Formatter) channelForSpan(spanID int) int {
	for idx, id := range f.spanChannels {
		if id == spanID {
			return idx
		}
	}
	panic(fmt.Sprintf("no channel holding span %d", spanID))
}

func getLoc(chanIdx int, toC int, fromC int) string {
	min
}

func maxLength(strs []string) int {
	maxLen := 0
	for _, s := range strs {
		if len(s) > maxLen {
			maxLen = len(s)
		}
	}
	return maxLen
}

func composite(strs []string) string {
	maxLen := maxLength(strs)
	out := []rune(strings.Repeat(" ", maxLen))
	for idx := range out {
		for _, s := range strs {
			c := s[idx]
			out[idx] = compositeChars(out[idx], c)
		}
	}
	return string(out)
}

func compositeChars(a rune, b rune) rune {
	if a == '─'
	if a == ' ' {
		return b
	}

}

func (f *Formatter) logLeftTrack(evtSpanID int, parentID int, evt string) {
	if evt == startSpanEvt {
		fromC := f.channelForSpan(parentID)
		toC := f.channelForSpan(evtSpanID)
		dir := "left"
		if toC > fromC {
			dir = "right"
		}

		for chanIdx, spanID := range f.spanChannels {
			loc := getLoc(chanIdx, toC, fromC)
		}
	}

	parentToLeft := false
	printedNode := false
	for _, spanID := range f.spanChannels {
		if spanID == -1 {
			if parentToLeft && !printedNode {
				fmt.Print("──")
			} else {
				fmt.Print("  ")
			}
			continue
		}
		if spanID == evtSpanID {
			printedNode = true
			switch evt {
			case logEvt:
				fmt.Print("*")
			case startSpanEvt:
				fmt.Print("O")
			case finishSpanEvt:
				fmt.Print("X")
			}
			fmt.Print(" ")
			continue
		}
		if spanID == parentID && evt == startSpanEvt {
			fmt.Print("├─")
			parentToLeft = true
		} else if parentToLeft && !printedNode {
			fmt.Print("┼─")
		} else {
			fmt.Print("│ ")
		}
	}
}

func (f *Formatter) addSpanToChannels(eventSpanID int) {
	for idx, spanID := range f.spanChannels {
		if spanID == -1 {
			f.spanChannels[idx] = eventSpanID
			return
		}
	}
	f.spanChannels = append(f.spanChannels, eventSpanID)
}

func (f *Formatter) removeFromTrack(evtSpanID int) {
	for idx, spanID := range f.spanChannels {
		if spanID == evtSpanID {
			f.spanChannels[idx] = -1
			return
		}
	}
}
