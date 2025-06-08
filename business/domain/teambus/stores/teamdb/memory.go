package teamdb

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type Memory struct {
	idSeq int
	mem   map[int]team
	mutex sync.RWMutex
}

func NewMemory() *Memory {
	return &Memory{
		mem: make(map[int]team),
	}

}

func (dt *Memory) Query(ctx context.Context, query teambus.QueryFilter) ([]teambus.Team, error) {
	dt.mutex.RLock()
	defer dt.mutex.RUnlock()
	if query.ID != nil {
		db, ok := dt.mem[*query.ID]
		if !ok {
			return []teambus.Team{}, errors.New("memory:GetTeamsBy:not_found")
		}
		return toBusTeamSlice([]team{db}), nil
	}

	return dt.all(), nil
}

func (dt *Memory) Update(ctx context.Context, updates teambus.UpdateTeam) (teambus.Team, error) {
	result, err := dt.Query(ctx, teambus.QueryFilter{ID: &updates.ID})
	if err != nil {
		return teambus.Team{}, fmt.Errorf("memory:UpdateTeam:%w", err)
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
	dt.mem[needsUpdate.ID] = toDBTeam(needsUpdate)

	return needsUpdate, nil
}

func (dt *Memory) Create(ctx context.Context, new teambus.CreateTeam) (teambus.Team, error) {
	dt.mutex.Lock()
	defer dt.mutex.Unlock()
	dt.idSeq++
	newId := dt.idSeq

	newTeam := team{
		ID:      newId,
		Name:    new.Name,
		Country: new.Country,
	}

	dt.mem[newId] = newTeam
	return toBusTeam(newTeam), nil
}

func (dt *Memory) all() []teambus.Team {
	allDb := []team{}
	for _, db := range dt.mem {
		allDb = append(allDb, db)
	}
	return toBusTeamSlice(allDb)
}
