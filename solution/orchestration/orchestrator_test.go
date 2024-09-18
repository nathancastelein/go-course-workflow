package orchestration

import (
	"errors"
	"testing"

	"github.com/nathancastelein/go-course-workflows/mocks"
	"github.com/nathancastelein/go-course-workflows/pokemon"
	"github.com/stretchr/testify/require"
)

func TestNewOrchestrator(t *testing.T) {
	// Arrange

	// Act
	worker := NewOrchestrator()

	// Assert
	require.NotNil(t, worker)
	require.NotNil(t, worker.workflows)
}

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
