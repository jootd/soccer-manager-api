package playerdb

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/jootd/soccer-manager/business/domain/playerbus"
	"github.com/jootd/soccer-manager/business/sdk/sqldb"
	"go.uber.org/zap"
)

type Store struct {
	log *zap.SugaredLogger
	db  sqlx.ExtContext
}

func NewStore(logger *zap.SugaredLogger, db *sqlx.DB) playerbus.Storer {
	return &Store{
		log: logger,
		db:  db,
	}
}

func (s *Store) NewWithTx(tx sqldb.CommitRollbacker) (playerbus.Storer, error) {
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

func (s *Store) GetByTeamID(ctx context.Context, teamID int) ([]playerbus.Player, error) {

	data := struct {
		TeamID int `db:"team_id"`
	}{
		TeamID: teamID,
	}
	const q = `
	SELECT
		id, team_id, first_name, last_name, age, country, value,  position	
	FROM
		players
	WHERE
		"team_id" = :team_id
	`

	var result []player
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, q, data, &result); err != nil {
		return []playerbus.Player{}, fmt.Errorf("NamedQuerySlice: %w", err)
	}

	return toPlayerBusSlice(result), nil
}

func (s *Store) Query(ctx context.Context, filter playerbus.QueryFilter) ([]playerbus.Player, error) {
	data := make(map[string]any)
	const q = `
	SELECT
	    id, team_id, first_name, last_name, age, country, value,  position
	FROM
		players`

	result := []player{}
	buf := bytes.NewBufferString(q)
	s.applyFilter(filter, data, buf)

	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &result); err != nil {
		return []playerbus.Player{}, fmt.Errorf("NamedQuerySlice: %w", err)
	}

	return toPlayerBusSlice(result), nil
}

func (s *Store) Update(ctx context.Context, updates playerbus.Player) error {
	const q = `
	UPDATE 
		teams
	SET
		team_id    = COALESCE(:team_id, team_id),
		first_name = COALESCE(:first_name, first_name),
		last_name  = COALESCE(:last_name, last_name),
		country    = COALESCE(:country, country),
		value      = COALESCE(:value, value),
		position   = COALESCE(:position, position)
	WHERE
		id = :id
	`
	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBPlayer(updates)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) Create(ctx context.Context, new playerbus.Player) (int, error) {
	const q = `
	INSERT INTO teams 
		(team_id, first_name, last_name, country, value, position)
	VALUES
		(:team_id, :first_name, :last_name, :country, :value, :position)
	RETURNING 
		id`

	result := struct {
		ID int `db:"id"`
	}{}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, toDBPlayer(new), &result); err != nil {
		return 0, fmt.Errorf("namedexeccontext: %w", err)
	}

	return result.ID, nil
}

func (s *Store) CreateBatch(ctx context.Context, players []playerbus.Player) error {
	if len(players) == 0 {
		return nil
	}

	playersDB := toDBPlayerSlice(players)

	const baseQuery = `
	INSERT INTO players (team_id, first_name, last_name, age, country, value, position)
	VALUES `

	valueStrings := make([]string, 0, len(players))
	valueArgs := make([]interface{}, 0, len(players)*7)

	for i, p := range playersDB {
		offset := i * 7
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			offset+1, offset+2, offset+3, offset+4, offset+5, offset+6, offset+7))
		valueArgs = append(valueArgs,
			p.TeamID,
			p.FirstName,
			p.LastName,
			p.Age,
			p.Country,
			p.Value,
			p.Position,
		)
	}

	query := baseQuery + strings.Join(valueStrings, ",")

	_, err := s.db.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return fmt.Errorf("batch insert players: %w", err)
	}

	return nil
}

func (s *Store) All(ctx context.Context) ([]playerbus.Player, error) {
	var result []player
	const q = `SELECT  * FROM players`

	if err := sqldb.QuerySlice(ctx, s.log, s.db, q, &result); err != nil {
		return []playerbus.Player{}, fmt.Errorf("namedexeccontext: %w", err)
	}

	return toPlayerBusSlice(result), nil
}
