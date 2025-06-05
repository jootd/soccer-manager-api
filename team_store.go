package main

import (
	"context"
	"crypto/rand"
	"sync"

	"github.com/jootd/soccer-manager/business"
)

type dbTeam struct {
	ID      int
	Name    string
	Country string
}

type TeamStore struct {
	mu  sync.RWMutex
	mem map[int]dbTeam
}

func toTeam(db dbTeam) business.Team {
	return business.Team{
		ID:      db.ID,
		Name:    db.Name,
		Country: db.Country,
	}
}

func toTeamSlice(dbTeams []dbTeam) []business.Team {
	teams := []business.Team{}
	for _, db := range dbTeams {
		teams = append(teams, toTeam(db))
	}
	return teams
}

func (dt *TeamStore) GetTeamsBy(ctx context.Context, query business.QueryTeam) ([]business.Team, bool) {

	if query.ID != nil {
		db, ok := dt.mem[*query.ID]
		if !ok {
			return []business.Team{}, false
		}

		return toTeamSlice([]dbTeam{db}), true
	}

	// if no filter applied return all
	allDb := []dbTeam{}
	for _, db := range dt.mem {
		allDb = append(allDb, db)
	}

	return toTeamSlice(allDb), true
}

func (dt *TeamStore) UpdateTeam(ctx context.Context, team business.UpdateTeam) (business.Team, bool) {
	dbTeam := dbTeam{}

	return toTeam(dbTeam), true
}
func (dt *TeamStore) CreateTeam(ctx context.Context) (business.Team, bool) {
	dbTeam := dbTeam{
		ID:      len(dt.mem) + 1,
		Name:    rand.Text(),
		Country: rand.Text()[:2],
	}
	return toTeam(dbTeam), true
}
