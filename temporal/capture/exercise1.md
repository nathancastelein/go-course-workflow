# Exercise 1

You will see a lot of similarities with everything you've done during this workshop.

This last exercises are a perfect conclusion to write proper workflows with an orchestrator. Let's do it!

## Worker

As usual, the worker will be in charge of executing locally your business code.

In our case, the worker will be in charge of executing each activities, and the workflow.

To do so, the worker will manipulate our three services: Status, Combat and Pokeball.

Create new files `worker.go` and `worker_test.go`.

Create a new `Worker` struct with all three services as fields of the structure (status, combat, pokeball).

Then create a constructor for your type named `NewWorker`.

Test your code with this:

```go
func TestNewWorker(t *testing.T) {
	// Arrange
	statusMock := new(mocks.StatusService)
	combatMock := new(mocks.CombatService)
	pokeballMock := new(mocks.PokeballService)
	expected := &Worker{
		status:   statusMock,
		combat:   combatMock,
		pokeball: pokeballMock,
	}

	// Act
	got := NewWorker(statusMock, combatMock, pokeballMock)

	// Assert
	require.Equal(t, expected, got)
}
```

## Activities

We know our workflow will have three interactions with external services: Paralyze, Attack and ThrowPokeball.

This kind of work fits perfectly in an activity!

Create new files `activities.go` and `activities_test.go`.

Let's start with the `ParalyzeActivity`.

While writing your code for Temporal, you need to understand the behavior of Temporal.

In an activity, inputs and outputs are constantly serialized and deserialized between worker and server.

But if you have a look on the `Paralyze` method of the `StatusService`, this method works with a pointer to a Pokémon, and modify the Pokémon inside the function.

For example, you should be tempted to write code like this:

```go
type Object struct {
	Value string
}

func SimpleActivity(ctx context.Context, object *Object) error {
	object.Value = "foobar"
	return nil
}
```

This will simply not work in a Temporal environment, due to the way Temporal is working.

To handle this situation, you will have to return the modified data as output of the activity:

```go
type Object struct {
	Value string
}

func SimpleActivityWorking(ctx context.Context, object *Object) (*Object, error) {
	object.Value = "foobar"
	return object, nil
}
```

We will need to use this mecanism for all of our activities.

### Paralyze activity

Understanding this, write a method `ParalyzeActivity` on your worker to paralyze a Pokémon.

This method:
- Accept a `context.Context` as first input, as required for an activity
- Accept a `*pokemon.Pokemon` as second input
- Return a `*pokemon.Pokemon`
- Return an `error`

Use this test to check your code:

```go
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
```

### Attack activity

The attack activity is quite the same than the paralyze one.

Write a `AttackActivity`.

Test:

```go
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
```

### Throw Pokéball activity

Last but not least, the Throw Pokéball activity.

It will be a bit different, because the `Throw` method of the `PokeballService` needs a `Trainer` and a `Pokemon` input.

So we will have to accept one more input, which is doable for Temporal.
But if we have two inputs, following what we said previously, we need to return two outputs (and an error).

Returning two outputs is not allowed by Temporal. So we will need to create a type to aggregate the two objects:

```go
type ThrowPokeballOutput struct {
	Pokemon *pokemon.Pokemon
	Trainer *pokemon.Trainer
}
```

Using three inputs parameters (`context.Context`, `*pokemon.Trainer` and `*pokemon.Pokemon`) and two outputs parameters (`ThrowPokeballOutput` and `error`), write the `ThrowPokeballActivity`.

To test your code:

```go
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
```