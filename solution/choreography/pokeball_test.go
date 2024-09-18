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
