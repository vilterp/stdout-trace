package main

import (
	"bufio"
	"encoding/json"
	"os"
	"time"

	"github.com/vilterp/stdout-trace/trace-format/format"
	"github.com/vilterp/stdout-trace/tracer"
)

func main() {
	f := format.NewFormatter()

	line := 1

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		evt := &tracer.TraceEvent{}
		err := json.Unmarshal([]byte(text), evt)
		if err != nil {
			//log.Printf("parse error line %d: %v", line, err)
			if f.FirstSpanID == "" {
				continue
			}
			evt = &tracer.TraceEvent{
				TraceEvent: tracer.LogEvt,
				SpanID:     f.FirstSpanID,
				Timestamp:  time.Now(), // TODO: doesn't work when replaying an old one...
				LogLine:    text,
			}
		}
		f.Handle(evt)
		line++
	}
}
