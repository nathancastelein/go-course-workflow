package orchestration

import (
	"log/slog"
)

type Orchestrator struct {
	workflows WorkflowDefinitions
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		workflows: make(WorkflowDefinitions, 0),
	}
}

func (w *Orchestrator) Register(definer WorkflowDefiner) {
	w.workflows = append(w.workflows, definer.Definition())
}

func (w *Orchestrator) RunWorkflow(workflowExecution *WorkflowExecution) error {
	workflowToExecute, err := w.workflows.FindByName(workflowExecution.WorkflowName)
	if err != nil {
		return err
	}

	slog.Info("executing workflow", slog.Any("workflow_definition", workflowToExecute), slog.Any("workflow_execution", workflowExecution))
	alreadyPerformedSteps := workflowExecution.PerformedSteps

	for _, step := range workflowToExecute.Steps {
		if len(alreadyPerformedSteps) > 0 && alreadyPerformedSteps[0] == step.Name {
			alreadyPerformedSteps = alreadyPerformedSteps[1:]
			continue
		}

		slog.Info("performing step", slog.Any("step", step), slog.Any("trainer", workflowExecution.Trainer), slog.Any("pokemon", workflowExecution.Pokemon))
		err := step.Do(workflowExecution.Trainer, workflowExecution.Pokemon)
		if err != nil {
			return err
		}
		workflowExecution.PerformedSteps = append(workflowExecution.PerformedSteps, step.Name)
	}
	slog.Info("workflow executed", slog.Any("trainer", workflowExecution.Trainer), slog.Any("pokemon", workflowExecution.Pokemon))

	return nil
}
