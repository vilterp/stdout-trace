package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/vilterp/stdout-trace/tracer"
)

var timestampRegex = regexp.MustCompile(`\[(\d\d:\d\d:\d\d)\].:\s+\[Step \d/\d\] (.*)`)

func main() {
	p := newProcessor()

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		line := s.Text()

		p.process(line)
	}
	p.finish()
}

type processor struct {
	firstLineProcessed bool
	lastTS             time.Time
}

func newProcessor() *processor {
	return &processor{}
}

func (p *processor) process(line string) {
	tsMatch := timestampRegex.FindStringSubmatch(line)
	if tsMatch == nil {
		return
	}

	rawTS := tsMatch[1]
	lineContent := tsMatch[2]

	ts, err := time.Parse("15:04:05", rawTS)
	if err != nil {
		panic(err)
	}
	p.lastTS = ts

	if !p.firstLineProcessed {
		p.firstLineProcessed = true

		evt := &tracer.TraceEvent{
			TraceEvent: tracer.StartSpanEvt,
			SpanID:     1,
			ParentID:   -1,
			Timestamp:  ts,
			Operation:  "Run TeamCity build",
			LogLine:    lineContent,
		}
		p.logEvt(evt)
	}

	p.logEvt(&tracer.TraceEvent{
		TraceEvent: tracer.LogEvt,
		SpanID:     1,
		Timestamp:  ts,
		LogLine:    lineContent,
	})
}

func (p *processor) finish() {
	p.logEvt(&tracer.TraceEvent{
		TraceEvent: tracer.FinishSpanEvt,
		SpanID:     1,
		Timestamp:  p.lastTS,
	})
}

func (p *processor) logEvt(evt *tracer.TraceEvent) {
	fmt.Println(evt.ToJSON())
}
