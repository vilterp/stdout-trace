package basic_goroutines

import (
	"context"
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
}

func a(ctx context.Context) {
	span, _ := tracer.StartSpan(ctx, "a")
	defer span.Finish()

	span.Log("AAA begin")
	time.Sleep(5 * time.Second)
	span.Log("AAA sleep moar")
	time.Sleep(5 * time.Second)
	span.Log("AAA done")
}

func b(ctx context.Context) {
	span, _ := tracer.StartSpan(ctx, "b")
	defer span.Finish()

	span.Log("BBB sup")
	time.Sleep(1 * time.Second)
	span.Log("BBB yo")
}

func c(ctx context.Context) {
	span, _ := tracer.StartSpan(ctx, "c")
	defer span.Finish()

	time.Sleep(1 * time.Second)
	span.Log("CCC blurp")
	time.Sleep(2 * time.Second)
	span.Log("CCC durp")
}
