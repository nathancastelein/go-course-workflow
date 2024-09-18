# Exercise 2

Now it's time to write your workflow, and test it with Temporal!

## Workflow

Writing your workflow must follow the Temporal's guidelines, as we did with the Helloworld workflow.

Create new `workflow.go` and `workflow_test.go` files.
 
Here's a skeleton for your workflow. For the same reason than the `ThrowPokeballActivity`, we created a new type for the workflow output.

```go
type CapturePokemonOutput struct {
	Trainer *pokemon.Trainer
	Pokemon *pokemon.Pokemon
}

func (w *Worker) CapturePokemonWorkflow(ctx workflow.Context, trainer *pokemon.Trainer, pokemon *pokemon.Pokemon) (*CapturePokemonOutput, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("CapturePokemon workflow started")

	// Start your activities sequentially (Paralyze, Attack, ThrowPokeball)

	logger.Info("CapturePokemon workflow completed.")

	return &CapturePokemonOutput{}, nil
}
```

To test your code:

```go
func TestWorkflow(t *testing.T) {
	// Arrange
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()
	worker := &Worker{}
	rattata := mocks.Rattata()
	blue := mocks.Blue()

	env.OnActivity(worker.ParalyzeActivity, mock.Anything, rattata).Return(rattata, nil)
	env.OnActivity(worker.AttackActivity, mock.Anything, rattata).Return(rattata, nil)
	env.OnActivity(worker.ThrowPokeballActivity, mock.Anything, blue, rattata).Return(ThrowPokeballOutput{
		Trainer: blue,
		Pokemon: rattata,
	}, nil)

	// Act
	env.ExecuteWorkflow(worker.CapturePokemonWorkflow, blue, rattata)

	// Assert
	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
	var result CapturePokemonOutput
	require.NoError(t, env.GetWorkflowResult(&result))
	require.Equal(t, rattata, result.Pokemon)
	require.Equal(t, blue, result.Trainer)
}
```

## Worker

And now it's time to create the last two files: the worker and the starter!

Here's the worker code you can add in [cmd folder](./cmd/worker/)

```go
package main

import (
	"log/slog"

	"github.com/nathancastelein/go-course-workflows/pokemon"
	"github.com/nathancastelein/go-course-workflows/temporal/capture"
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
```

Then run your worker:

```bash
$ go run main.go
```

## Starter

Here's the starter code you can add in [cmd folder](./cmd/starter/)

```go
package main

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/nathancastelein/go-course-workflows/pokemon"
	"github.com/nathancastelein/go-course-workflows/solution/capture"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/log"
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

	workflowOptions := client.StartWorkflowOptions{
		ID:        uuid.New().String(),
		TaskQueue: "capture-pokemon",
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, "CapturePokemonWorkflow", pokemon.Sacha(), pokemon.Mewtwo())
	if err != nil {
		slog.Error("unable to execute workflow", slog.Any("error", err))
		return
	}

	slog.Info("started workflow", slog.String("workflow_id", we.GetID()), slog.String("run_id", we.GetRunID()))

	var result capture.CapturePokemonOutput
	err = we.Get(context.Background(), &result)
	if err != nil {
		slog.Error("unable to get workflow result", slog.Any("error", err))
		return
	}
	slog.Info("workflow result", slog.Any("trainer", result.Trainer), slog.Any("pokemon", result.Pokemon))
}
```

Then run your starter:

```bash
$ go run main.go
```

Check the workflow execution [on the UI](http://localhost:8233/namespaces/default/workflows)!