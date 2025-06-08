package transferdb

import (
	"context"
	"log"
)

type Store struct {
	log *log.Logger
}

func (s *Store) NewStore(logger *log.Logger) *Store {
	return &Store{
		log: logger,
	}
}
func (s *Store) Create(ctx context.Context, t transferbus.Transfer) error {

	return nil
}
func (s *Store) Query(ctx context.Context, query transferbus.QueryFilter) ([]transferbus.Transfer, error) {

	return []transferbus.Transfer{}, nil
}
func (s *Store) Update(ctx context.Context, update transferbus.UpdateTransfer) (transferbus.Transfer, error) {

	return transferbus.Transfer{}, nil
}
