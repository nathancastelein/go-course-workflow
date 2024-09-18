package orchestration

import (
	"errors"
	"log/slog"
	"strings"
)

var (
	ErrWorkflowNotFound = errors.New("workflow not found")
)

type WorkflowDefinition struct {
	Name  string
	Steps []*Step
}

type WorkflowDefiner interface {
	Definition() *WorkflowDefinition
}

type WorkflowDefinitions []*WorkflowDefinition

func (w WorkflowDefinitions) FindByName(workflowName string) (*WorkflowDefinition, error) {
	for _, workflowDefinition := range w {
		if workflowDefinition.Name == workflowName {
			return workflowDefinition, nil
		}
	}

	return nil, ErrWorkflowNotFound
}

func (w WorkflowDefinition) LogValue() slog.Value {
	steps := make([]string, len(w.Steps))
	for i, step := range w.Steps {
		steps[i] = step.Name
	}
	return slog.GroupValue(
		slog.String("name", w.Name),
		slog.String("steps", strings.Join(steps, ",")),
	)
}
