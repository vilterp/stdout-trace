package tracer

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type Tracer struct {
	mu     sync.Mutex
	nextID int
}

func NewTracer() *Tracer {
	return &Tracer{}
}

var globalTracer = NewTracer()

type spanIDKey struct{}
type parentIDKey struct{}

const (
	startSpanEvt  = "start_span"
	logEvt        = "log"
	finishSpanEvt = "finish_span"
)

func StartSpan(ctx context.Context, operation string) (*Span, context.Context) {
	return globalTracer.StartSpan(ctx, operation)
}

func (t *Tracer) StartSpan(ctx context.Context, operation string) (*Span, context.Context) {
	globalTracer.mu.Lock()
	defer globalTracer.mu.Unlock()

	parentID, ok := ctx.Value(spanIDKey{}).(int)
	if !ok {
		parentID = -1
	}

	spanID := globalTracer.nextID
	ctx = context.WithValue(ctx, spanIDKey{}, spanID)
	ctx = context.WithValue(ctx, parentIDKey{}, parentID)

	globalTracer.nextID += 1

	now := time.Now()

	globalTracer.logEvent(&TraceEvent{
		TraceEvent: startSpanEvt,
		Timestamp:  now,
		SpanID:     spanID,
		ParentID:   parentID,
		Operation:  operation,
	})

	return &Span{
		id:         spanID,
		parentID:   parentID,
		operation:  operation,
		startedAt:  now,
		finishedAt: nil,
		logs:       nil,
	}, ctx
}

type Span struct {
	id         int
	parentID   int
	operation  string
	startedAt  time.Time
	finishedAt *time.Time
	logs       []*LogLine
}

type LogLine struct {
	time time.Time
	line string
}

func (s *Span) Log(line string) {
	now := time.Now()

	s.logs = append(s.logs, &LogLine{
		time: now,
		line: line,
	})

	globalTracer.logEvent(&TraceEvent{
		TraceEvent: logEvt,
		SpanID:     s.id,
		Timestamp:  now,
		LogLine:    line,
	})
}

func (s *Span) Finish() {
	now := time.Now()
	s.finishedAt = &now

	globalTracer.logEvent(&TraceEvent{
		TraceEvent: finishSpanEvt,
		Timestamp:  now,
		SpanID:     s.id,
	})
}

// Event

type TraceEvent struct {
	TraceEvent string    `json:"trace_evt"`
	SpanID     int       `json:"id"`
	ParentID   int       `json:"parent_id,omitempty"`
	Timestamp  time.Time `json:"ts"`
	LogLine    string    `json:"line,omitempty"`
	Operation  string    `json:"op,omitempty"`
}

func (t *Tracer) logEvent(e *TraceEvent) {
	bytes, err := json.Marshal(e)
	if err != nil {
		panic(err) // TODO: something else
	}
	fmt.Println(string(bytes))
}
