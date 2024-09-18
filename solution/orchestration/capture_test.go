package orchestration

import (
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/nathancastelein/go-course-workflows/mocks"
	"github.com/stretchr/testify/require"
)

func TestNewCapturePokemonWorker(t *testing.T) {
	// Arrange
	statusMock := new(mocks.StatusService)
	combatMock := new(mocks.CombatService)
	pokeballMock := new(mocks.PokeballService)
	expected := &CapturePokemonWorker{
		status:   statusMock,
		combat:   combatMock,
		pokeball: pokeballMock,
	}

	// Act
	got := NewCapturePokemonWorker(statusMock, combatMock, pokeballMock)

	// Assert
	require.Equal(t, expected, got)
}

func TestCaptureWorker_Paralyze(t *testing.T) {
	// Arrange
	statusMocks := new(mocks.StatusService)
	workflow := NewCapturePokemonWorker(statusMocks, nil, nil)
	blue := mocks.Blue()
	rattata := mocks.Rattata()

	statusMocks.On("Paralyze", rattata).Return(nil)

	// Act
	err := workflow.paralyze(blue, rattata)

	// Assert
	require.NoError(t, err)
}

func TestCaptureWorker_Attack(t *testing.T) {
	// Arrange
	combatMocks := new(mocks.CombatService)
	workflow := NewCapturePokemonWorker(nil, combatMocks, nil)
	blue := mocks.Blue()
	rattata := mocks.Rattata()

	combatMocks.On("Attack", rattata).Return(nil)

	// Act
	err := workflow.attack(blue, rattata)

	// Assert
	require.NoError(t, err)
}

func TestCaptureWorker_ThrowPokeball(t *testing.T) {
	// Arrange
	pokeballMocks := new(mocks.PokeballService)
	workflow := NewCapturePokemonWorker(nil, nil, pokeballMocks)
	blue := mocks.Blue()
	rattata := mocks.Rattata()

	pokeballMocks.On("Throw", blue, rattata).Return(nil)

	// Act
	err := workflow.throwPokeball(blue, rattata)

	// Assert
	require.NoError(t, err)
}

func TestCaptureWorker_Definition(t *testing.T) {
	// Arrange
	captureWorker := NewCapturePokemonWorker(nil, nil, nil)
	require := require.New(t)

	// Act
	workflowDefinition := captureWorker.Definition()

	// Assert
	require.NotNil(workflowDefinition)
	require.Equal(CapturePokemonWorkflowName, workflowDefinition.Name)
	require.Len(workflowDefinition.Steps, 3)
	require.Equal("Paralyze", workflowDefinition.Steps[0].Name)
	require.True(strings.Contains(runtime.FuncForPC(reflect.ValueOf(workflowDefinition.Steps[0].Do).Pointer()).Name(), "paralyze"))
	require.Equal("Attack", workflowDefinition.Steps[1].Name)
	require.True(strings.Contains(runtime.FuncForPC(reflect.ValueOf(workflowDefinition.Steps[1].Do).Pointer()).Name(), "attack"))
	require.Equal("ThrowPokeball", workflowDefinition.Steps[2].Name)
	require.True(strings.Contains(runtime.FuncForPC(reflect.ValueOf(workflowDefinition.Steps[2].Do).Pointer()).Name(), "throwPokeball"))
}

func TestNewCapturePokemonWorkflowExecution(t *testing.T) {
	// Arrange
	blue := mocks.Blue()
	rattata := mocks.Rattata()

	// Act
	got := NewCapturePokemonWorkflowExecution(blue, rattata)

	// Assert
	require.Equal(t, CapturePokemonWorkflowName, got.WorkflowName)
	require.Equal(t, []string{}, got.PerformedSteps)
	require.Equal(t, blue, got.Trainer)
	require.Equal(t, rattata, got.Pokemon)
}
