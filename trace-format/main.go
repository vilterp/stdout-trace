package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

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
			log.Fatalf("parse error line %d: %v", line, err)
		}
		f.Handle(evt)
		line++
	}
	fmt.Println("really done?", scanner.Scan())
	fmt.Println("last text", scanner.Text())
}
