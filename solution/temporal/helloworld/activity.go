package helloworld

import (
	"context"
	"fmt"
	"io"
	"net/http"

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

func SayHelloToProfessorOak(ctx context.Context) (string, error) {
	resp, err := http.Get("localhost:8080/hello")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(result), nil
}
