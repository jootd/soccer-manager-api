package teamdb

import (
	"bytes"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/jootd/soccer-manager/business/domain/teambus"
	"github.com/jootd/soccer-manager/business/sdk/sqldb"
	"go.uber.org/zap"
)

type Store struct {
	log *zap.SugaredLogger
	db  sqlx.ExtContext
}

func NewStore(logger *zap.SugaredLogger, db *sqlx.DB) teambus.Storer {
	return &Store{
		log: logger,
		db:  db,
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

func (s *Store) GetByID(ctx context.Context, id int) (teambus.Team, error) {

	data := struct {
		ID int `db:"id"`
	}{ID: id}
	const q = `
	SELECT
		id, name, country, budget
	FROM
		teams
	WHERE
		"id" = :id
	`

	var result team
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &result); err != nil {
		return teambus.Team{}, fmt.Errorf("NamedQueryStruct: %w", err)
	}

	return toBusTeam(result), nil
}

func (s *Store) Query(ctx context.Context, filter teambus.QueryFilter) ([]teambus.Team, error) {
	//TODO:
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
		"name" = COALESCE(:name, "name"),
    	"country" = COALESCE(:country, "country"),
    	"budget" = COALESCE(:budget, "budget")
	WHERE
		id = :id
	`
	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBTeam(updates)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) Create(ctx context.Context, new teambus.Team) (int, error) {
	const q = `
	INSERT INTO teams 
		(name, country, budget)
	VALUES
		(:name, :country, :budget)
	RETURNING 
		id`

	result := struct {
		ID int `db:"id"`
	}{}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, toDBTeam(new), &result); err != nil {
		return 0, fmt.Errorf("namedexeccontext: %w", err)
	}

	return result.ID, nil
}

func (dt *Store) all() []teambus.Team {
	return []teambus.Team{}
}
