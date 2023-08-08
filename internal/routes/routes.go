package routes

import (
	"database/sql"
	"github.com/gorilla/mux"
	"internlink/internal/handlers"
)

func SetupRoutes(r *mux.Router, db *sql.DB) {
	r.HandleFunc("/register", handlers.RegisterUserHandler(db)).Methods("POST")
	// Outras rotas
}
