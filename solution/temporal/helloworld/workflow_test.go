package helloworld

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
	blue := mocks.Blue()
	rattata := mocks.Rattata()

	env.OnActivity(SayHelloToTrainer, mock.Anything, blue).Return("Hello Blue!", nil)
	env.OnActivity(SayHelloToPokemon, mock.Anything, rattata).Return("Hello Rattata!", nil)
	env.OnActivity(SayHelloToProfessorOak, mock.Anything).Return("Hello Professor Oak!", nil)

	// Act
	env.ExecuteWorkflow(Helloworld, blue, rattata)

	// Assert
	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
	var result string
	require.NoError(t, env.GetWorkflowResult(&result))
	require.Equal(t, "Hello Blue! Hello Rattata! Hello Professor Oak!", result)
}
