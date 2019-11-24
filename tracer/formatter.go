package tracer

import "fmt"

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
		f.logLeftTrack(evt.SpanID)
		fmt.Print("\t")
		fmt.Printf("start %d: %s\n", evt.SpanID, evt.Operation)
	case logEvt:
		f.logLeftTrack(evt.SpanID)
		fmt.Print("\t")
		fmt.Println(evt.LogLine)
	case finishSpanEvt:
		span := f.openSpans[evt.SpanID]
		span.finishedAt = &evt.Timestamp
		f.logLeftTrack(evt.SpanID)
		fmt.Print("\t")
		duration := span.finishedAt.Sub(span.startedAt)
		fmt.Printf("finish %d: %s (%v)\n", evt.SpanID, span.operation, duration)
		f.removeFromTrack(evt.SpanID)
	}
	//fmt.Println(f.spanChannels)
	//fmt.Println(f.openSpans)
}

func (f *Formatter) logLeftTrack(evtSpanID int) {
	for _, spanID := range f.spanChannels {
		if spanID == -1 {
			fmt.Print("  ")
			continue
		}
		if spanID == evtSpanID {
			fmt.Print("X ")
			continue
		}
		fmt.Print("| ")
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
