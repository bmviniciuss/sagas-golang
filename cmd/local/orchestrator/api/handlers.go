package api

import (
	"net/http"
	"time"

	"github.com/bmviniciuss/sagas-golang/internal/saga"
	"github.com/bmviniciuss/sagas-golang/internal/saga/service"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type HandlersPort interface {
	Health(w http.ResponseWriter, r *http.Request)
	CreateOrder(w http.ResponseWriter, r *http.Request)
}

type Handlers struct {
	logger              *zap.SugaredLogger
	createOrderWorkflow *saga.Workflow
	workflowService     service.Port
}

var (
	_ HandlersPort = (*Handlers)(nil)
)

func NewHandlers(logger *zap.SugaredLogger, createOrderWorkflow *saga.Workflow, workflowService service.Port) *Handlers {
	return &Handlers{
		logger:              logger,
		createOrderWorkflow: createOrderWorkflow,
		workflowService:     workflowService,
	}
}

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, map[string]string{"status": "ok", "time": time.Now().Format(time.RFC3339)})
}

func (h *Handlers) CreateOrder(w http.ResponseWriter, r *http.Request) {
	lggr := h.logger
	lggr.Info("Creating order")
	globalID, err := h.workflowService.Start(r.Context(), h.createOrderWorkflow, nil)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error starting workflow")
		render.Status(r, http.StatusInternalServerError)
	}
	lggr.Infof("Successfully started workflow with global ID: %s", globalID.String())

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, map[string]string{"id": globalID.String()})
}
