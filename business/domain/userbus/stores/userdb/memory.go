package userdb

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/jootd/soccer-manager/business/domain/userbus"
)

type Memory struct {
	mem   map[string]user
	mutex sync.RWMutex
}

func NewMemory() *Memory {
	return &Memory{
		mem: make(map[string]user),
	}
}

func (us *Memory) Get(ctx context.Context, username string) (userbus.User, error) {
	us.mutex.RLock()
	defer us.mutex.RUnlock()

	user, exists := us.mem[username]
	if !exists {
		return userbus.User{}, errors.New("store:GetUser:not_found")
	}

	return toBusUser(user), nil
}

func (us *Memory) Create(ctx context.Context, username string, passHash string) (userbus.User, error) {
	us.mutex.Lock()
	defer us.mutex.Unlock()
	_, exists := us.mem[username]
	if exists {
		return userbus.User{}, errors.New("store:CreateUser:duplicate_username")
	}

	newuser := user{
		Username: username,
		Password: passHash,
	}
	us.mem[username] = newuser

	return toBusUser(newuser), nil
}

func (us *Memory) Update(ctx context.Context, username string, teamId int) (userbus.User, error) {
	user, err := us.Get(ctx, username)
	if err != nil {
		return userbus.User{}, fmt.Errorf("store:UpdateUser:%w", err)
	}

	us.mutex.Lock()
	defer us.mutex.Unlock()
	user.TeamID = teamId
	us.mem[username] = toDBUser(user)

	return user, nil
}
