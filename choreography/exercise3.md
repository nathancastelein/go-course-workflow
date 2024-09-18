# The Pokeball worker

You should now be used to write workers for our choreography pattern.

Here is the last worker:

```go
type PokeballWorker struct {
	pokeball      pokemon.PokeballService
	pokeballTopic MessageBrokerTopic
	captureTopic    MessageBrokerTopic
}
```

Write all the missing code in a `pokeball.go` file to make the worker usable:
- its constructor
- its Run method

You can use this test suite to test your code:

```go
package choreography

import (
	"context"
	"testing"

	"github.com/nathancastelein/go-course-workflows/mocks"
	"github.com/stretchr/testify/require"
)

func TestPokeballWorker_Run(t *testing.T) {
	// Arrange
	// Data
	trainer := mocks.Blue()
	rattata := mocks.Rattata()
	pokeballTopic := make(MessageBrokerTopic, 1)
	captureTopic := make(MessageBrokerTopic, 1)

	// Mocks
	pokeballMocks := new(mocks.PokeballService)
	pokeballMocks.On("Throw", trainer, rattata).Return(nil)

	combatWorker := NewPokeballWorker(
		pokeballMocks,
		pokeballTopic,
		captureTopic,
	)

	ctx, cancel := context.WithCancel(context.Background())
	var captureEvent Event

	go func() {
		captureEvent = <-captureTopic
		cancel()
	}()

	// Act
	pokeballTopic <- Event{
		Type:    PokemonWeakened,
		Pokemon: rattata,
		Trainer: trainer,
	}
	combatWorker.Run(ctx)

	// Assert
	require.Equal(t, PokemonCaptured, captureEvent.Type)
	require.Equal(t, trainer, captureEvent.Trainer)
	require.Equal(t, rattata, captureEvent.Pokemon)
	require.Empty(t, pokeballTopic)
	require.Empty(t, captureTopic)
}

```

# Test the complete workflow

It's now time to test your workflow from start to end.

Add this test in `capture_test.go` file:

```go
package choreography

import (
	"context"
	"testing"

	"github.com/nathancastelein/go-course-workflows/mocks"
	"github.com/stretchr/testify/require"
)

func TestCatchChoreography(t *testing.T) {
	// Arrange
	// Data
	trainer := mocks.Blue()
	rattata := mocks.Rattata()
	statusTopic := make(MessageBrokerTopic, 1)
	combatTopic := make(MessageBrokerTopic, 1)
	pokeballTopic := make(MessageBrokerTopic, 1)
	captureTopic := make(MessageBrokerTopic, 1)

	// Mocks
	statusMocks := new(mocks.StatusService)
	statusMocks.On("Paralyze", rattata).Return(nil)
	combatMocks := new(mocks.CombatService)
	combatMocks.On("Attack", rattata).Return(nil)
	pokeballMocks := new(mocks.PokeballService)
	pokeballMocks.On("Throw", trainer, rattata).Return(nil)

	statusWorker := NewStatusWorker(
		statusMocks,
		statusTopic,
		combatTopic,
	)

	combatWorker := NewCombatWorker(
		combatMocks,
		combatTopic,
		pokeballTopic,
	)

	pokeballWorker := NewPokeballWorker(
		pokeballMocks,
		pokeballTopic,
		captureTopic,
	)

	ctx, cancel := context.WithCancel(context.Background())

	// Act
	statusTopic <- Event{
		Type:    PokemonEncountered,
		Pokemon: rattata,
		Trainer: trainer,
	}
	go statusWorker.Run(ctx)
	go combatWorker.Run(ctx)
	go pokeballWorker.Run(ctx)

	captureEvent := <-captureTopic
	cancel()

	// Assert
	require.Equal(t, PokemonCaptured, captureEvent.Type)
	require.Equal(t, trainer, captureEvent.Trainer)
	require.Equal(t, rattata, captureEvent.Pokemon)
	require.Empty(t, pokeballTopic)
	require.Empty(t, captureTopic)
}
```

Run the test: everything should work!

You can buzz now.

# BONUS: The main program

If you have a bit of time, creates a folder `choreography/cmd` and add a `main.go`.

Try to write a main function to start your three workers, properly connected with their topics, and send an event to capture a Pokémon!

You can use the previous test as a base.

Use this snippet to manually wait for the end of your program:

```go
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Waiting input to exit")
	reader.ReadString('\n')
```