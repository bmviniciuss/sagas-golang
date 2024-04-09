package saga

import (
	"fmt"

	"github.com/google/uuid"
)

type Workflow struct {
	ID           uuid.UUID
	Name         string
	Steps        *StepsList
	ReplyChannel string
}

func (w Workflow) EventTypes() map[string]uuid.UUID {
	ts := map[string]uuid.UUID{}
	steps := w.Steps.ToList()
	for _, step := range steps {
		ts[fmt.Sprintf("%s.%s.success", w.Name, step.Name)] = w.ID
		ts[fmt.Sprintf("%s.%s.failure", w.Name, step.Name)] = w.ID
	}
	return ts
}
