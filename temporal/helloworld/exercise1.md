# Exercise 1

Let's write a simple Helloworld workflow to say hello to the Trainer and the Pokemon.

Our workflow will be quite simple, with two activities:
- Say hello to trainer
- Say hello to Pokémon

The output of the workflow will be a simple concatenation of the two activity's results with a space.

For example, with the Trainer Sacha and the Pokémon Pikachu:
- Activity 1 returns `Hello Sacha!`
- Activity 2 returns `Hello Pikachu!`
- Workflow return `Hello Sacha! Hello Pikachu!`

## SayHello activities

Start by writing two activities.

An activity has special rules for its function signature:
- The activity must accept a `context.Context` as first input, then as many parameters as needed
- The activity may return an output, and must then return an error

Create the `activity.go` and `activity_test.go` file.

Write the two functions `SayHelloToTrainer` and `SayHelloToPokemon`.

Test your functions with this code:

SayHelloToTrainer
```go
func TestSayHelloToTrainer(t *testing.T) {
	// Arrange
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestActivityEnvironment()
	env.RegisterActivity(SayHelloToTrainer)
	blue := mocks.Blue()

	// Act
	val, err := env.ExecuteActivity(SayHelloToTrainer, blue)

	// Assert
	require.NoError(t, err)
	var res string
	require.NoError(t, val.Get(&res))
	require.Equal(t, "Hello Blue!", res)
}
```

SayHelloToPokemon
```go
func TestSayHelloToPokemon(t *testing.T) {
	// Arrange
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestActivityEnvironment()
	env.RegisterActivity(SayHelloToPokemon)
	rattata := mocks.Rattata()

	// Act
	val, err := env.ExecuteActivity(SayHelloToPokemon, rattata)

	// Assert
	require.NoError(t, err)
	var res string
	require.NoError(t, val.Get(&res))
	require.Equal(t, "Hello Rattata!", res)
}
```

## Helloworld workflow

Let's now write our workflow Helloworld.
Create `workflow.go` and a `workflow_test.go` files.

A Temporal workflow has special rules for its function signature:
- The workflow must accept a `workflow.Context` as first input, then as many parameters as needed
- The activity may return an output, and must then return an error

*Notice: it is recommended to accept only one input parameter after the context.*

Calling an activity from a workflow is done using the Temporal SDK's function `workflow.ExecuteActivity` (see [https://pkg.go.dev/go.temporal.io/sdk/workflow#ExecuteActivity](https://pkg.go.dev/go.temporal.io/sdk/workflow#ExecuteActivity)).

Here's a bit of help to write your workflow, where the `SayHelloToTrainer` activity is already called:

```go
func Helloworld(ctx workflow.Context, trainer *pokemon.Trainer, pokemon *pokemon.Pokemon) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var result string
	err := workflow.ExecuteActivity(ctx, SayHelloToTrainer, trainer).Get(ctx, &result)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return "", err
	}

	return result, nil
}
```

Take a bit of time to understand the code. Then modify it to add the call to `SayHelloToPokemon` and return the expected result.

You can use this code to test your work:

```go
func TestWorkflow(t *testing.T) {
	// Arrange
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()
	blue := mocks.Blue()
	rattata := mocks.Rattata()

	env.OnActivity(SayHelloToTrainer, mock.Anything, blue).Return("Hello Blue!", nil)
	env.OnActivity(SayHelloToPokemon, mock.Anything, rattata).Return("Hello Rattata!", nil)

	// Act
	env.ExecuteWorkflow(Helloworld, blue, rattata)

	// Assert
	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
	var result string
	require.NoError(t, env.GetWorkflowResult(&result))
	require.Equal(t, "Hello Blue! Hello Rattata!", result)
}
```

## Worker

Your code is now ready to be executed by Temporal!

Let's now create a worker to execute the code.

Add a new file `main.go` in the [worker folder](./worker/).

Copy paste this function:

```go
package main

import (
	"log/slog"

	"github.com/nathancastelein/go-course-workflows/temporal/helloworld"
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
```

Take time to read the content.

Then launch your worker:

```bash
$ go run main.go
2024/08/11 21:36:42 INFO Started Worker Namespace=default TaskQueue=helloworld WorkerID=62239@mbp-de-nathan.home@
```

## Starter

The starter is in charge of asking Temporal to launch a new execution of a workflow.

Add a new file `main.go` in the [starter folder](./starter/).

Add this content:

```go
package main

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/nathancastelein/go-course-workflows/pokemon"
	"github.com/nathancastelein/go-course-workflows/solution/temporal/helloworld"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/log"
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

	workflowOptions := client.StartWorkflowOptions{
		ID:        uuid.New().String(),
		TaskQueue: "helloworld",
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, helloworld.Helloworld, pokemon.Sacha(), pokemon.Pikachu())
	if err != nil {
		slog.Error("unable to execute workflow", slog.Any("error", err))
		return
	}

	slog.Info("started workflow", slog.String("workflow_id", we.GetID()), slog.String("run_id", we.GetRunID()))

	var result string
	err = we.Get(context.Background(), &result)
	if err != nil {
		slog.Error("unable to get workflow result", slog.Any("error", err))
		return
	}
	slog.Info("workflow result", slog.Any("result", result))
}
```

Take time to read the function.

Then run your starter!

```bash
$ go run main.go 
2024/08/11 21:36:59 INFO started workflow workflow_id=ea0d2634-2429-4c94-9202-60032d3606b8 run_id=13efe5d5-b8b8-448e-9114-28f05d789260
2024/08/11 21:36:59 INFO workflow result result="Hello Sacha! Hello Pikachu!"
```

Check the workflow execution [on the UI](http://localhost:8233/namespaces/default/workflows).