package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bmviniciuss/sagas-golang/cmd/local/orchestrator/appcontext"
	"github.com/bmviniciuss/sagas-golang/internal/saga"
	"github.com/bmviniciuss/sagas-golang/internal/saga/service"
	"github.com/bmviniciuss/sagas-golang/pkg/responses"
	"github.com/go-chi/render"
	goval "github.com/go-playground/validator/v10"

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

func StructToMap(data interface{}) (map[string]interface{}, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

type CreateOrderRequest struct {
	CustomerID   string `json:"customer_id" validate:"required,uuid"`
	Card         string `json:"card" validate:"required"`
	Amount       *int64 `json:"amount" validate:"required,gt=0"`
	CurrencyCode string `json:"currency_code" validate:"required"`
	Items        []Item `json:"items" validate:"required,min=1,dive"`
}

type Item struct {
	ID        string `json:"id" validate:"required,uuid"`
	Quantity  *int32 `json:"quantity" validate:"required,gt=0"`
	UnitPrice *int64 `json:"unit_price" validate:"required,gt=0"`
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

	createOrderWorkflow, err := h.workflowRepository.Find(ctx, "create_order_v1")
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error finding workflow")
		errRes := responses.NewInternalServerErrorResponse(reqID)
		responses.RenderError(w, r, errRes)
		return
	}

	data, err := StructToMap(req)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error decoding create order request to map")
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
