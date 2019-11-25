package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/vilterp/stdout-trace/tracer"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: trace-replay <file.ldjson>")
		os.Exit(1)
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	s := bufio.NewScanner(f)

	var lastTS *time.Time

	for s.Scan() {
		line := s.Text()

		evt := &tracer.TraceEvent{}
		if err := json.Unmarshal([]byte(line), evt); err != nil {
			panic(err)
		}

		if lastTS == nil {
			lastTS = &evt.Timestamp
			fmt.Println(line)
			continue
		}

		time.Sleep(evt.Timestamp.Sub(*lastTS))
		lastTS = &evt.Timestamp
		fmt.Println(line)
	}
}
