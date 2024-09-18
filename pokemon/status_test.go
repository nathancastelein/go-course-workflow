package pokemon_test

import (
	"testing"

	"github.com/nathancastelein/go-course-workflows/mocks"
	"github.com/nathancastelein/go-course-workflows/pokemon"
	"github.com/stretchr/testify/require"
)

func TestStatusService_Paralyze(t *testing.T) {
	// Arrange
	statusService := pokemon.NewStatusService()
	rattata := mocks.Rattata()

	// Act
	err := statusService.Paralyze(rattata)

	// Assert
	require.NoError(t, err)
	require.Equal(t, pokemon.StatusParalyzed, rattata.Status)
}
