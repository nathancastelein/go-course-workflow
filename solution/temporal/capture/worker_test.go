package capture

import (
	"testing"

	"github.com/nathancastelein/go-course-workflows/mocks"
	"github.com/stretchr/testify/require"
)

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
