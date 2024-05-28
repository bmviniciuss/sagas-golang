package executions

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/bmviniciuss/sagas-golang/cmd/local/orchestrator/adapters/repositories/executions/generated"
	"github.com/bmviniciuss/sagas-golang/internal/saga"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type RepositoryAdapter struct {
	lggr               *zap.SugaredLogger
	pool               *pgxpool.Pool
	workflowRepository saga.WorkflowRepository
}

func NewRepositoryAdapter(lggr *zap.SugaredLogger, pool *pgxpool.Pool, workflowRepository saga.WorkflowRepository) *RepositoryAdapter {
	return &RepositoryAdapter{
		lggr:               lggr,
		pool:               pool,
		workflowRepository: workflowRepository,
	}
}

var (
	_ saga.ExecutionRepository = (*RepositoryAdapter)(nil)
)

func (r *RepositoryAdapter) Insert(ctx context.Context, execution *saga.Execution) error {
	lggr := r.lggr
	lggr.Info("RepositoryAdapter.Insert")

	db, err := r.pool.Acquire(ctx)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error acquiring connection")
		return err
	}
	defer db.Release()
	queries := generated.New(db)

	state, err := json.Marshal(execution.State)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error Marshalling state")
		return err
	}

	err = queries.InsertExecution(ctx, generated.InsertExecutionParams{
		Uuid:         execution.ID,
		WorkflowName: execution.Workflow.Name,
		State:        state,
	})
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error inserting workflow execution")
		return nil
	}

	return nil
}

func (r *RepositoryAdapter) Save(ctx context.Context, execution *saga.Execution) error {
	lggr := r.lggr
	lggr.Info("RepositoryAdapter.Save")

	db, err := r.pool.Acquire(ctx)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error acquiring connection")
		return err
	}
	defer db.Release()
	queries := generated.New(db)

	state, err := json.Marshal(execution.State)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error Marshalling state")
		return err
	}

	err = queries.UpdateExecution(ctx, generated.UpdateExecutionParams{
		Uuid:  execution.ID,
		State: state,
	})
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error updating workflow execution")
		return nil
	}

	return nil
}

func (r *RepositoryAdapter) Find(ctx context.Context, globalID string) (*saga.Execution, error) {
	lggr := r.lggr
	lggr.Info("RepositoryAdapter.Find")

	gid, err := uuid.Parse(globalID)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error parsing global id to uuid")
		return nil, err
	}

	db, err := r.pool.Acquire(ctx)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error acquiring connection")
		return nil, err
	}
	defer db.Release()
	queries := generated.New(db)

	execRow, err := queries.FindExecutionByUUID(ctx, gid)
	if err == pgx.ErrNoRows {
		return &saga.Execution{}, nil
	}
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error finding workflow execution")
		return nil, err
	}

	var state map[string]interface{}
	err = json.Unmarshal(execRow.State, &state)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error Unmarshalling state")
		return nil, err
	}

	wflw, err := r.workflowRepository.Find(ctx, execRow.WorkflowName)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error finding workflow")
		return nil, err
	}
	if wflw.IsEmpty() {
		return nil, errors.New("workflow not found")
	}

	return &saga.Execution{
		ID:       execRow.Uuid,
		Workflow: wflw,
		State:    state,
	}, nil
}
