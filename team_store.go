package main

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/jootd/soccer-manager/business"
)

type dbTeam struct {
	ID      int
	Name    string
	Country string
}

type TeamStore struct {
	mutex sync.RWMutex
	mem   map[int]dbTeam
	idSeq int
}

func toTeam(db dbTeam) business.Team {
	return business.Team{
		ID:      db.ID,
		Name:    db.Name,
		Country: db.Country,
	}
}

func fromTeam(team business.Team) dbTeam {
	return dbTeam{
		ID:      team.ID,
		Name:    team.Name,
		Country: team.Country,
	}
}

func toTeamSlice(dbTeams []dbTeam) []business.Team {
	teams := []business.Team{}
	for _, db := range dbTeams {
		teams = append(teams, toTeam(db))
	}
	return teams
}

func (dt *TeamStore) all(ctx context.Context) []business.Team {
	allDb := []dbTeam{}
	for _, db := range dt.mem {
		allDb = append(allDb, db)
	}
	return toTeamSlice(allDb)
}

func (dt *TeamStore) GetTeamsBy(ctx context.Context, query business.QueryTeam) ([]business.Team, error) {
	dt.mutex.RLock()
	defer dt.mutex.RUnlock()
	if query.ID != nil {
		db, ok := dt.mem[*query.ID]
		if !ok {
			return []business.Team{}, errors.New("store:GetTeamsBy:not_found")
		}
		return toTeamSlice([]dbTeam{db}), nil
	}

	return dt.all(ctx), nil
}

func (dt *TeamStore) UpdateTeam(ctx context.Context, updates business.UpdateTeam) (business.Team, error) {
	result, err := dt.GetTeamsBy(ctx, business.QueryTeam{ID: &updates.Id})
	if err != nil {
		return business.Team{}, fmt.Errorf("store:UpdateTeam:%w", err)
	}

	needsUpdate := result[0]
	if updates.Country != nil {
		needsUpdate.Country = *updates.Country
	}
	if updates.Name != nil {
		needsUpdate.Name = *updates.Name
	}

	dt.mutex.Lock()
	defer dt.mutex.Unlock()
	dt.mem[needsUpdate.ID] = fromTeam(needsUpdate)

	return needsUpdate, nil
}

func (dt *TeamStore) CreateTeam(ctx context.Context, new business.CreateTeam) (business.Team, error) {
	dt.mutex.Lock()
	defer dt.mutex.Unlock()
	dt.idSeq++
	newId := dt.idSeq

	newTeam := dbTeam{
		ID:      newId,
		Name:    new.Name,
		Country: new.Country,
	}

	dt.mem[newId] = newTeam
	return toTeam(newTeam), nil
}
