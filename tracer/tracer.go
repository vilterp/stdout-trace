package tracer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type Tracer struct {
	nextID int
}

func NewTracer() *Tracer {
	return &Tracer{}
}

var globalTracer = NewTracer()

type spanIDKey struct{}
type parentIDKey struct{}

func StartSpan(ctx context.Context, operation string) (*Span, context.Context) {
	parentID, ok := ctx.Value(parentIDKey{}).(int)
	if !ok {
		parentID = -1
	}

	spanID := globalTracer.nextID
	ctx = context.WithValue(ctx, spanIDKey{}, spanID)
	ctx = context.WithValue(ctx, parentIDKey{}, parentID)

	globalTracer.nextID += 1

	now := time.Now()

	globalTracer.logEvent(&TraceEvent{
		TraceEvent: "start_span",
		Timestamp:  now,
		SpanID:     spanID,
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
		TraceEvent: "log",
		SpanID:     s.id,
		Timestamp:  now,
		LogLine:    line,
	})
}

func (s *Span) Finish() {
	now := time.Now()
	s.finishedAt = &now

	globalTracer.logEvent(&TraceEvent{
		TraceEvent: "finish_span",
		Timestamp:  now,
		SpanID:     s.id,
	})
}

// Event

type TraceEvent struct {
	TraceEvent string    `json:"evt"`
	SpanID     int       `json:"id"`
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
