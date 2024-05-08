package api

import (
	"net/http"
	"time"

	"github.com/bmviniciuss/sagas-golang/cmd/local/orchestrator/appcontext"
	"github.com/bmviniciuss/sagas-golang/cmd/local/order/application/usecases"
	"github.com/bmviniciuss/sagas-golang/cmd/local/order/presentation"
	"github.com/bmviniciuss/sagas-golang/pkg/responses"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type OrderHandlers interface {
	Health(w http.ResponseWriter, r *http.Request)
	ListAll(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
}

type Handlers struct {
	lggr           *zap.SugaredLogger
	listUseCase    *usecases.ListOrders
	getByIDUseCase *usecases.GetOrderByID
}

var (
	_ OrderHandlers = (*Handlers)(nil)
)

func NewHandlers(lggr *zap.SugaredLogger, listUseCase *usecases.ListOrders, getByIDUseCase *usecases.GetOrderByID) *Handlers {
	return &Handlers{
		lggr:           lggr,
		listUseCase:    listUseCase,
		getByIDUseCase: getByIDUseCase,
	}
}

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, map[string]string{"status": "ok", "time": time.Now().UTC().Format(time.RFC3339)})
}

func (h *Handlers) ListAll(w http.ResponseWriter, r *http.Request) {
	h.lggr.Info("Listing all orders")
	results, err := h.listUseCase.Execute(r.Context())
	if err != nil {
		h.lggr.With(zap.Error(err)).Error("Got error listing orders")
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    "internal_error",
				"message": "Got error listing orders",
			},
		})
		return
	}
	var res presentation.OrderList
	res.Content = results.Orders
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, res)
}

func (h *Handlers) GetByID(w http.ResponseWriter, r *http.Request) {
	var (
		lggr     = h.lggr
		ctx      = r.Context()
		reqID, _ = appcontext.RequestID(ctx)
	)
	lggr.Info("Getting order by ID")
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error decoding order id")
		errRes := responses.NewBadRequestErrorResponse(reqID, []responses.FieldError{
			{
				Field:   "id",
				Message: "Invalid order ID",
			},
		})
		responses.RenderError(w, r, errRes)
		return
	}
	order, err := h.getByIDUseCase.Execute(ctx, id)
	if err != nil {
		h.lggr.With(zap.Error(err)).Error("Got error finding order by ID")
		errRes := responses.NewInternalServerErrorResponse(reqID)
		responses.RenderError(w, r, errRes)
		return
	}
	if order.IsEmpty() {
		h.lggr.With(zap.Error(err)).Error("Order not found")
		errRes := responses.NewNotFoundErrorResponse(reqID)
		responses.RenderError(w, r, errRes)
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, order)
}
