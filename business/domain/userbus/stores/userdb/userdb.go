package userdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/jootd/soccer-manager/business/domain/userbus"
	"github.com/jootd/soccer-manager/business/sdk/sqldb"
	"go.uber.org/zap"
)

type Store struct {
	log *zap.SugaredLogger
	db  sqlx.ExtContext
}

func NewStore(log *zap.SugaredLogger, db *sqlx.DB) *Store {
	return &Store{
		log: log,
		db:  db,
	}
}

// value with a sqlx DB value that is currently inside a transaction.
func (s *Store) NewWithTx(tx sqldb.CommitRollbacker) (userbus.Storer, error) {
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

func (s *Store) Get(ctx context.Context, username string) (userbus.User, error) {
	data := struct {
		Username string `db:"username"`
	}{
		Username: username,
	}
	const q = `
	SELECT 
	    username, password_hash , team_id, date_created, date_updated
	FROM 
		users
	WHERE
		username = :username
	`
	var dbUser user
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbUser); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return userbus.User{}, fmt.Errorf("db: %w", userbus.ErrNotFound)
		}
		return userbus.User{}, fmt.Errorf("db: %w", err)
	}

	return toBusUser(dbUser), nil
}

func (s *Store) Create(ctx context.Context, new userbus.User) error {
	const q = `
	INSERT INTO users
		(username, password_hash, date_created, date_updated)
	VALUES
		(:username, :password_hash, :date_created, :date_updated)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBUser(new)); err != nil {
		if errors.Is(err, sqldb.ErrDBDuplicatedEntry) {
			return fmt.Errorf("namedexeccontext: %w", userbus.ErrUniqueUsername)
		}
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) Update(ctx context.Context, user userbus.User) error {
	const q = `
	UPDATE
		users
	SET 
		username       = COALESCE(:username, username),
		team_id        = COALESCE(:team_id, team_id),
		date_updated   = COALESCE(:date_updated, date_updated),
		date_created   = COALESCE(:date_created, date_created)
	WHERE
		username = :username
	`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBUser(user)); err != nil {
		if errors.Is(err, sqldb.ErrDBDuplicatedEntry) {
			return userbus.ErrUniqueUsername
		}
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}
