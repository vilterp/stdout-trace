package main

import (
	"context"
	"time"

	"github.com/vilterp/stdout-trace/tracer"
)

func main() {
	ctx := context.Background()

	span, ctx := tracer.StartSpan(ctx, "main")
	defer span.Finish()

	a(ctx)
	go b(ctx)
	time.Sleep(2 * time.Second)
	go c(ctx)
}

func a(ctx context.Context) {
	span, _ := tracer.StartSpan(ctx, "a")
	defer span.Finish()

	span.Log("begin")
	time.Sleep(5 * time.Second)
	span.Log("sleep moar")
	time.Sleep(5 * time.Second)
	span.Log("done")
}

func b(ctx context.Context) {
	span, _ := tracer.StartSpan(ctx, "b")
	defer span.Finish()

	span.Log("sup")
	time.Sleep(1 * time.Second)
	span.Log("yo")
}

func c(ctx context.Context) {
	span, _ := tracer.StartSpan(ctx, "c")
	defer span.Finish()

	time.Sleep(1 * time.Second)
	span.Log("blurp")
	time.Sleep(2 * time.Second)
	span.Log("durp")
}
