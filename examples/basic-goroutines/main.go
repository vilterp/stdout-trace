package main

import (
	"context"
	"fmt"
	"time"

	"github.com/vilterp/stdout-trace/tracer"
)

func main() {
	ctx := context.Background()

	span, ctx := tracer.StartSpan(ctx, "main")
	defer span.Finish()

	go b(ctx)
	time.Sleep(2 * time.Second)
	go c(ctx)
	a(ctx)
	span.Log("sleeping")
	time.Sleep(5 * time.Second)
}

func a(ctx context.Context) {
	span, ctx := tracer.StartSpan(ctx, "a")
	defer span.Finish()

	span.Log("AAA begin")
	time.Sleep(5 * time.Second)
	span.Log("AAA sleep moar")
	time.Sleep(5 * time.Second)
	e(ctx)
	span.Log("AAA done")
}

func b(ctx context.Context) {
	span, ctx := tracer.StartSpan(ctx, "b")
	defer span.Finish()

	f(ctx)
	span.Log("BBB sup")
	time.Sleep(1 * time.Second)
	span.Log("BBB yo")
}

func c(ctx context.Context) {
	span, ctx := tracer.StartSpan(ctx, "c")
	defer span.Finish()

	time.Sleep(1 * time.Second)
	d(ctx)
	span.Log("CCC blurp")
	time.Sleep(2 * time.Second)
	span.Log("CCC durp")
}

func d(ctx context.Context) {
	span, ctx := tracer.StartSpan(ctx, "d")
	defer span.Finish()

	time.Sleep(1 * time.Second)
	span.Log("DDD blurp")
	time.Sleep(2 * time.Second)
	span.Log("DDD durp")
}

func e(ctx context.Context) {
	span, ctx := tracer.StartSpan(ctx, "e")
	defer span.Finish()

	time.Sleep(1 * time.Second)
	span.Log("EEE blurp")
	time.Sleep(1 * time.Second)
	span.Log("EEE durp")
}

func f(ctx context.Context) {
	span, ctx := tracer.StartSpan(ctx, "f")
	defer span.Finish()

	time.Sleep(1 * time.Second)
	span.Log("FFF blurp")
	go g(ctx)
	time.Sleep(2 * time.Second)
	span.Log("FFF durp")
}

func g(ctx context.Context) {
	span, ctx := tracer.StartSpan(ctx, "g")
	defer span.Finish()

	for i := 0; i < 3; i++ {
		time.Sleep(500 * time.Millisecond)
		go h(ctx)
	}
}

func h(ctx context.Context) {
	span, ctx := tracer.StartSpan(ctx, "h")
	defer span.Finish()

	go a(ctx)

	for i := 0; i < 2; i++ {
		span.Log(fmt.Sprintf("h %d", i))

		go d(ctx)

		time.Sleep(200 * time.Millisecond)
	}
}
