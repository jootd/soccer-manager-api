package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jootd/soccer-manager/business/domain/teambus"
	"github.com/jootd/soccer-manager/business/domain/teambus/stores/teamdb"
	"github.com/jootd/soccer-manager/business/domain/userbus"
	"github.com/jootd/soccer-manager/business/domain/userbus/stores/userdb"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	log.SetOutput(os.Stdout)

	userStore := userdb.NewMemory()
	teamStore := teamdb.NewMemory()

	userBus := userbus.NewUserBus(userStore)
	teamBus := teambus.NewTeamBus(teamStore)

	authHandler := NewAuthHandler(userBus, teamBus)
	authMiddleware := CreateAuthMiddleware(userBus, teamBus)
	mux := mux.NewRouter()

	mux.HandleFunc("/signup", authHandler.Signup)
	mux.HandleFunc("/login", authHandler.Login)

	teamHandler := NewHandler(teamBus)

	mux.Handle("/team", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			teamHandler.Get(w, r)
		case http.MethodPut:
			teamHandler.Update(w, r)
		}

	})))

	handler := LoggingMiddleware(mux)

	fmt.Println("Server started at 8080")
	http.ListenAndServe(":8080", handler)

}
