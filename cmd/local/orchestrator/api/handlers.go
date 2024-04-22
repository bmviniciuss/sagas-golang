package api

import (
	"net/http"
	"time"

	"github.com/bmviniciuss/sagas-golang/cmd/local/orchestrator/appcontext"
	"github.com/bmviniciuss/sagas-golang/internal/saga"
	"github.com/bmviniciuss/sagas-golang/internal/saga/service"
	"github.com/bmviniciuss/sagas-golang/pkg/responses"
	"github.com/go-chi/render"
	goval "github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"

	"go.uber.org/zap"
)

type HandlersPort interface {
	Health(w http.ResponseWriter, r *http.Request)
	CreateOrder(w http.ResponseWriter, r *http.Request)
}

type Handlers struct {
	logger             *zap.SugaredLogger
	workflowRepository saga.WorkflowRepository
	workflowService    service.Port
	validator          *goval.Validate
}

var (
	_ HandlersPort = (*Handlers)(nil)
)

func NewHandlers(
	logger *zap.SugaredLogger,
	workflowRepository saga.WorkflowRepository,
	workflowService service.Port,
	validator *goval.Validate,
) *Handlers {
	return &Handlers{
		logger:             logger,
		workflowRepository: workflowRepository,
		workflowService:    workflowService,
		validator:          validator,
	}
}

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, map[string]string{"status": "ok", "time": time.Now().Format(time.RFC3339)})
}

type CreateOrderRequest struct {
	CustomerID   string `json:"customer_id" mapstructure:"customer_id" validate:"required,uuid"`
	Date         string `json:"date" mapstructure:"date" validate:"required"`
	Total        *int64 `json:"total" mapstructure:"total" validate:"required,gt=0"`
	CurrencyCode string `json:"currency_code" mapstructure:"currency_code" validate:"required"`
}

func (h *Handlers) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var (
		lggr     = h.logger
		ctx      = r.Context()
		reqID, _ = appcontext.RequestID(r.Context())
	)
	lggr = lggr.With("request_id", reqID)
	lggr.Info("Creating order")

	var req CreateOrderRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		lggr.With(zap.Error(err)).Error("Got error decoding request")
		errRes := responses.ParseErrorToResponse(reqID, err)
		responses.RenderError(w, r, errRes)
		return
	}

	err := h.validator.StructCtx(ctx, req)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error validating request")
		fieldErrs := responses.ValidatorErrorToFieldError(err)
		errRes := responses.NewBadRequestErrorResponse(reqID, fieldErrs)
		responses.RenderError(w, r, errRes)
		return
	}

	lggr.Infof("Received request: %+v", req)

	data := map[string]interface{}{}
	err = mapstructure.Decode(req, &data)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error decoding create order request to map")
		errRes := responses.NewInternalServerErrorResponse(reqID)
		responses.RenderError(w, r, errRes)
		return
	}

	createOrderWorkflow, err := h.workflowRepository.Find(ctx, "create-order-v1")
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error finding workflow")
		errRes := responses.NewInternalServerErrorResponse(reqID)
		responses.RenderError(w, r, errRes)
		return
	}

	if createOrderWorkflow.IsEmpty() {
		lggr.Error("workflow not found")
		errRes := responses.NewInternalServerErrorResponse(reqID)
		responses.RenderError(w, r, errRes)
		return
	}

	globalID, err := h.workflowService.Start(r.Context(), createOrderWorkflow, data)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error starting workflow")
		errRes := responses.NewInternalServerErrorResponse(reqID)
		responses.RenderError(w, r, errRes)
		return
	}

	lggr.Infof("Successfully started workflow with global ID: %s", globalID.String())
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, map[string]string{"id": globalID.String()})
}
