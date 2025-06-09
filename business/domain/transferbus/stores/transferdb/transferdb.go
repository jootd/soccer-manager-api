package transferdb

import (
	"bytes"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/jootd/soccer-manager/business/domain/transferbus"
	"github.com/jootd/soccer-manager/business/sdk/sqldb"
	"go.uber.org/zap"
)

type Store struct {
	log *zap.SugaredLogger
	db  sqlx.ExtContext
}

func NewStore(logger *zap.SugaredLogger, db *sqlx.DB) transferbus.Storer {
	return &Store{
		log: logger,
		db:  db,
	}
}

func (s *Store) NewWithTx(tx sqldb.CommitRollbacker) (transferbus.Storer, error) {
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

func (s *Store) GetByPlayerID(ctx context.Context, playerID int) ([]transferbus.Transfer, error) {
	data := struct {
		PlayerID int `db:"player_id"`
	}{
		PlayerID: playerID,
	}
	const q = `
	SELECT
		id, player_id, seller_id, asking_price, status
	FROM
		transfers
	WHERE
		player_id = :player_id
	`

	var result []transfer
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, q, data, &result); err != nil {
		return []transferbus.Transfer{}, fmt.Errorf("NamedQuerySlice: %w", err)
	}

	return toTransferBusSlice(result), nil
}

func (s *Store) Query(ctx context.Context, filter transferbus.QueryFilter) ([]transferbus.Transfer, error) {
	data := make(map[string]any)
	const q = `
	SELECT
		id, player_id, seller_id, asking_price, status
	FROM
		transfers`

	var result []transfer
	buf := bytes.NewBufferString(q)
	s.applyFilter(filter, data, buf)

	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &result); err != nil {
		return []transferbus.Transfer{}, fmt.Errorf("NamedQuerySlice: %w", err)
	}

	return toTransferBusSlice(result), nil
}

func (s *Store) Update(ctx context.Context, updates transferbus.Transfer) error {
	const q = `
	UPDATE 
		transfers
	SET
		player_id    = COALESCE(:player_id, player_id),
		seller_id    = COALESCE(:seller_id, seller_id),
		asking_price = COALESCE(:asking_price, asking_price),
		status       = COALESCE(:status, status)
	WHERE
		id = :id
	`
	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBTransfer(updates)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) Create(ctx context.Context, new transferbus.Transfer) (int, error) {
	const q = `
	INSERT INTO transfers
		(player_id, seller_id, asking_price, status)
	VALUES
		(:player_id, :seller_id, :asking_price, :status)
	RETURNING
		id`

	result := struct {
		ID int `db:"id"`
	}{}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, toDBTransfer(new), &result); err != nil {
		return 0, fmt.Errorf("namedexeccontext: %w", err)
	}

	return result.ID, nil
}

func (s *Store) All(ctx context.Context) ([]transferbus.Transfer, error) {
	var result []transfer
	const q = `SELECT * FROM transfers`

	if err := sqldb.QuerySlice(ctx, s.log, s.db, q, &result); err != nil {
		return []transferbus.Transfer{}, fmt.Errorf("queryslice: %w", err)
	}

	return toTransferBusSlice(result), nil
}
