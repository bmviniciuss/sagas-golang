package api

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type KitchenApiHandler interface {
	Health(w http.ResponseWriter, r *http.Request)
}

type Handlers struct {
	lggr *zap.SugaredLogger
}

var (
	_ KitchenApiHandler = (*Handlers)(nil)
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
