package main

import (
	"context"

	"github.com/pmorelli92/demo-rabbitmq-streams/order/app"
)

func main() {
	app.Run(context.Background())
}
