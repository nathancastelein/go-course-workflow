package durabletask

import "github.com/nathancastelein/go-course-workflows/pokemon"

type Worker struct {
	status   pokemon.StatusService
	combat   pokemon.CombatService
	pokeball pokemon.PokeballService
}

func NewWorker(
	status pokemon.StatusService,
	combat pokemon.CombatService,
	pokeball pokemon.PokeballService,
) *Worker {
	worker := &Worker{
		status:   status,
		combat:   combat,
		pokeball: pokeball,
	}

	return worker
}
