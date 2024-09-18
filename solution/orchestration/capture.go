package orchestration

import (
	"github.com/google/uuid"
	"github.com/nathancastelein/go-course-workflows/pokemon"
)

var (
	CapturePokemonWorkflowName = "CapturePokemon"
)

type CapturePokemonWorker struct {
	status   pokemon.StatusService
	combat   pokemon.CombatService
	pokeball pokemon.PokeballService
}

func NewCapturePokemonWorker(
	status pokemon.StatusService,
	combat pokemon.CombatService,
	pokeball pokemon.PokeballService,
) *CapturePokemonWorker {
	return &CapturePokemonWorker{
		status:   status,
		combat:   combat,
		pokeball: pokeball,
	}
}

func (c *CapturePokemonWorker) attack(trainer *pokemon.Trainer, pokemon *pokemon.Pokemon) error {
	return c.combat.Attack(pokemon)
}

func (c *CapturePokemonWorker) paralyze(trainer *pokemon.Trainer, pokemon *pokemon.Pokemon) error {
	return c.status.Paralyze(pokemon)
}

func (c *CapturePokemonWorker) throwPokeball(trainer *pokemon.Trainer, pokemon *pokemon.Pokemon) error {
	return c.pokeball.Throw(trainer, pokemon)
}

func (c *CapturePokemonWorker) Definition() *WorkflowDefinition {
	return &WorkflowDefinition{
		Name: CapturePokemonWorkflowName,
		Steps: []*Step{
			{
				Name: "Paralyze",
				Do:   c.paralyze,
			},
			{
				Name: "Attack",
				Do:   c.attack,
			},
			{
				Name: "ThrowPokeball",
				Do:   c.throwPokeball,
			},
		},
	}
}

func NewCapturePokemonWorkflowExecution(trainer *pokemon.Trainer, pokemon *pokemon.Pokemon) *WorkflowExecution {
	return &WorkflowExecution{
		ID:             uuid.New().String(),
		WorkflowName:   CapturePokemonWorkflowName,
		PerformedSteps: []string{},
		Trainer:        trainer,
		Pokemon:        pokemon,
	}
}
