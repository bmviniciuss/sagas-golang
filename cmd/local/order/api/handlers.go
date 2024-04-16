package api

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type OrderHandlers interface {
	Health(w http.ResponseWriter, r *http.Request)
	ListAll(w http.ResponseWriter, r *http.Request)
}

type Handlers struct {
	lggr *zap.SugaredLogger
}

var (
	_ OrderHandlers = (*Handlers)(nil)
)

func NewHandlers(lggr *zap.SugaredLogger) *Handlers {
	return &Handlers{
		lggr: lggr,
	}
}

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, map[string]string{"status": "ok", "time": time.Now().UTC().Format(time.RFC3339)})
}

func (h *Handlers) ListAll(w http.ResponseWriter, r *http.Request) {
	h.lggr.Info("Listing all orders")
	w.WriteHeader(http.StatusOK)
}
