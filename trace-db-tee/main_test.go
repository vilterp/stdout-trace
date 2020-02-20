package main

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vilterp/stdout-trace/tracer"
)

func TestInsertEvent(t *testing.T) {
	db, err := sql.Open("postgres", "postgres://root@localhost:26257/traces?sslmode=disable")
	require.NoError(t, err)

	p := newProcessor(db)
	require.NoError(t, p.updateDB(&tracer.TraceEvent{
		TraceEvent: tracer.StartSpanEvt,
		SpanID:     "abc-123",
		ParentID:   "root",
		Timestamp:  time.Now(),
		Operation:  "foo",
	}))
}
