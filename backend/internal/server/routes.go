package server

import (
	"context"
	"net/http"

	"github.com/WieseChristoph/go-oauth2-backend/internal/auth"
	"github.com/WieseChristoph/go-oauth2-backend/internal/config"
	"github.com/WieseChristoph/go-oauth2-backend/internal/database/models"
	"github.com/WieseChristoph/go-oauth2-backend/internal/utils/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)

	r.Get("/", s.GetIndexHandler)

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(s.middleware.AuthMiddleware)
		r.Route("/users", func(r chi.Router) {
			r.Get("/me", s.api.GetUsersMe)
			r.Group(func(r chi.Router) {
				r.Use(s.middleware.RoleMiddleware(models.RoleAdmin))
				r.Get("/", s.api.GetUsersAll)
			})
		})
	})

	r.Route("/auth", func(r chi.Router) {
		r.Get("/{provider}", s.getAuthBeginHandler)
		r.Get("/{provider}/callback", s.getAuthCallbackHandler)
		r.With(s.middleware.AuthMiddleware).Get("/logout", s.getLogoutHandler)
	})

	return r
}

func (s *Server) GetIndexHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello, World!"))
	if err != nil {
		log.Errorf("Error writing response. Err: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (s *Server) getAuthBeginHandler(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(context.Background(), config.ProviderContextKey, provider))

	err := s.auth.BeginOAuth2Flow(w, r)
	if err != nil {
		if err == auth.ErrAlreadyLoggedIn {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		log.Errorf("Error beginning oauth flow for provider '%s'. Err: %v", provider, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (s *Server) getAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(context.Background(), config.ProviderContextKey, provider))

	err := s.auth.HandleOAuth2Callback(w, r)
	if err != nil {
		log.Errorf("Error handling oauth callback for provider '%s'. Err: %v", provider, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (s *Server) getLogoutHandler(w http.ResponseWriter, r *http.Request) {
	err := s.auth.Logout(w, r)
	if err != nil {
		log.Errorf("Error logging out. Err: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
