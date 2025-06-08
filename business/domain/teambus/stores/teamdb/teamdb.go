package teamdb

import (
	"bytes"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/jootd/soccer-manager/business/sdk/sqldb"
	"go.uber.org/zap"
)

type Store struct {
	log *zap.SugaredLogger
	db  sqlx.ExtContext
}

func NewStore(logger *zap.SugaredLogger) *Store {
	return &Store{
		log: logger,
	}
}

func (s *Store) NewWithTx(tx sqldb.CommitRollbacker) (teambus.Storer, error) {
	ec, err := sqldb.GetExtContext(tx)
	if err != nil {
		return nil, err
	}

	store := Store{
		log: s.log,
		db:  ec,
	}

	return &store, nil
}

func (s *Store) Query(ctx context.Context, filter teambus.QueryFilter) ([]teambus.Team, error) {
	data := make(map[string]any)
	const q = `
	SELECT
	    id, name, country, budget, date_created, date_updated
	FROM
		teams`

	buf := bytes.NewBufferString(q)
	s.applyFilter(filter, data, buf)

	return []teambus.Team{}, nil
}

func (s *Store) Update(ctx context.Context, updates teambus.Team) error {
	const q = `
	UPDATE 
		teams
	SET
		"name" = :name
		"country" = :country
		"budget" = :budget
	WHERE
		id = :id
	`
	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBTeam(updates)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) Create(ctx context.Context, new teambus.Team) error {
	const q = `
	INSERT INTO users
		(name, country, budget)
	VALUES
		(:name, :country, :budget)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBTeam(new)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (dt *Store) all() []teambus.Team {
	return []teambus.Team{}
}
