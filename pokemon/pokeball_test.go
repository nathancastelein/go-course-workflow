package pokemon_test

import (
	"testing"

	"github.com/nathancastelein/go-course-workflows/mocks"
	"github.com/nathancastelein/go-course-workflows/pokemon"
	"github.com/stretchr/testify/require"
)

func TestPokeballService_Throw(t *testing.T) {
	// Arrange
	statusService := pokemon.NewPokeballService()
	rattata := mocks.Rattata()
	blue := mocks.Blue()

	// Act
	err := statusService.Throw(blue, rattata)

	// Assert
	require.NoError(t, err)
	require.Equal(t, blue.Name, rattata.TrainerName)
	require.Len(t, blue.Pokemons, 2)
	require.Equal(t, rattata, blue.Pokemons[1])
}
