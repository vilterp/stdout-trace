package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/vilterp/stdout-trace/tracer"
)

var failRegex *regexp.Regexp
var passRegex *regexp.Regexp
var runRegex *regexp.Regexp
var pauseRegex *regexp.Regexp
var contRegex *regexp.Regexp
var logRegex *regexp.Regexp

func init() {
	runRegex = regexp.MustCompile(`=== RUN   (.*)`)
	pauseRegex = regexp.MustCompile(`=== PAUSE (.*)`)
	contRegex = regexp.MustCompile(`=== CONT  (.*)`)
	// TODO: want to avoid regex for test names if at all possible...
	logRegex = regexp.MustCompile(`=== LOG   (Test[a-zA-Z_\-0-9]+): (.*)`)
	passRegex = regexp.MustCompile(`--- PASS: (.*) \(.*\)`)
	failRegex = regexp.MustCompile(`--- FAIL: (.*) \(.*\)`)
}

func main() {
	fromTC := false
	if len(os.Args) == 2 && os.Args[1] == "--from-tc" {
		fromTC = true
	}

	c := newConverter()

	s := bufio.NewScanner(os.Stdin)
	if fromTC {
		doFromTC(s, c)
	} else {
		doNormal(s, c)
	}
}

func doNormal(s *bufio.Scanner, c *converter) {
	for s.Scan() {
		c.process(s.Text(), time.Now())
	}
}

func doFromTC(s *bufio.Scanner, c *converter) {
	for s.Scan() {
		if c.rootSpan != nil && c.rootSpan.FinishedAt != nil {
			return
		}
		line := s.Text()
		evt := &tracer.TraceEvent{}
		if err := json.Unmarshal([]byte(line), evt); err != nil {
			panic(err)
		}
		if evt.TraceEvent != tracer.LogEvt {
			continue
		}
		c.process(evt.LogLine, evt.Timestamp)
	}
}

type converter struct {
	// should probably just let span IDs be strings
	testNameToSpan map[string]*tracer.Span
	rootSpan       *tracer.Span
	mostRecentSpan *tracer.Span
	nextSpanID     int

	gettingFailureMessageFor string
}

func newConverter() *converter {
	return &converter{
		testNameToSpan: map[string]*tracer.Span{},
		nextSpanID:     2,
	}
}

func (c *converter) process(line string, ts time.Time) {
	if c.rootSpan == nil {
		rootSpan := &tracer.Span{
			ID:         "1",
			ParentID:   "",
			Operation:  "run go tests",
			StartedAt:  ts,
			FinishedAt: nil,
		}
		c.rootSpan = rootSpan
		c.mostRecentSpan = rootSpan
		c.logEvt(&tracer.TraceEvent{
			TraceEvent: tracer.StartSpanEvt,
			Operation:  "run go tests",
			SpanID:     "1",
			ParentID:   "",
			Timestamp:  ts,
		})
	}

	if line == "FAIL" {
		c.logInSpan(c.rootSpan, ts, "FAIL")
		c.finishSpan(c.rootSpan, ts)
		return
	}
	if line == "PASS" {
		c.rootSpan.Log("PASS")
		c.logInSpan(c.rootSpan, ts, "PASS")
		c.finishSpan(c.rootSpan, ts)
		return
	}

	if c.gettingFailureMessageFor != "" {
		span, ok := c.testNameToSpan[c.gettingFailureMessageFor]
		if !ok {
			panic(fmt.Sprintf("couldn't find span for `%s`", c.gettingFailureMessageFor))
		}
		if span.FinishedAt != nil {
			panic("wut")
		}
		if strings.HasPrefix(line, "    ") {
			c.logInSpan(span, ts, line[4:])
			return
		}
		c.logInSpan(span, ts, "FAIL") // TODO: finish with error or something
		c.finishSpan(span, ts)
		c.mostRecentSpan = c.rootSpan
		c.gettingFailureMessageFor = ""
		return
	}

	runMatch := runRegex.FindStringSubmatch(line)
	if runMatch != nil {
		span := c.startSpan(runMatch[1], "1", ts)
		c.testNameToSpan[runMatch[1]] = span
		c.mostRecentSpan = span
		return
	}
	failMatch := failRegex.FindStringSubmatch(line)
	if failMatch != nil {
		c.gettingFailureMessageFor = failMatch[1]
		return
	}
	passMatch := passRegex.FindStringSubmatch(line)
	if passMatch != nil {
		span, ok := c.testNameToSpan[passMatch[1]]
		if !ok {
			panic(fmt.Sprintf("couldn't find span for `%s` on line `%s`", failMatch, line))
		}
		c.finishSpan(span, ts)
		c.mostRecentSpan = c.rootSpan
		return
	}
	contMatch := contRegex.FindStringSubmatch(line)
	if contMatch != nil {
		span, ok := c.testNameToSpan[contMatch[1]]
		if !ok {
			panic(fmt.Sprintf("couldn't find span for `%s` on line `%s`", failMatch, line))
		}
		c.logInSpan(span, ts, "CONT")
		return
	}
	pauseMatch := pauseRegex.FindStringSubmatch(line)
	if pauseMatch != nil {
		span, ok := c.testNameToSpan[pauseMatch[1]]
		if !ok {
			panic(fmt.Sprintf("couldn't find span for `%s` on line `%s`", failMatch, line))
		}
		c.logInSpan(span, ts, "PAUSE")
		return
	}
	logMatch := logRegex.FindStringSubmatch(line)
	if logMatch != nil {
		span, ok := c.testNameToSpan[logMatch[1]]
		if !ok {
			panic(fmt.Sprintf("couldn't find span for `%s` on line `%s`", logMatch[1], line))
		}
		c.logInSpan(span, ts, logMatch[2])
		return
	}
	c.logInSpan(c.mostRecentSpan, ts, line)
}

func (c *converter) logEvt(evt *tracer.TraceEvent) {
	fmt.Println(evt.ToJSON())
}

func (c *converter) startSpan(opName string, parentID string, ts time.Time) *tracer.Span {
	rawID := c.nextSpanID
	c.nextSpanID++
	id := fmt.Sprintf("%d", rawID)

	c.logEvt(&tracer.TraceEvent{
		SpanID:     id,
		ParentID:   parentID,
		TraceEvent: tracer.StartSpanEvt,
		Timestamp:  ts,
		Operation:  opName,
	})

	return &tracer.Span{
		ID:        id,
		ParentID:  parentID,
		StartedAt: ts,
		Operation: opName,
	}
}

func (c *converter) logInSpan(s *tracer.Span, ts time.Time, line string) {
	c.logEvt(&tracer.TraceEvent{
		TraceEvent: tracer.LogEvt,
		SpanID:     s.ID,
		Timestamp:  ts,
		LogLine:    line,
	})
}

func (c *converter) finishSpan(s *tracer.Span, ts time.Time) {
	s.FinishedAt = &ts
	c.logEvt(&tracer.TraceEvent{
		TraceEvent: tracer.FinishSpanEvt,
		SpanID:     s.ID,
		Timestamp:  ts,
	})
}
