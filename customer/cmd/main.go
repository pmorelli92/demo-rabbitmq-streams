package main

import (
	"context"

	"github.com/pmorelli92/demo-rabbitmq-streams/app"
)

func main() {
	app.Run(context.Background())
}
