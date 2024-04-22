package workflows

import (
	"context"
	"sync"

	"github.com/bmviniciuss/sagas-golang/internal/saga"
)

type InmemRepository struct {
	mu   *sync.Mutex
	data map[string]*saga.Workflow
}

func NewInmemRepository(workflows []saga.Workflow) *InmemRepository {
	data := make(map[string]*saga.Workflow)
	for _, w := range workflows {
		data[w.Name] = &w
	}

	return &InmemRepository{
		mu:   &sync.Mutex{},
		data: data,
	}
}

func (r *InmemRepository) Find(ctx context.Context, name string) (*saga.Workflow, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if execution, ok := r.data[name]; ok {
		return execution, nil
	}
	return &saga.Workflow{}, nil
}
