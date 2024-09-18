package orchestration

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindByName(t *testing.T) {
	// Arrange
	workflowDefinitions := WorkflowDefinitions{
		&WorkflowDefinition{
			Name: "TestWorkflow",
		},
	}

	// Act
	found, err := workflowDefinitions.FindByName("TestWorkflow")

	// Assert
	require.NoError(t, err)
	require.Equal(t, workflowDefinitions[0], found)
}

func TestFindByNameNotFound(t *testing.T) {
	// Arrange
	workflowDefinitions := WorkflowDefinitions{
		&WorkflowDefinition{
			Name: "TestWorkflow",
		},
	}

	// Act
	found, err := workflowDefinitions.FindByName("InvalidName")

	// Assert
	require.Error(t, err)
	require.Nil(t, found)
	require.ErrorIs(t, err, ErrWorkflowNotFound)
}
