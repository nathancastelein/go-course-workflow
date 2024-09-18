package sequential

import (
	"context"
	"testing"

	"github.com/nathancastelein/go-course-workflows/mocks"
	"github.com/stretchr/testify/require"
)

func TestCatchProcess(t *testing.T) {
	// Arrange
	statusMock := new(mocks.StatusService)
	combatMock := new(mocks.CombatService)
	pokeballMock := new(mocks.PokeballService)
	catchProcess := NewCatchProcess(statusMock, combatMock, pokeballMock)
	trainer := mocks.Blue()
	pokemon := mocks.Rattata()

	statusMock.On("Paralyze", pokemon).Return(nil)
	combatMock.On("Attack", pokemon).Return(nil)
	pokeballMock.On("Throw", trainer, pokemon).Return(nil)

	// Act
	err := catchProcess.Do(context.Background(), trainer, pokemon)

	// Assert
	require.NoError(t, err)
}
