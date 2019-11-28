package format

import (
	"fmt"
	"strings"

	"github.com/vilterp/stdout-trace/tracer"
)

type Formatter struct {
	spanChannels []int
	openSpans    map[int]*tracer.Span
}

func NewFormatter() *Formatter {
	return &Formatter{
		openSpans: map[int]*tracer.Span{},
	}
}

func (f *Formatter) Handle(evt *tracer.TraceEvent) {
	switch evt.TraceEvent {
	case tracer.StartSpanEvt:
		f.openSpans[evt.SpanID] = &tracer.Span{
			StartedAt: evt.Timestamp,
			Operation: evt.Operation,
			ID:        evt.SpanID,
			ParentID:  evt.ParentID,
		}
		f.addSpanToChannels(evt.SpanID)
		f.logLeftTrack(evt.SpanID, evt.ParentID, tracer.StartSpanEvt)
		fmt.Print("\t")
		fmt.Printf("start: %s\n", evt.Operation)
	case tracer.LogEvt:
		f.guardOpenSpans(evt)
		f.logLeftTrack(evt.SpanID, -1, tracer.LogEvt)
		fmt.Print("\t")
		fmt.Println(evt.LogLine)
	case tracer.FinishSpanEvt:
		f.guardOpenSpans(evt)
		span := f.openSpans[evt.SpanID]
		span.FinishedAt = &evt.Timestamp
		f.logLeftTrack(evt.SpanID, -1, tracer.FinishSpanEvt)
		fmt.Print("\t")
		duration := span.FinishedAt.Sub(span.StartedAt)
		fmt.Printf("finish: %s (%v)\n", span.Operation, duration)
		f.removeFromTrack(evt.SpanID)
	}
}

func (f *Formatter) channelForSpan(spanID int) (int, bool) {
	for idx, id := range f.spanChannels {
		if id == spanID {
			return idx, true
		}
	}
	return 0, false
}

func (f *Formatter) logLeftTrack(evtSpanID int, parentID int, evt string) {
	fmt.Print(f.getLeftTrack(evtSpanID, parentID, evt).spaceOut().String())
}

func (f *Formatter) getLeftTrack(evtSpanID int, parentID int, evt string) Line {
	switch evt {
	case tracer.StartSpanEvt:
		return compositeLines([]Line{
			f.existingChannelsLine(evtSpanID),
			f.spawnLine(parentID, evtSpanID),
		})
	case tracer.LogEvt:
		fallthrough
	case tracer.FinishSpanEvt:
		return compositeLines([]Line{
			f.existingChannelsLine(-1), // ugh sentinels
			f.evtLine(evtSpanID, evt),
		})
	default:
		panic(fmt.Sprintf("unrecognized event %s", evt))
	}
}

func (f *Formatter) spawnLine(fromSpanID int, toSpanID int) Line {
	if fromSpanID == -1 { // root span
		return nil
	}
	fromC, ok := f.channelForSpan(fromSpanID)
	if !ok {
		return nil
	}
	toC, ok := f.channelForSpan(toSpanID)
	if !ok {
		return nil
	}
	left := min(fromC, toC)
	right := max(fromC, toC)
	return horizLine(left, right)
}

func (f *Formatter) existingChannelsLine(newSpanID int) Line {
	var out Line
	for _, spanID := range f.spanChannels {
		if spanID == -1 {
			out = append(out, Empty)
		} else if spanID == newSpanID {
			out = append(out, VertDownHalf)
		} else {
			out = append(out, VertLine)
		}
	}
	return out
}

func (f *Formatter) evtLine(spanID int, evt string) Line {
	c, ok := f.channelForSpan(spanID)
	if !ok {
		panic(fmt.Sprintf("can't find channel for event in span %d. channels: %v. evt: %v", spanID, f.spanChannels, evt))
	}
	out := Line(strings.Repeat(" ", c))
	switch evt {
	case tracer.LogEvt:
		return append(out, Log)
	case tracer.FinishSpanEvt:
		return append(out, Finish)
	default:
		return Line("")
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

func (f *Formatter) guardOpenSpans(evt *tracer.TraceEvent) {
	if _, ok := f.openSpans[evt.SpanID]; !ok {
		panic(fmt.Sprintf("span id %d already closed or never opened. evt: %+v", evt.SpanID, evt))
	}
}
