package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/jootd/soccer-manager/business/data/dbschema"
	"github.com/jootd/soccer-manager/business/sdk/sqldb"
)

// Seed loads test data into the database.
func Seed(cfg sqldb.Config) error {
	db, err := sqldb.Open(cfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := dbschema.Seed(ctx, db); err != nil {
		return fmt.Errorf("seed database: %w", err)
	}

	fmt.Println("seed data complete")
	return nil
}
