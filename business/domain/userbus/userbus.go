package userbus

import (
	"context"
	"errors"
	"fmt"

	"github.com/jootd/soccer-manager/business/sdk/sqldb"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNotFound       = errors.New("user_not_found")
	ErrUniqueUsername = errors.New("user_already_exists")
)

type Storer interface {
	NewWithTx(tx sqldb.CommitRollbacker) (Storer, error)
	Get(ctx context.Context, username string) (User, error)
	Create(ctx context.Context, new User) error
	Update(ctx context.Context, update User) error
}

type ExtBusiness interface {
	NewWithTx(tx sqldb.CommitRollbacker) (ExtBusiness, error)
	Get(ctx context.Context, username string) (User, error)
	Create(ctx context.Context, new CreateUser) error
	Update(ctx context.Context, upd UpdateUser) error
}

type Extension func(ExtBusiness) ExtBusiness

type Business struct {
	store Storer
	log   *zap.SugaredLogger
}

func NewUserBus(store Storer, log *zap.SugaredLogger, extensions ...Extension) ExtBusiness {
	b := ExtBusiness(&Business{
		store: store,
		log:   log,
	})
	for i := len(extensions) - 1; i >= 0; i-- {
		ext := extensions[i]
		if ext != nil {
			b = ext(b)
		}
	}
	return b
}

func (b *Business) NewWithTx(tx sqldb.CommitRollbacker) (ExtBusiness, error) {
	storer, err := b.store.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	bus := Business{
		log:   b.log,
		store: storer,
	}

	return &bus, nil
}

func (ub *Business) Get(ctx context.Context, username string) (User, error) {
	user, err := ub.store.Get(ctx, username)
	if err != nil {
		return User{}, fmt.Errorf("bus:Get:%w", err)
	}

	return user, nil
}

func (ub *Business) Create(ctx context.Context, new CreateUser) error {
	if len(new.Username) == 0 || len(new.Password) == 0 {
		return fmt.Errorf("username and password cannot be empty ")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(new.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hashing failed")
	}

	user := User{
		Username: new.Username,
		Password: string(hash),
		TeamID:   new.TeamID,
	}
	err = ub.store.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("bus:Create:%w", err)
	}

	return nil
}
func (ub *Business) Update(ctx context.Context, upd UpdateUser) error {

	user := User{}

	if upd.Username != nil {
		user.Username = *upd.Username
	}

	if upd.PasswordHash != nil {
		user.Password = *upd.PasswordHash
	}

	if upd.TeamID != nil {
		user.TeamID = *upd.TeamID
	}

	err := ub.store.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("bus:Update:%w", err)
	}

	return nil
}
