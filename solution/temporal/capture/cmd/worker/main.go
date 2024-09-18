package main

import (
	"log/slog"

	"github.com/nathancastelein/go-course-workflows/pokemon"
	"github.com/nathancastelein/go-course-workflows/solution/temporal/capture"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.Dial(client.Options{
		Logger: log.NewStructuredLogger(slog.Default()),
	})
	if err != nil {
		slog.Error("unable to create client", slog.Any("error", err))
		return
	}
	defer c.Close()

	w := worker.New(c, "capture-pokemon", worker.Options{})

	pokemonWorker := capture.NewWorker(
		pokemon.NewStatusService(),
		pokemon.NewCombatService(),
		pokemon.NewPokeballService(),
	)

	w.RegisterWorkflow(pokemonWorker.CapturePokemonWorkflow)
	w.RegisterActivity(pokemonWorker.ParalyzeActivity)
	w.RegisterActivity(pokemonWorker.AttackActivity)
	w.RegisterActivity(pokemonWorker.ThrowPokeballActivity)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		slog.Error("unable to start worker", slog.Any("error", err))
		return
	}
}
