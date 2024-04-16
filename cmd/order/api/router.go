package api

import "github.com/go-chi/chi/v5"

type Router struct {
	handlers OrderHandlers
}

func NewRouter(handlers OrderHandlers) *Router {
	return &Router{
		handlers: handlers,
	}
}

func (rr *Router) Build() *chi.Mux {
	router := chi.NewRouter()
	router.Route("/v1", func(r chi.Router) {
		r.Get("/health", rr.handlers.Health)
		r.Get("/orders", rr.handlers.ListAll)
	})
	return router
}
