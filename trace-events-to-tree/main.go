package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"

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

	openSpans map[string]*tracer.Span
}

func newTreeBuilder() *treeBuilder {
	root := &tracer.Span{
		ID:         "__ROOT__",
		ParentID:   "",
		Operation:  "root",
		StartedAt:  time.Time{},
		FinishedAt: nil,
		Logs:       nil,
		Children:   nil,
	}
	return &treeBuilder{
		rootSpan: root,
		openSpans: map[string]*tracer.Span{
			root.ID: root,
		},
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
		parent := tb.openSpans[span.ParentID]
		if parent != nil {
			parent.Children = append(parent.Children, span)
		} else {
			tb.rootSpan.Children = append(tb.rootSpan.Children, span)
		}
	case tracer.LogEvt:
		span, ok := tb.openSpans[evt.SpanID]
		if !ok {
			panic(fmt.Sprintf("span not found for evt: %v", evt))
		}
		span.Logs = append(span.Logs, &tracer.LogLine{
			Time: evt.Timestamp,
			Line: evt.LogLine,

			Tags:  evt.Tags,
			Attrs: evt.Attrs,
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
