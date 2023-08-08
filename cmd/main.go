package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"internlink/internal/routes"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	connData := "user=postgres password=funyarinpa999 dbname=internlink sslmode=disable"
	db, err := sql.Open("postgres", connData)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Database connection test failed:", err)
	}

	fmt.Println("Connected to the database!")

	routes.SetupRoutes(r, db)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've requested the InternLink API!")
	})

	http.ListenAndServe(":8080", r)
}
