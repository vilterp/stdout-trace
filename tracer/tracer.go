package tracer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/xid"
)

// TODO: not sure this type should event exist.
//   maybe it should hold onto a logger or something.
type Tracer struct {
	writer io.Writer
}

func NewTracer(w io.Writer) *Tracer {
	return &Tracer{
		writer: w,
	}
}

var globalTracer = NewTracer(os.Stdout)

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

func SetOutputWriter(w io.Writer) {
	globalTracer.writer = w
}

func (t *Tracer) StartSpan(ctx context.Context, operation string) (*Span, context.Context) {
	parentID, ok := ctx.Value(spanIDKey{}).(string)
	if !ok {
		parentID = ""
	}

	spanID := xid.New().String()
	ctx = context.WithValue(ctx, spanIDKey{}, spanID)
	ctx = context.WithValue(ctx, parentIDKey{}, parentID)

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
		Logs:       nil,
	}, ctx
}

type Span struct {
	ID         string     `json:"id"`
	ParentID   string     `json:"parent_id"`
	Operation  string     `json:"operation"`
	StartedAt  time.Time  `json:"started_at"`
	FinishedAt *time.Time `json:"finished_at"`
	Logs       []*LogLine `json:"logs"`
	Children   []*Span    `json:"children"`
}

type LogLine struct {
	Time time.Time `json:"time"`
	Line string    `json:"line"`
}

// TODO: remove references to globaltracer

func (s *Span) Log(line string) {
	now := time.Now()

	s.Logs = append(s.Logs, &LogLine{
		Time: now,
		Line: line,
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
	TraceEvent string    `json:"trace_evt"`
	SpanID     string    `json:"id"`
	ParentID   string    `json:"parent_id,omitempty"`
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
	fmt.Fprintln(t.writer, e.ToJSON())
}
