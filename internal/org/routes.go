package org

import "github.com/go-chi/chi/v5"

func UserRoutes(r chi.Router, handler *UserHandler) {
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

func TeamRoutes(r chi.Router, handler *TeamHandler) {
	r.Route("/teams", func(r chi.Router) {
		r.Get("/", handler.List)
		r.Post("/", handler.Create)

		r.Route("/{teamId}", func(r chi.Router) {
			r.Get("/", handler.GetByID)
			r.Put("/", handler.UpdateByID)
			r.Delete("/", handler.DeleteByID)

			r.Route("/members", func(r chi.Router) {
				r.Get("/", handler.ListTeamMembers)
				r.Post("/", handler.AddMember)
				r.Delete("/", handler.RemoveMember)
			})
		})
	})
}
