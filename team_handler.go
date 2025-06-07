package main

import (
	"encoding/json"
	"net/http"

	"github.com/jootd/soccer-manager/business"
)

type UpdateTeamRequest struct {
	Name    string
	Country string
}

type TeamHandler struct {
	teamBus *business.TeamBus
}

func NewTeamHandler(teamBus *business.TeamBus) *TeamHandler {
	return &TeamHandler{
		teamBus: teamBus,
	}
}

func (tb *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	team := r.Context().Value(UserTeamContextKey)
	teamBytes, err := json.Marshal(team)
	if err != nil {
		http.Error(w, "something went wrong, please try again", http.StatusInternalServerError)
		return
	}
	w.Write(teamBytes)
}

func (tb *TeamHandler) UpdateTeam(w http.ResponseWriter, r *http.Request) {
	team := r.Context().Value(UserTeamContextKey)
	team = team.(business.Team)

	var updates UpdateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {

	}
	tb.teamBus.UpdateTeam(r.Context(), business.UpdateTeam{})

}
