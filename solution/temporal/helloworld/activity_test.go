package helloworld

import (
	"testing"

	"github.com/nathancastelein/go-course-workflows/mocks"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

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
