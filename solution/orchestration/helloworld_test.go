package orchestration

import (
	"io"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/nathancastelein/go-course-workflows/mocks"
	"github.com/stretchr/testify/require"
)

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

func TestHelloworldWorker_Definition(t *testing.T) {
	// Arrange
	helloworldWorker := HelloworldWorker{}
	require := require.New(t)

	// Act
	workflowDefinition := helloworldWorker.Definition()

	// Assert
	require.NotNil(workflowDefinition)
	require.Equal(HelloworldWorkflowName, workflowDefinition.Name)
	require.Len(workflowDefinition.Steps, 2)
	require.Equal("Hello Trainer", workflowDefinition.Steps[0].Name)
	require.True(strings.Contains(runtime.FuncForPC(reflect.ValueOf(workflowDefinition.Steps[0].Do).Pointer()).Name(), "helloTrainer"))
	require.Equal("Hello Pokémon", workflowDefinition.Steps[1].Name)
	require.True(strings.Contains(runtime.FuncForPC(reflect.ValueOf(workflowDefinition.Steps[1].Do).Pointer()).Name(), "helloPokemon"))
}
