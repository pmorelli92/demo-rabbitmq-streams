package accept_test

import (
	"context"
	"testing"
	"time"

	"github.com/pmorelli92/demo-rabbitmq-streams/customer/app"
	"go.uber.org/goleak"
)

func TestNoLeak(t *testing.T) {
	ctx, cancelCtx := context.WithCancel(context.TODO())
	go app.Run(ctx)
	time.Sleep(2 * time.Second)
	cancelCtx()
	time.Sleep(2 * time.Second)
	goleak.VerifyNone(t)
}
