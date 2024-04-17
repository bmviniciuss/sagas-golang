package api

import "github.com/go-chi/chi/v5"

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
	router.Get("/v1/health", r.handlers.Health)
	router.Post("/v1/create-orders", r.handlers.CreateOrder)
	return router
}
