package helloworld

import (
	"time"

	"github.com/nathancastelein/go-course-workflows/pokemon"
	"go.temporal.io/sdk/workflow"
)

func Helloworld(ctx workflow.Context, trainer *pokemon.Trainer, pokemon *pokemon.Pokemon) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)

	var finalResult string

	var result string
	err := workflow.ExecuteActivity(ctx, SayHelloToTrainer, trainer).Get(ctx, &result)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return "", err
	}
	finalResult += result

	err = workflow.ExecuteActivity(ctx, SayHelloToPokemon, pokemon).Get(ctx, &result)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return "", err
	}
	finalResult += " " + result

	err = workflow.ExecuteActivity(ctx, SayHelloToProfessorOak).Get(ctx, &result)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return "", err
	}
	finalResult += " " + result

	logger.Info("HelloWorld workflow completed.")

	return finalResult, nil
}
