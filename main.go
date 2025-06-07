package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jootd/soccer-manager/business"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	log.SetOutput(os.Stdout)

	userStore := UserStore{
		mem: make(map[string]dbUser),
	}
	teamStore := TeamStore{
		mem: make(map[int]dbTeam),
	}

	userBus := business.NewUserBus(&userStore)
	teamBus := business.NewTeamBus(&teamStore)

	authHandler := NewAuthHandler(userBus, teamBus)
	authMiddleware := CreateAuthMiddleware(userBus, teamBus)
	mux := mux.NewRouter()

	mux.HandleFunc("/signup", authHandler.Signup)
	mux.HandleFunc("/login", authHandler.Login)

	teamHandler := NewTeamHandler(teamBus)

	mux.Handle("/team", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			teamHandler.GetTeam(w, r)
		case http.MethodPut:
			teamHandler.UpdateTeam(w, r)
		}

	})))

	handler := LoggingMiddleware(mux)

	fmt.Println("Server started at 8080")
	http.ListenAndServe(":8080", handler)

}
