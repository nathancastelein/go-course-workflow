# Exercise 1

In this first exercise, we will focus on our orchestrator.

It will be an important part to manage a workflow. We will try to write an orchestrator agnostic of our capture workflow.
Having this in mind will help to re-use our orchestrator for different workflows!

Let's first have a look on some predefined and provided structures.

## Workflow definition

The [workflowDefinition.go](workflowDefinition.go) files provides two structures (`WorkflowDefinition` and `WorkflowDefinitions`) and an interface `WorkflowDefiner`.

`WorkflowDefinition` is a structure which contains the definition of a workflow: its name, and the steps to execute sequentially.

`WorkflowDefinitions` is just an array of `WorkflowDefinition` with a `FindByName` method to search for a given workflow name in an array.

`WorkflowDefiner` is the interface to meet to be considered as a workflow by our orchestrator. It contains a simple method to return the `WorkflowDefinition` of the workflow.

## Step

The [step.go](./step.go) provides the definition of a single workflow's step.

A `Step` is defined by its `Name` and a `Do` function to hold the code linked to the step.

## Workflow execution

Finally, the [workflowExecution.go](workflowExecution.go) file provides the structure representing a workflow execution.

A `WorkflowExecution` has a unique `ID`, the `WorkflowName` to execute, the already `PerformedSteps`, and the input `Trainer` and `Pokemon` this execution handle.

*Notice: the WorkflowExecution holds a trainer and a pokemon, two objects that are a bit business oriented. If we want an orchestrator more generic, those two fields should be changed by an `Input []byte` field, holding input serialized properly. It is a trade-off accepted for readibility in this exercise*

# Orchestrator

It's time to focus on our orchestrator. Don't worry, we will go step by step!

## The structure

The `Orchestrator` will be a structure. The `Orchestrator` struct maintains only one field: the list of `WorkflowDefinition` the orchestrator can handle.

Creates a new file `orchestrator.go`. Add a structure `Orchestrator` with a single field `workflows`, of `WorkflowDefinitions` type.

Provide a constructor `NewOrchestrator() *Orchestrator` to return an properly initialyzed object.

Add a `orchestrator_test.go` file and add this test to check your code:

```go
func TestNewOrchestrator(t *testing.T) {
	// Arrange

	// Act
	worker := NewOrchestrator()

	// Assert
	require.NotNil(t, worker)
	require.NotNil(t, worker.workflows)
}
```

## Register workflows

Add now a method on your orchestrator to register a new workflow definition to handle.

This is the signature of the expected method: 
```go
func (o *Orchestrator) Register(definer WorkflowDefiner)
```

The new definition should be append to the `workflows` list of the orchestrator.

Write this method!

Add this code on your test file to check your code:

```go
type TestWorkflow struct {
	steps []*Step
}

func (t *TestWorkflow) Definition() *WorkflowDefinition {
	return &WorkflowDefinition{
		Name:  "TestWorkflow",
		Steps: t.steps,
	}
}

func TestOrchestrator_Register(t *testing.T) {
	// Arrange
	worker := NewOrchestrator()
	testWorkflow := &TestWorkflow{
		steps: []*Step{
			{
				Name: "Dummy",
				Do: func(trainer *pokemon.Trainer, pokemon *pokemon.Pokemon) error {
					return nil
				},
			},
		},
	}

	// Act
	worker.Register(testWorkflow)

	// Assert
	require.NotNil(t, worker)
	require.NotNil(t, worker.workflows)
	require.Len(t, worker.workflows, 1)
	require.Equal(t, testWorkflow.Definition(), worker.workflows[0])
}
```

## Run workflow

Last but not least, it's now time to write the core method of your orchestrator: the `RunWorkflow` method!

Add a new method on your orchestrator:

```go
func (o *Orchestrator) RunWorkflow(workflowExecution *WorkflowExecution) error
```

Write the algorithm to execute the workflow:
- Find the proper `WorkflowDefinition` to execute based on the `workflowExecution.WorkflowName`
- Loop through each `Step` of the workflow sequentially, and execute the `Do` function, providing the `Trainer` and the `Pokemon` from the `workflowExecution` input

Use this test to check your code:

```go
func TestOrchestrator_RunWorkflow(t *testing.T) {
	// Arrange
	var step1done, step2done bool
	blue := mocks.Blue()
	rattata := mocks.Rattata()
	testWorkflow := &TestWorkflow{
		steps: []*Step{
			{
				Name: "Step1",
				Do: func(trainer *pokemon.Trainer, pokemon *pokemon.Pokemon) error {
					require.Equal(t, blue, trainer)
					require.Equal(t, rattata, pokemon)

					step1done = true
					return nil
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
		PerformedSteps: []string{},
		Pokemon:        rattata,
		Trainer:        blue,
	})

	// Assert
	require.NoError(t, err)
	require.True(t, step1done)
	require.True(t, step2done)
}
```