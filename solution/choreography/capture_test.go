package choreography

import (
	"context"
	"testing"

	"github.com/nathancastelein/go-course-workflows/mocks"
	"github.com/stretchr/testify/require"
)

func TestCaptureChoreography(t *testing.T) {
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
	statusMocks.AssertExpectations(t)
	combatMocks.AssertExpectations(t)
	pokeballMocks.AssertExpectations(t)
}
