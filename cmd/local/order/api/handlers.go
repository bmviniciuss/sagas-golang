package api

import (
	"net/http"
	"time"

	"github.com/bmviniciuss/sagas-golang/cmd/local/order/application/usecases"
	"github.com/bmviniciuss/sagas-golang/cmd/local/order/presentation"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type OrderHandlers interface {
	Health(w http.ResponseWriter, r *http.Request)
	ListAll(w http.ResponseWriter, r *http.Request)
}

type Handlers struct {
	lggr        *zap.SugaredLogger
	listUseCase *usecases.ListOrders
}

var (
	_ OrderHandlers = (*Handlers)(nil)
)

func NewHandlers(lggr *zap.SugaredLogger, listUseCase *usecases.ListOrders) *Handlers {
	return &Handlers{
		lggr:        lggr,
		listUseCase: listUseCase,
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
