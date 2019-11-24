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

func (f *Formatter) logLeftTrack(evtSpanID int, parentID int, evt string) {
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
