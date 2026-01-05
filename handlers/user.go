package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/tarunmonga/hello-world/db"
	"golang.org/x/crypto/bcrypt"
)

type UserResponse struct {
	ID         string    `json:"id"`
	Username   string    `json:"username"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	ProfilePic string    `json:"profilePic,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Conn.Query(context.Background(), "SELECT id, username, name, email, COALESCE(profilePic, ''), created_at, updated_at FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []UserResponse
	for rows.Next() {
		var u UserResponse
		err := rows.Scan(&u.ID, &u.Username, &u.Name, &u.Email, &u.ProfilePic, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	type RegisterRequest struct {
		Username string `json:"username"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic Validation
	if len(req.Username) < 3 {
		http.Error(w, "Username must be at least 3 characters", http.StatusBadRequest)
		return
	}
	if len(req.Password) < 6 {
		http.Error(w, "Password must be at least 6 characters", http.StatusBadRequest)
		return
	}
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(req.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	// Insert into DB
	var id string
	err := db.Conn.QueryRow(context.Background(),
		"INSERT INTO users (username, name, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		req.Username, req.Name, req.Email, string(hashedPassword), time.Now(), time.Now(),
	).Scan(&id)

	if err != nil {
		fmt.Printf("CreateUser DB Error: %v\n", err) // Log to terminal
		http.Error(w, fmt.Sprintf("Failed to create user: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func GetMe(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	var u UserResponse
	err := db.Conn.QueryRow(context.Background(),
		"SELECT id, username, name, email, COALESCE(profilePic, ''), created_at, updated_at FROM users WHERE id=$1",
		userID,
	).Scan(&u.ID, &u.Username, &u.Name, &u.Email, &u.ProfilePic, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	userIDFromToken := r.Context().Value("userID").(string)
	targetID := chi.URLParam(r, "id")

	// Authorization
	if userIDFromToken != targetID {
		http.Error(w, "Unauthorized to update this user", http.StatusForbidden)
		return
	}

	type UpdateRequest struct {
		Username   string `json:"username"`
		Name       string `json:"name"`
		Email      string `json:"email"`
		Password   string `json:"password"` // Optional
		ProfilePic string `json:"profilePic"`
	}
	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Password != "" {
		if len(req.Password) < 6 {
			http.Error(w, "Password must be at least 6 characters", http.StatusBadRequest)
			return
		}
		hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		_, err := db.Conn.Exec(context.Background(), "UPDATE users SET password=$1, updated_at=$2 WHERE id=$3", string(hashed), time.Now(), targetID)
		if err != nil {
			http.Error(w, "Update failed", http.StatusInternalServerError)
			return
		}
	}

	_, err := db.Conn.Exec(context.Background(),
		"UPDATE users SET username=$1, name=$2, email=$3, profilePic=$4, updated_at=$5 WHERE id=$6",
		req.Username, req.Name, req.Email, req.ProfilePic, time.Now(), targetID,
	)

	if err != nil {
		fmt.Printf("UpdateUser DB Error: %v\n", err) // Log to terminal
		http.Error(w, "Update failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User updated successfully"))
}
