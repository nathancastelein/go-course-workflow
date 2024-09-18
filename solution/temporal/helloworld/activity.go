package helloworld

import (
	"context"
	"fmt"

	"github.com/nathancastelein/go-course-workflows/pokemon"
	"go.temporal.io/sdk/activity"
)

func SayHelloToTrainer(ctx context.Context, trainer *pokemon.Trainer) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("SayHelloToTrainer", "name", trainer.Name)
	return fmt.Sprintf("Hello %s!", trainer.Name), nil
}

func SayHelloToPokemon(ctx context.Context, pokemon *pokemon.Pokemon) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("SayHelloToPokemon", "name", pokemon.Name)
	return fmt.Sprintf("Hello %s!", pokemon.Name), nil
}
