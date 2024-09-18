package capture

import (
	"context"
	"testing"

	"github.com/nathancastelein/go-course-workflows/mocks"
	"github.com/nathancastelein/go-course-workflows/pokemon"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

func TestParalyzeActivity(t *testing.T) {
	// Arrange
	testSuite := &testsuite.WorkflowTestSuite{}
	statusService := new(mocks.StatusService)
	rattata := mocks.Rattata()
	paralyzedRattata := mocks.Rattata()
	paralyzedRattata.Status = pokemon.StatusParalyzed
	env := testSuite.NewTestActivityEnvironment()
	worker := NewWorker(
		statusService,
		nil,
		nil,
	)
	env.RegisterActivity(worker.ParalyzeActivity)
	statusService.On("Paralyze", rattata).Run(func(args mock.Arguments) {
		argRattata := args.Get(0).(*pokemon.Pokemon)
		argRattata.Status = pokemon.StatusParalyzed
	}).Return(nil)

	// Act
	val, err := env.ExecuteActivity(worker.ParalyzeActivity, rattata)

	// Assert
	require.NoError(t, err)
	var res pokemon.Pokemon
	require.NoError(t, val.Get(&res))
	require.Equal(t, paralyzedRattata, &res)
}

func TestAttackActivity(t *testing.T) {
	// Arrange
	testSuite := &testsuite.WorkflowTestSuite{}
	combatService := new(mocks.CombatService)
	rattata := mocks.Rattata()
	env := testSuite.NewTestActivityEnvironment()
	worker := NewWorker(
		nil,
		combatService,
		nil,
	)
	env.RegisterActivity(worker.AttackActivity)
	combatService.On("Attack", rattata).Return(nil)

	// Act
	val, err := env.ExecuteActivity(worker.AttackActivity, rattata)

	// Assert
	require.NoError(t, err)
	var res pokemon.Pokemon
	require.NoError(t, val.Get(&res))
	require.Equal(t, rattata, &res)
}

func TestThrowPokeballActivity(t *testing.T) {
	// Arrange
	testSuite := &testsuite.WorkflowTestSuite{}
	pokeballService := new(mocks.PokeballService)
	rattata := mocks.Rattata()
	blue := mocks.Blue()
	env := testSuite.NewTestActivityEnvironment()
	worker := NewWorker(
		nil,
		nil,
		pokeballService,
	)
	env.RegisterActivity(worker.ThrowPokeballActivity)
	pokeballService.On("Throw", blue, rattata).Return(nil)

	// Act
	val, err := env.ExecuteActivity(worker.ThrowPokeballActivity, blue, rattata)

	// Assert
	require.NoError(t, err)
	var res ThrowPokeballOutput
	require.NoError(t, val.Get(&res))
	require.Equal(t, rattata, res.Pokemon)
	require.Equal(t, blue, res.Trainer)
}

type Object struct {
	Value string
}

func SimpleActivity(ctx context.Context, object *Object) error {
	object.Value = "foobar"
	return nil
}

func SimpleActivityWorking(ctx context.Context, object *Object) (*Object, error) {
	object.Value = "foobar"
	return object, nil
}

func TestSimpleActivity(t *testing.T) {
	// Arrange
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestActivityEnvironment()
	env.RegisterActivity(SimpleActivity)
	object := &Object{Value: "foo"}

	// Act
	_, err := env.ExecuteActivity(SimpleActivity, object)

	// Assert
	require.NoError(t, err)
	// Not working due to Temporal serialization and deserialization
	//require.Equal(t, "foobar", object.Value)
}

func TestSimpleActivityWorking(t *testing.T) {
	// Arrange
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestActivityEnvironment()
	env.RegisterActivity(SimpleActivityWorking)
	object := &Object{Value: "foo"}

	// Act
	val, err := env.ExecuteActivity(SimpleActivityWorking, object)

	// Assert
	require.NoError(t, err)
	require.NoError(t, val.Get(object))
	require.Equal(t, "foobar", object.Value)
}
