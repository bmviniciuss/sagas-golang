package executions

import (
	"context"
	"sync"

	"github.com/bmviniciuss/sagas-golang/internal/saga"
)

type InmemRepository struct {
	mu   *sync.Mutex
	data map[string]*saga.Execution
}

func NewInmemRepository() *InmemRepository {
	return &InmemRepository{
		mu:   &sync.Mutex{},
		data: make(map[string]*saga.Execution),
	}
}

func (r *InmemRepository) Find(ctx context.Context, globalID string) (*saga.Execution, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if execution, ok := r.data[globalID]; ok {
		return execution, nil
	}
	return &saga.Execution{}, nil
}

func (r *InmemRepository) Save(ctx context.Context, execution *saga.Execution) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[execution.ID.String()] = execution
	return nil
}
