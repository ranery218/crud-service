package http

import (
	"crud/internal/transport/http/middleware"
	"net/http"

	"github.com/go-chi/chi"
)

func NewRouter(userHandler *UserHandler, authMiddleware *middleware.AuthMiddleware) http.Handler {
	r := chi.NewRouter()
	r.Route("/users", func(r chi.Router) {
		r.Post("/register", userHandler.Register)
		r.Post("/login", userHandler.Login)
	})
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.RequireAuth)
		r.Post("/users/logout", userHandler.Logout)
		r.Patch("/users/me", userHandler.Update)
	})
	return r
}
