package orchestration

import (
	"fmt"

	"github.com/nathancastelein/go-course-workflows/pokemon"
)

var (
	HelloworldWorkflowName = "HelloWorld"
)

type HelloworldWorker struct{}

func (c *HelloworldWorker) helloTrainer(trainer *pokemon.Trainer, pokemon *pokemon.Pokemon) error {
	fmt.Printf("Hello %s\n", trainer.Name)
	return nil
}

func (c *HelloworldWorker) helloPokemon(trainer *pokemon.Trainer, pokemon *pokemon.Pokemon) error {
	fmt.Printf("Hello %s\n", pokemon.Name)
	return nil
}

func (c *HelloworldWorker) Definition() *WorkflowDefinition {
	return &WorkflowDefinition{
		Name: HelloworldWorkflowName,
		Steps: []*Step{
			{
				Name: "Hello Trainer",
				Do:   c.helloTrainer,
			},
			{
				Name: "Hello Pokémon",
				Do:   c.helloPokemon,
			},
		},
	}
}
