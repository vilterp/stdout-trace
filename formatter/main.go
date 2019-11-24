package main

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/vilterp/stdout-trace/tracer"
)

func main() {
	f := tracer.NewFormatter()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		evt := &tracer.TraceEvent{}
		if err := json.Unmarshal([]byte(text), evt); err != nil {
			panic(err)
		}
		f.Handle(evt)
	}
}
