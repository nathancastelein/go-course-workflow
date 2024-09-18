package pokemon

import (
	"log/slog"
	"math/rand/v2"
)

func Mewtwo() *Pokemon {
	return &Pokemon{
		ID:            rand.IntN(10000),
		Name:          "Mewtwo",
		Level:         50,
		CurrentHealth: 200,
		MaxHealth:     200,
		Status:        StatusHealthy,
	}
}

func Pikachu() *Pokemon {
	return &Pokemon{
		ID:            rand.IntN(10000),
		Name:          "Pikachu",
		Level:         60,
		CurrentHealth: 230,
		MaxHealth:     230,
		Status:        StatusHealthy,
	}
}

type Pokemon struct {
	ID            int
	Name          string
	Level         int
	CurrentHealth int
	MaxHealth     int
	Status        Status
	TrainerName   string
}

func (p *Pokemon) LogValue() slog.Value {
	group := slog.GroupValue(
		slog.Int("id", p.ID),
		slog.String("name", p.Name),
		slog.Any("status", p.Status),
		slog.Int("current_health", p.CurrentHealth),
		slog.Int("max_health", p.MaxHealth),
		slog.Bool("wild", p.TrainerName == ""),
	)

	return group
}
