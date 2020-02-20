package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/vilterp/stdout-trace/tracer"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage: trace-db-tee <db conn string>")
	}

	dbConnString := os.Args[1]

	db, err := sql.Open("postgres", dbConnString)
	if err != nil {
		log.Fatal(err)
	}

	p := newProcessor(db)
	if err := p.ensureSchema(); err != nil {
		log.Fatal("error ensuring schema: ", err)
	}

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		line := s.Bytes()

		fmt.Println(string(line))
		if err := p.process(line); err != nil {
			fmt.Fprintln(os.Stderr, "error processing event:", err)
		}
	}
}

type processor struct {
	conn *sql.DB
}

func newProcessor(db *sql.DB) *processor {
	return &processor{conn: db}
}

func (p *processor) process(b []byte) error {
	evt := &tracer.TraceEvent{}
	if err := json.Unmarshal(b, evt); err != nil {
		return errors.Wrap(err, "unmarshalling event")
	}
	return p.updateDB(evt)
}

func (p *processor) updateDB(evt *tracer.TraceEvent) error {
	switch evt.TraceEvent {
	case tracer.StartSpanEvt:
		_, err := p.conn.Exec(
			"INSERT INTO spans (id, parent_id, started_at) VALUES ($1, $2, $3)",
			evt.SpanID, evt.ParentID, evt.Timestamp,
		)
		return err
	case tracer.LogEvt:
		_, err := p.conn.Exec(
			"INSERT INTO logs (span_id, timestamp, text) VALUES ($1, $2, $3)",
			evt.SpanID, evt.Timestamp, evt.LogLine,
		)
		return err
	case tracer.FinishSpanEvt:
		_, err := p.conn.Exec("UPDATE spans SET finished_at = $1 WHERE id = $2", evt.Timestamp, evt.SpanID)
		return err
	default:
		return nil
	}
}

var schema = `
CREATE TABLE IF NOT EXISTS spans (
	id TEXT PRIMARY KEY,
	parent_id TEXT, -- TODO: foreign key?
	started_at TIMESTAMPTZ,
	finished_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS logs (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	span_id TEXT,
	timestamp TIMESTAMPTZ,
	text TEXT
);
`

func (p *processor) ensureSchema() error {
	_, err := p.conn.Exec(schema)
	return err
}
