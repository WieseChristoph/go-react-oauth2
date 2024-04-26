package main

import (
	"github.com/WieseChristoph/go-oauth2-backend/internal/server"
	"github.com/WieseChristoph/go-oauth2-backend/internal/utils/log"
)

func main() {
	server := server.NewServer()

	log.Infof("Server started on %s", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Error starting server. Err: %v", err)
	}
}
