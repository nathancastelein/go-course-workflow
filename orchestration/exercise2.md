# Exercise 2

Now that we have a working orchestrator, it's time to write your first workflow: the `Helloworld` workflow.

The workflow is a simple workflow with two steps:
- The first step which says Hello to the trainer
- The second step which says Hello to the Pokémon

*Notice: This exercise will mix two concepts: the `worker` and the `workflow definition`. Keep in mind that those two concepts can be splitted if needed.*

## Worker

Let's first create a new struct `HelloworldWorker` in a `helloworld.go` file. It`s an empty struct, to hold data.

As previously mentioned, the workflow has two steps:
- hello trainer
- hello Pokémon

Each step will be a method attached to the worker.

As the method is a workflow `Step`, it must follow the definition of the `Do` function of a step!

Add two methods to your struct:

```go
func (c *HelloworldWorker) helloTrainer(trainer *pokemon.Trainer, pokemon *pokemon.Pokemon) error
```

```go
func (c *HelloworldWorker) helloPokemon(trainer *pokemon.Trainer, pokemon *pokemon.Pokemon) error
```

Each method must print `Hello ` with the trainer name or the Pokémon name, ending with a `\n` using the `fmt` package.
For example, `helloTrainer` prints `Hello Sacha\n` if the input is Sacha.

Add a `helloworld_test.go` file and use this test to check your code:

```go
func TestHelloworldWorker_HelloTrainer(t *testing.T) {
	// Arrange
	helloworldWorker := HelloworldWorker{}
	trainer := mocks.Blue()
	pokemon := mocks.Rattata()
	orig := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w
	defer func() {
		os.Stdout = orig
	}()
	expectedOutput := "Hello " + trainer.Name + "\n"

	// Act
	err = helloworldWorker.helloTrainer(trainer, pokemon)

	// Assert
	w.Close()
	os.Stdout = orig
	require.NoError(t, err)
	out, err := io.ReadAll(r)
	require.NoError(t, err)
	require.Equal(t, expectedOutput, string(out))
}

func TestHelloworldWorker_HelloPokemon(t *testing.T) {
	// Arrange
	helloworldWorker := HelloworldWorker{}
	trainer := mocks.Blue()
	pokemon := mocks.Rattata()
	orig := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w
	defer func() {
		os.Stdout = orig
	}()
	expectedOutput := "Hello " + pokemon.Name + "\n"

	// Act
	err = helloworldWorker.helloPokemon(trainer, pokemon)

	// Assert
	w.Close()
	os.Stdout = orig
	require.NoError(t, err)
	out, err := io.ReadAll(r)
	require.NoError(t, err)
	require.Equal(t, expectedOutput, string(out))
}
```

## Workflow definition

Last but not least, let's implement the `WorkflowDefiner` interface on your structure.

Add this method to your object:

```go
func (c *HelloworldWorker) Definition() *WorkflowDefinition
```

And implement it to return a `WorkflowDefinition` containing your two steps (first is helloTrainer, second is helloPokemon).

Some information:
- Workflow name: `HelloWorld`
- Step 1 name: `Hello Trainer`
- Step 2 name: `Hello Pokémon`

Add this test to check your code:

```go
func TestHelloworldWorker_Definition(t *testing.T) {
	// Arrange
	helloworldWorker := HelloworldWorker{}
	require := require.New(t)

	// Act
	workflowDefinition := helloworldWorker.Definition()

	// Assert
	require.NotNil(workflowDefinition)
	require.Equal("HelloWorld", workflowDefinition.Name)
	require.Len(workflowDefinition.Steps, 2)
	require.Equal("Hello Trainer", workflowDefinition.Steps[0].Name)
	require.True(strings.Contains(runtime.FuncForPC(reflect.ValueOf(workflowDefinition.Steps[0].Do).Pointer()).Name(), "helloTrainer"))
	require.Equal("Hello Pokémon", workflowDefinition.Steps[1].Name)
	require.True(strings.Contains(runtime.FuncForPC(reflect.ValueOf(workflowDefinition.Steps[1].Do).Pointer()).Name(), "helloPokemon"))
}
```

Exercise 2 is now done! You can buzz.

## BONUS: feature to check for already performed steps

Now you've buzzed, if you still have time, you can implement the check of already performed steps.

The goal is to avoid executing a step that have already been executed if your orchestrator restart and run the same workflow execution.

If you check the `WorkflowExecution` struct, there's a field `PerformedSteps` you can use to know what steps have already been performed on the current execution.

You can test your code with this test:

```go
func TestOrchestrator_RunWorkflow_StepAlreadyPerformed(t *testing.T) {
	// Arrange
	var step2done bool
	blue := mocks.Blue()
	rattata := mocks.Rattata()
	testWorkflow := &TestWorkflow{
		steps: []*Step{
			{
				Name: "Step1",
				Do: func(trainer *pokemon.Trainer, pokemon *pokemon.Pokemon) error {
					return errors.New("should not be executed")
				},
			},
			{
				Name: "Step2",
				Do: func(trainer *pokemon.Trainer, pokemon *pokemon.Pokemon) error {
					require.Equal(t, blue, trainer)
					require.Equal(t, rattata, pokemon)

					step2done = true
					return nil
				},
			},
		},
	}
	worker := NewOrchestrator()
	worker.Register(testWorkflow)

	// Act
	err := worker.RunWorkflow(&WorkflowExecution{
		WorkflowName:   "TestWorkflow",
		PerformedSteps: []string{"Step1"},
		Pokemon:        rattata,
		Trainer:        blue,
	})

	// Assert
	require.NoError(t, err)
	require.True(t, step2done)
}
```