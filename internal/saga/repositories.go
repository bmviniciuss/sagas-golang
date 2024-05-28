package saga

import "context"

type ExecutionRepository interface {
	Insert(ctx context.Context, execution *Execution) error
	Find(ctx context.Context, globalID string) (*Execution, error)
	Save(ctx context.Context, execution *Execution) error
}

type WorkflowRepository interface {
	Find(ctx context.Context, name string) (*Workflow, error)
}
