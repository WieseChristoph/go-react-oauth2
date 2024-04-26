package server

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/robfig/cron/v3"

	"github.com/WieseChristoph/go-oauth2-backend/internal/api"
	"github.com/WieseChristoph/go-oauth2-backend/internal/auth"
	"github.com/WieseChristoph/go-oauth2-backend/internal/database"
	"github.com/WieseChristoph/go-oauth2-backend/internal/database/repositories"
	"github.com/WieseChristoph/go-oauth2-backend/internal/middleware"
	"github.com/WieseChristoph/go-oauth2-backend/internal/utils/env"
	"github.com/WieseChristoph/go-oauth2-backend/internal/utils/log"
)

type Server struct {
	port int

	db         *database.DB
	auth       auth.Auth
	middleware middleware.Middleware
	api        api.Api
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(env.GetEnvOrDefault("PORT", "80"))

	db := database.New()
	repos := repositories.New(db)
	NewServer := &Server{
		port: port,

		db:         db,
		auth:       auth.New(repos),
		middleware: middleware.New(repos),
		api:        api.New(repos),
	}

	// Migrate Database
	err := db.Migrate()
	if err != nil {
		log.Fatalf("Error migrating database. Err: %v", err)
	}

	// Register cron jobs
	c := cron.New()
	c.AddFunc("@daily", func() {
		log.Infoln("Runing cron job: DeleteExpiredSessions")
		err := repos.Session.DeleteExpiredSessions()
		if err != nil {
			log.Errorf("Error deleting expired sessions from database. Err: %v", err)
		}
	})
	c.Start()

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
