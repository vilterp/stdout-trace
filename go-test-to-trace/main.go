package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"

	"github.com/vilterp/stdout-trace/tracer"
)

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
}

func newConverter(rootSpan *tracer.Span, ctx context.Context) *converter {
	return &converter{
		testNameToSpan: map[string]*tracer.Span{},
		rootSpan:       rootSpan,
		rootCtx:        ctx,
		mostRecentSpan: rootSpan,
	}
}

var failRegex *regexp.Regexp
var passRegex *regexp.Regexp
var runRegex *regexp.Regexp

func init() {
	failRegex = regexp.MustCompile(`=== FAIL: (.*)`)
	passRegex = regexp.MustCompile(`--- PASS: (.*) \(.*\)`)
	runRegex = regexp.MustCompile(`=== RUN   (.*)`)
}

func (c *converter) process(line string) {
	runMatch := runRegex.FindStringSubmatch(line)
	if runMatch != nil {
		span, _ := tracer.StartSpan(c.rootCtx, runMatch[1])
		c.testNameToSpan[runMatch[1]] = span
		c.mostRecentSpan = span
		return
	}
	failMatch := failRegex.FindStringSubmatch(line)
	if failMatch != nil {
		span, ok := c.testNameToSpan[failMatch[1]]
		if !ok {
			panic(fmt.Sprintf("couldn't find span for `%s` on line `%s`", failMatch, line))
		}
		span.Log("FAIL") // TODO: finish with error or something
		span.Finish()
		c.mostRecentSpan = c.rootSpan
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
	c.mostRecentSpan.Log(line) // TODO: ...
}
