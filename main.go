package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var (
	userStore = make(map[string]string)
	mu        sync.Mutex
	jwtKey    = []byte("test_secret_key")
)

func generateJWT(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(1 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "cannot read body", http.StatusInternalServerError)
		return
	}

	var req RegisterRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if len(req.Username) == 0 || len(req.Password) == 0 {
		http.Error(w, "username and password required", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if _, exists := userStore[req.Username]; exists {
		http.Error(w, "user already exists", http.StatusConflict)
		return
	}
	userStore[req.Username] = string(hash)

}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "cannot read body", http.StatusInternalServerError)
		return
	}

	var req RegisterRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	hash, exists := userStore[req.Username]
	if !exists || bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)) != nil {
		http.Error(w, "username or password incorrect", http.StatusUnauthorized)
		return
	}

	token, err := generateJWT(req.Username)
	if err != nil {
		http.Error(w, "something went wrong, please try again", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
	w.Write([]byte(fmt.Sprintf("Logged in. Token: %s, userStore %+v", token, userStore)))
}

func main() {
	http.HandleFunc("/register", Register)
	http.HandleFunc("/login", Login)

	fmt.Println("Server started at 8080")
	http.ListenAndServe(":8080", nil)

}
