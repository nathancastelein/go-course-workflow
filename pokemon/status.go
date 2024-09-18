package pokemon

import "log/slog"

const (
	StatusParalyzed Status = "PAR"
	StatusHealthy   Status = "HEALTHY"
)

type StatusService interface {
	Paralyze(pokemon *Pokemon) error
}

type Status string

func (s Status) LogValue() slog.Value {
	return slog.StringValue(string(s))
}

func NewStatusService() StatusService {
	return &statusService{}
}

type statusService struct{}

func (s *statusService) Paralyze(pokemon *Pokemon) error {
	pokemon.Status = StatusParalyzed
	return nil
}
