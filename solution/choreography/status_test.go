package choreography

import (
	"context"
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
	statusMocks.AssertExpectations(t)
}

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
	statusMocks.AssertExpectations(t)
}
