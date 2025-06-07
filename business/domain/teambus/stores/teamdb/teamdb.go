package teamdb

import (
	"context"
	"log"

	"github.com/jootd/soccer-manager/business/domain/teambus"
)

type Store struct {
	//sqlx

	log *log.Logger
}

func NewStore(logger *log.Logger) *Store {
	return &Store{
		log: logger,
	}
}

func (s *Store) Query(ctx context.Context, query teambus.QueryFilter) ([]teambus.Team, error) {

	return []teambus.Team{}, nil
}

func (s *Store) UpdateTeam(ctx context.Context, updates teambus.UpdateTeam) (teambus.Team, error) {

	return teambus.Team{}, nil
}

func (s *Store) Create(ctx context.Context, new teambus.CreateTeam) (teambus.Team, error) {

	return teambus.Team{}, nil
}

func (dt *Store) all() []teambus.Team {
	return []teambus.Team{}
}
