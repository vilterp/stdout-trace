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
	return &Tracer{
		// start at 1 to avoid Go mistake of conflating 0 with empty... facepalm
		// TODO: maybe use separate structs for different types of events, instead of
		//   using one struct to represent their union...
		nextID: 1,
	}
}

var globalTracer = NewTracer()

type spanIDKey struct{}
type parentIDKey struct{}

const (
	StartSpanEvt  = "start_span"
	LogEvt        = "log"
	FinishSpanEvt = "finish_span"
)

func StartSpan(ctx context.Context, operation string) (*Span, context.Context) {
	return globalTracer.StartSpan(ctx, operation)
}

func (t *Tracer) StartSpan(ctx context.Context, operation string) (*Span, context.Context) {
	t.mu.Lock()
	defer t.mu.Unlock()

	parentID, ok := ctx.Value(spanIDKey{}).(int)
	if !ok {
		parentID = -1
	}

	spanID := t.nextID
	ctx = context.WithValue(ctx, spanIDKey{}, spanID)
	ctx = context.WithValue(ctx, parentIDKey{}, parentID)

	t.nextID += 1

	now := time.Now()

	t.logEvent(&TraceEvent{
		TraceEvent: StartSpanEvt,
		Timestamp:  now,
		SpanID:     spanID,
		ParentID:   parentID,
		Operation:  operation,
	})

	return &Span{
		ID:         spanID,
		ParentID:   parentID,
		Operation:  operation,
		StartedAt:  now,
		FinishedAt: nil,
		logs:       nil,
	}, ctx
}

type Span struct {
	ID         int
	ParentID   int
	Operation  string
	StartedAt  time.Time
	FinishedAt *time.Time
	logs       []*LogLine
}

type LogLine struct {
	time time.Time
	line string
}

// TODO: remove references to globaltracer

func (s *Span) Log(line string) {
	now := time.Now()

	s.logs = append(s.logs, &LogLine{
		time: now,
		line: line,
	})

	globalTracer.logEvent(&TraceEvent{
		TraceEvent: LogEvt,
		SpanID:     s.ID,
		Timestamp:  now,
		LogLine:    line,
	})
}

func (s *Span) Finish() {
	now := time.Now()
	s.FinishedAt = &now

	globalTracer.logEvent(&TraceEvent{
		TraceEvent: FinishSpanEvt,
		Timestamp:  now,
		SpanID:     s.ID,
	})
}

// Event

type TraceEvent struct {
	TraceEvent string    `json:"evt"`
	SpanID     int       `json:"id"`
	ParentID   int       `json:"parent_id,omitempty"`
	Timestamp  time.Time `json:"ts"`
	LogLine    string    `json:"line,omitempty"`
	Operation  string    `json:"op,omitempty"`
}

func (e *TraceEvent) ToJSON() string {
	bytes, err := json.Marshal(e)
	if err != nil {
		panic(err) // TODO: something else
	}
	return string(bytes)
}

func (t *Tracer) logEvent(e *TraceEvent) {
	fmt.Println(e.ToJSON())
}
