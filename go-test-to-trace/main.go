package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

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
	s := bufio.NewScanner(os.Stdin)

	rootSpan, ctx := tracer.StartSpan(context.Background(), "run tests!")

	c := newConverter(rootSpan, ctx)

	for s.Scan() {
		c.process(s.Text())
	}
}

type converter struct {
	// should probably just let span IDs be strings
	testNameToSpan map[string]*tracer.Span
	rootSpan       *tracer.Span
	rootCtx        context.Context
	mostRecentSpan *tracer.Span

	gettingFailureMessageFor string
}

func newConverter(rootSpan *tracer.Span, ctx context.Context) *converter {
	return &converter{
		testNameToSpan: map[string]*tracer.Span{},
		rootSpan:       rootSpan,
		rootCtx:        ctx,
		mostRecentSpan: rootSpan,
	}
}

func (c *converter) process(line string) {
	if line == "FAIL" {
		c.rootSpan.Log("FAIL")
		c.rootSpan.Finish()
		return
	}
	if line == "PASS" {
		c.rootSpan.Log("PASS")
		c.rootSpan.Finish()
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
			c.mostRecentSpan = c.rootSpan
			span.Log(line[4:])
			return
		}
		span.Log("FAIL") // TODO: finish with error or something
		span.Finish()
		c.gettingFailureMessageFor = ""
	}

	runMatch := runRegex.FindStringSubmatch(line)
	if runMatch != nil {
		span, _ := tracer.StartSpan(c.rootCtx, runMatch[1])
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
		span.Finish()
		c.mostRecentSpan = c.rootSpan
		return
	}
	contMatch := contRegex.FindStringSubmatch(line)
	if contMatch != nil {
		span, ok := c.testNameToSpan[contMatch[1]]
		if !ok {
			panic(fmt.Sprintf("couldn't find span for `%s` on line `%s`", failMatch, line))
		}
		span.Log("CONT")
		return
	}
	pauseMatch := pauseRegex.FindStringSubmatch(line)
	if pauseMatch != nil {
		span, ok := c.testNameToSpan[pauseMatch[1]]
		if !ok {
			panic(fmt.Sprintf("couldn't find span for `%s` on line `%s`", failMatch, line))
		}
		span.Log("PAUSE")
		return
	}
	logMatch := logRegex.FindStringSubmatch(line)
	if logMatch != nil {
		span, ok := c.testNameToSpan[logMatch[1]]
		if !ok {
			panic(fmt.Sprintf("couldn't find span for `%s` on line `%s`", logMatch[1], line))
		}
		span.Log(logMatch[2])
		return
	}
	//fmt.Fprintln(os.Stderr, "no match for", line)
	c.mostRecentSpan.Log(line) // TODO: ...
}
