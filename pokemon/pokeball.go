package pokemon

type PokeballService interface {
	Throw(trainer *Trainer, pokemon *Pokemon) error
}

func NewPokeballService() PokeballService {
	return &pokeballService{}
}

type pokeballService struct{}

func (s *pokeballService) Throw(trainer *Trainer, pokemon *Pokemon) error {
	pokemon.TrainerName = trainer.Name
	trainer.Pokemons = append(trainer.Pokemons, pokemon)
	return nil
}
