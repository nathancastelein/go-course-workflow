# Exercise 2: The Combat worker

This worker will have quite the same look than the status one.

## Write the combat worker

The `CombatWorker` structure requires three information to work properly:
- a `pokemon.Combat` to call the `Attack` method
- a first `MessageBrokerTopic` to subscribe for `PokemonParalyzed` events
- a second `MessageBrokerTopic` to publish `PokemonWeakened` events

Define this structure in a `combat.go` file.

## Provide a constructor

As we did for the `StatusWorker`, provide a constructor with this signature:

```go
func NewCombatWorker(combat pokemon.CombatService, combatTopic MessageBrokerTopic, pokeballTopic MessageBrokerTopic) CombatWorker
```

Create a `combat_test.go` file and add this first test:

```go
func TestCombatWorker_New(t *testing.T) {
	// Arrange
	// Data
	combatTopic := make(MessageBrokerTopic, 1)
	pokeballTopic := make(MessageBrokerTopic, 1)

	// Mocks
	combatMocks := new(mocks.CombatService)

	// Act
	combatWorker := NewCombatWorker(
		combatMocks,
		combatTopic,
		pokeballTopic,
	)

	// Assert
	require.NotNil(t, combatWorker)
}
```

## Write the Run method

It's now time to write the Run method:

```go
func (c *CombatWorker) Run(ctx context.Context) 
```

You can start from this skeleton:

```go
func (c *CombatWorker) Run(ctx context.Context) {
	slog.Info("starting combat worker")
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-c.combatTopic:
			// Add missing code here
		}
	}
}
```

You just need to add some missing code:
- check the event type
- if event type is `PokemonParalyzed`, then attack the Pokémon using the Combat
- in case of error, log then continue
- send an new event of type `PokemonWeakened` in the pokeball topic

Add this test on your `combat_test.go` to test your function:

```go
func TestCombatWorker_Run(t *testing.T) {
	// Arrange
	// Data
	trainer := mocks.Blue()
	rattata := mocks.Rattata()
	combatTopic := make(MessageBrokerTopic, 1)
	pokeballTopic := make(MessageBrokerTopic, 1)

	// Mocks
	combatMocks := new(mocks.CombatService)
	combatMocks.On("Attack", rattata).Return(nil)

	combatWorker := NewCombatWorker(
		combatMocks,
		combatTopic,
		pokeballTopic,
	)

	ctx, cancel := context.WithCancel(context.Background())
	var pokeballEvent Event

	go func() {
		pokeballEvent = <-pokeballTopic
		cancel()
	}()

	// Act
	combatTopic <- Event{
		Type:    PokemonParalyzed,
		Pokemon: rattata,
		Trainer: trainer,
	}
	combatWorker.Run(ctx)

	// Assert
	require.Equal(t, PokemonWeakened, pokeballEvent.Type)
	require.Equal(t, trainer, pokeballEvent.Trainer)
	require.Equal(t, rattata, pokeballEvent.Pokemon)
	require.Empty(t, combatTopic)
	require.Empty(t, pokeballTopic)
	combatMocks.AssertExpectations(t)
}
```