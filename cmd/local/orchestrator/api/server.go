package api

import (
	"net/http"

	"github.com/bmviniciuss/sagas-golang/cmd/local/orchestrator/appcontext"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Router struct {
	handlers HandlersPort
}

func NewRouter(handlers HandlersPort) *Router {
	return &Router{
		handlers: handlers,
	}
}

func (r *Router) Build() *chi.Mux {
	router := chi.NewRouter()
	router.Use(requestIDMiddleware)
	router.Get("/v1/health", r.handlers.Health)
	router.Post("/v1/create-orders", r.handlers.CreateOrder)
	return router
}

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		hReqID := r.Header.Get("X-Request-ID")
		if hReqID != "" {
			requestID = hReqID
		}
		ctx := appcontext.WithRequestID(r.Context(), requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
