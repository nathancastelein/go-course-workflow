package main

import (
	"log/slog"

	"github.com/nathancastelein/go-course-workflows/solution/temporal/helloworld"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.Dial(client.Options{
		HostPort: "localhost:7233",
		Logger:   log.NewStructuredLogger(slog.Default()),
	})
	if err != nil {
		slog.Error("unable to create client", slog.Any("error", err))
		return
	}
	defer c.Close()

	w := worker.New(c, "helloworld", worker.Options{})

	w.RegisterWorkflow(helloworld.Helloworld)
	w.RegisterActivity(helloworld.SayHelloToTrainer)
	w.RegisterActivity(helloworld.SayHelloToPokemon)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		slog.Error("unable to start worker", slog.Any("error", err))
		return
	}
}
