package capture

import (
	"testing"

	"github.com/nathancastelein/go-course-workflows/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

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
