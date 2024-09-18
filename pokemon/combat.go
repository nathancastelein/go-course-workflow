package pokemon

type CombatService interface {
	Attack(pokemon *Pokemon) error
}

func NewCombatService() CombatService {
	return &combatService{}
}

type combatService struct{}

func (s *combatService) Attack(pokemon *Pokemon) error {
	pokemon.CurrentHealth = 1
	return nil
}
