package choreography

import (
	"context"
	"testing"

	"github.com/nathancastelein/go-course-workflows/mocks"
	"github.com/stretchr/testify/require"
)

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
}
