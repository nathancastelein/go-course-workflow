package pokemon_test

import (
	"testing"

	"github.com/nathancastelein/go-course-workflows/mocks"
	"github.com/nathancastelein/go-course-workflows/pokemon"
	"github.com/stretchr/testify/require"
)

func TestCombatService_Attack(t *testing.T) {
	// Arrange
	statusService := pokemon.NewCombatService()
	rattata := mocks.Rattata()

	// Act
	err := statusService.Attack(rattata)

	// Assert
	require.NoError(t, err)
	require.Equal(t, 1, rattata.CurrentHealth)
}
