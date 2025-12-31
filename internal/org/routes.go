package org

import "github.com/go-chi/chi/v5"

func Routes(r chi.Router, handler *UserHandler) {
	r.Route("/users", func(r chi.Router) {
		r.Get("/", handler.List)
		r.Post("/", handler.Create)

		r.Route("/{userId}", func(r chi.Router) {
			r.Get("/", handler.GetByID)
			r.Put("/", handler.UpdateByID)
			r.Delete("/", handler.DeleteByID)
		})
	})
}
