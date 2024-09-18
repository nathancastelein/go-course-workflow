# Exercise 3

Following everything you've learnt from the previous exercises, it's now time to write your business workflow: Capture a Pokémon!

It's quite the same work we did for `Helloworld` workflow, except that our worker will use our three services `StatusService`, `CombatService`, and `PokeballService` to call the proper functions.

## Worker

Create new files `capture.go` and `capture_test.go`.

Write a structure `CapturePokemonWorker`. This struct has three fields, one for each service:
- status
- combat
- pokeball

Add a constructor `NewCapturePokemonWorker` with the three services as inputs, to instantiate your worker.

```go
func NewCapturePokemonWorker(status pokemon.StatusService, combat pokemon.CombatService, pokeball pokemon.PokeballService) *CapturePokemonWorker
```

Test with this code:

```go
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
```

## Attack, paralyze and throw pokeball

Add three methods on your worker, defined as `Step`, to attack, paralyze and throw pokeball:

```go
func (c *CapturePokemonWorker) attack(trainer *pokemon.Trainer, pokemon *pokemon.Pokemon) error

func (c *CapturePokemonWorker) paralyze(trainer *pokemon.Trainer, pokemon *pokemon.Pokemon) error

func (c *CapturePokemonWorker) throwPokeball(trainer *pokemon.Trainer, pokemon *pokemon.Pokemon) error
```

Implement each method.

Test your code:

```go
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
```

## Workflow definition

Last but no least, implement the `WorkflowDefiner` interface to return the definition of your capture workflow:

```go
func (c *CapturePokemonWorker) Definition() *WorkflowDefinition
```

Some information:
- Workflow name: CapturePokemon
- Step 1 name: Paralyze
- Step 2 name: Attack
- Step 3 name: ThrowPokeball

You can use this test to test your work:

```go
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
```


Exercise 3 is now done! You can buzz.

## BONUS: provide a constructure for your workflow execution

To make things easier for people who want to start your workflow, you can provide a constructor for a `WorkflowExecution` for your workflow.

The signature is this one:

```go
func NewCapturePokemonWorkflowExecution(trainer *pokemon.Trainer, pokemon *pokemon.Pokemon) *WorkflowExecution
```

You can use `uuid.New().String()` from `github.com/google/uuid` to generate an UUID.

Check your code with this test:

```go
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
```