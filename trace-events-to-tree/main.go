package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/vilterp/stdout-trace/tracer"
)

func main() {
	s := bufio.NewScanner(os.Stdin)

	b := newTreeBuilder()

	for s.Scan() {
		line := s.Text()

		evt := &tracer.TraceEvent{}
		if err := json.Unmarshal([]byte(line), evt); err != nil {
			panic(fmt.Sprintf("couldn't parse trace event: %v", err))
		}
		b.process(evt)
	}

	outBytes, err := json.Marshal(b.rootSpan)
	if err != nil {
		panic(fmt.Sprintf("couldn't marshal span: %v", err))
	}

	_, err = os.Stdout.Write(outBytes)
	if err != nil {
		panic(fmt.Sprintf("couldn't write output: %v", err))
	}
}

type treeBuilder struct {
	rootSpan *tracer.Span

	openSpans map[int]*tracer.Span
}

func newTreeBuilder() *treeBuilder {
	return &treeBuilder{
		openSpans: map[int]*tracer.Span{},
	}
}

func (tb *treeBuilder) process(evt *tracer.TraceEvent) {
	switch evt.TraceEvent {
	case tracer.StartSpanEvt:
		span := &tracer.Span{
			ID:        evt.SpanID,
			ParentID:  evt.ParentID,
			Operation: evt.Operation,
			StartedAt: evt.Timestamp,
		}
		tb.openSpans[span.ID] = span
		if tb.rootSpan == nil {
			tb.rootSpan = span
		} else {
			parent := tb.openSpans[span.ParentID]
			parent.Children = append(parent.Children, span)
		}
	case tracer.LogEvt:
		span, ok := tb.openSpans[evt.SpanID]
		if !ok {
			panic(fmt.Sprintf("span not found for evt: %v", evt))
		}
		span.Logs = append(span.Logs, &tracer.LogLine{
			Time: evt.Timestamp,
			Line: evt.LogLine,
		})
	case tracer.FinishSpanEvt:
		span, ok := tb.openSpans[evt.SpanID]
		if !ok {
			panic(fmt.Sprintf("span not found for evt: %v", evt))
		}
		span.FinishedAt = &evt.Timestamp
	default:
		panic(fmt.Sprintf("unknown trace evt: %v", evt.TraceEvent))
	}
}
