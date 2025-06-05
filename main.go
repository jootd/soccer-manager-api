package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jootd/soccer-manager/business"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var jwtKey = []byte("test_secret_key")

func SignupHandler(userBus *business.UserBus, teamBus *business.TeamBus) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		_, err = userBus.CreateUser(r.Context(), req.Username, string(hash))
		if err != nil {
			http.Error(w, "user already exists", http.StatusConflict)
			return
		}

		team, err := teamBus.CreateTeam(r.Context())
		if err != nil {
			http.Error(w, "something went wrong, please try again", http.StatusInternalServerError)
			return
		}

		user, err := userBus.UpdateUser(r.Context(), req.Username, team.ID)
		if err != nil {
			http.Error(w, "something went wrong, please try again", http.StatusInternalServerError)
			return
		}

		w.Write([]byte(fmt.Sprintf("registered, user table %+v", user)))

	}
}

func LoginHandler(userBus *business.UserBus) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		user, err := userBus.GetUser(r.Context(), req.Username)
		if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
			http.Error(w, "username or password incorrect", http.StatusUnauthorized)
			return
		}

		token, err := generateJWT(req.Username)
		if err != nil {
			http.Error(w, "something went wrong, please try again", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Authorization", "Bearer "+token)
		w.Write([]byte(fmt.Sprintf("Logged in. Token: %s, userStore %+v", token, user)))
	}
}

func TeamHandler(teamBus *business.TeamBus) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func main() {
	userStore := UserStore{
		mem: make(map[string]dbUser),
	}
	teamStore := TeamStore{
		mem: make(map[int]dbTeam),
	}

	userBus := business.NewUserBus(&userStore)
	teamBus := business.NewTeamBus(&teamStore)

	jwtMiddleware := CreateJWTMiddleware(userBus, teamBus)

	_ = jwtMiddleware

	http.HandleFunc("/signup", SignupHandler(userBus, teamBus))
	http.HandleFunc("/login", LoginHandler(userBus))

	fmt.Println("Server started at 8080")
	http.ListenAndServe(":8080", nil)

}
