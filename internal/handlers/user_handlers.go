package handlers

import (
	"database/sql"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"internlink/internal/models"
	"log"
	"net/http"
	"time"
)

func RegisterUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newUser models.User
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&newUser)
		if err != nil {
			http.Error(w, "Invalid request data", http.StatusBadRequest)
			return
		}

		emailExistsQuery := "SELECT COUNT(*) FROM users WHERE email = $1"
		var emailCount int
		err = db.QueryRow(emailExistsQuery, newUser.Email).Scan(&emailCount)
		if err != nil {
			log.Println("Failed to check email:", err)
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
			return
		}
		if emailCount > 0 {
			http.Error(w, "This email is already in use", http.StatusBadRequest)
			return
		}

		usernameExistsQuery := "SELECT COUNT(*) FROM users WHERE username = $1"
		var usernameCount int
		err = db.QueryRow(usernameExistsQuery, newUser.Username).Scan(&usernameCount)
		if err != nil {
			log.Println("Failed to check username:", err)
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
			return
		}
		if usernameCount > 0 {
			http.Error(w, "This username is already in use", http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password),
			bcrypt.DefaultCost)
		if err != nil {
			log.Println("Failed to hash password:", err)
			http.Error(w, "Failed to register user:", http.StatusInternalServerError)
		}

		insertQuery := `
            INSERT INTO users (username, email, password, created_at, admin)
            VALUES ($1, $2, $3, $4, $5)
        `
		_, err = db.Exec(insertQuery, newUser.Username, newUser.Email, hashedPassword, time.Now(),
			newUser.Admin)
		if err != nil {
			log.Println("Failed to insert user:", err)
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User registered successfully"))
	}
}
