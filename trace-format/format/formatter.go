package format

import (
	"fmt"

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
		f.logLeftTrack(evt.SpanID, -1, tracer.LogEvt)
		fmt.Print("\t")
		fmt.Println(evt.LogLine)
	case tracer.FinishSpanEvt:
		span := f.openSpans[evt.SpanID]
		span.FinishedAt = &evt.Timestamp
		f.logLeftTrack(evt.SpanID, -1, tracer.FinishSpanEvt)
		fmt.Print("\t")
		duration := span.FinishedAt.Sub(span.StartedAt)
		fmt.Printf("finish: %s (%v)\n", span.Operation, duration)
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

func (f *Formatter) logLeftTrack(evtSpanID int, parentID int, evt string) {
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
