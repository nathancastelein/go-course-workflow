# Exercise 1: The Status worker

Let's start with our first worker, the status worker, in charge of paralyzing the Pokémon.

## Understanding events

First of all, have a look on the [event file](./event.go).
You will find the definition of an `Event` in our system.

Events are sent to a `MessageBrokerTopic`. It represents a very simple implementation of a message broker topic, using Go channels.

Each event has an `EventType`, representing one of the four types of event we will need in our system.

## Write the status worker

Now we will create a `StatusWorker` structure.

This struct requires three information to work properly:
- a `pokemon.StatusService` to call the `Paralyze` method
- a first `MessageBrokerTopic` to subscribe for `PokemonEncountered` events
- a second `MessageBrokerTopic` to publish `PokemonParalyzed` events

Define this structure in a `status.go` file.

## Provide a constructor

Write a new method to instantiate your `StatusWorker`.

This is the expected signature for your constructor:

```go
func NewStatusWorker(status pokemon.StatusService, statusTopic MessageBrokerTopic, combatTopic MessageBrokerTopic) StatusWorker
```

In order to test your method, create a `status_test.go` file.

Add this simple unit test:

```go
package choreography

import (
	"testing"

	"github.com/nathancastelein/go-course-workflows/mocks"
	"github.com/stretchr/testify/require"
)

func TestStatusWorker_New(t *testing.T) {
	// Arrange
	// Data
	statusTopic := make(MessageBrokerTopic, 1)
	combatTopic := make(MessageBrokerTopic, 1)

	// Mocks
	statusMocks := new(mocks.StatusService)

	// Act
	statusWorker := NewStatusWorker(
		statusMocks,
		statusTopic,
		combatTopic,
	)

	// Assert
	require.NotNil(t, statusWorker)
}
```

Run the test, everything should be alright.

## Write the Run method

Now it's time to implement a new method on your `StatusWorker` structure.

This method will have this signature:

```go
func (s *StatusWorker) Run(ctx context.Context) 
```

This method will be in charge of the following algorithm running in an infinite loop:
- Listen to events from the status topic
- In case of a `PokemonEncountered` event, paralyze the Pokémon
- Then send an event of type `PokemonParalyzed` to the combat topic
- And don't forget to listen to the `ctx.Done()` channel to stop the process when context is cancelled

### Handle context cancellation

Let's start simply by listening to `ctx.Done()` event:
- Start an infinite loop
- Listen to events with a `select` statement
- If an event comes from the `ctx.Done()` channel, just `return`

You can add this test on your test file to test your code:

```go
type spyContext struct {
	context.Context
	doneCalled bool
}

func (t *spyContext) Done() <-chan struct{} {
	t.doneCalled = true
	return t.Context.Done()
}

func TestStatusWorker_RunWithCancel(t *testing.T) {
	// Arrange
	// Data
	statusTopic := make(MessageBrokerTopic, 1)
	combatTopic := make(MessageBrokerTopic, 1)

	// Mocks
	statusMocks := new(mocks.StatusService)

	statusWorker := NewStatusWorker(
		statusMocks,
		statusTopic,
		combatTopic,
	)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	spyContext := &spyContext{Context: ctx}

	// Act
	statusWorker.Run(spyContext)

	// Assert
	require.True(t, spyContext.doneCalled)
}
```

### Listen to events

Let's now implement the business code to handle events:
- add a case on your select statement to listen from events on the status topic channel
- check the event type
- if event type is `PokemonEncountered`, then paralyze the Pokémon using the StatusService
- in case of error, log then continue
- send an new event of type `PokemonParalyzed` in the combat topic

Add this test on your `status_test.go` to test your function:

```go
func TestStatusWorker_Run(t *testing.T) {
	// Arrange
	// Data
	trainer := mocks.Blue()
	rattata := mocks.Rattata()
	statusTopic := make(MessageBrokerTopic, 1)
	combatTopic := make(MessageBrokerTopic, 1)

	// Mocks
	statusMocks := new(mocks.StatusService)
	statusMocks.On("Paralyze", rattata).Return(nil)

	statusWorker := NewStatusWorker(
		statusMocks,
		statusTopic,
		combatTopic,
	)

	ctx, cancel := context.WithCancel(context.Background())
	var combatEvent Event

	go func() {
		combatEvent = <-combatTopic
		cancel()
	}()

	// Act
	statusTopic <- Event{
		Type:    PokemonEncountered,
		Pokemon: rattata,
		Trainer: trainer,
	}
	statusWorker.Run(ctx)

	// Assert
	require.Equal(t, PokemonParalyzed, combatEvent.Type)
	require.Equal(t, trainer, combatEvent.Trainer)
	require.Equal(t, rattata, combatEvent.Pokemon)
	require.Empty(t, combatTopic)
	require.Empty(t, statusTopic)
}
```