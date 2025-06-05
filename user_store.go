package main

import (
	"context"
	"sync"

	"github.com/jootd/soccer-manager/business"
)

type dbUser struct {
	Username string
	Password string
	TeamId   int
}

type UserStore struct {
	mem map[string]dbUser
	mu  sync.RWMutex
}

func toUser(db dbUser) business.User {
	return business.User{
		Username: db.Username,
		Password: db.Password,
		TeamId:   db.TeamId,
	}
}

func (us *UserStore) GetUser(ctx context.Context, username string) (business.User, bool) {
	us.mu.RLock()
	defer us.mu.RUnlock()

	dbUser, exists := us.mem[username]
	if !exists {
		return business.User{}, false
	}

	return toUser(dbUser), true
}

func (us *UserStore) CreateUser(ctx context.Context, username string, passHash string) (business.User, bool) {
	us.mu.Lock()
	defer us.mu.Unlock()

	newUser := dbUser{Username: username, Password: passHash}
	us.mem[username] = newUser

	_, exists := us.mem[username]

	if !exists {
		return business.User{}, false
	}

	return toUser(newUser), true
}

func (us *UserStore) UpdateUser(ctx context.Context, username string, teamId int) (business.User, bool) {
	us.mu.Lock()
	defer us.mu.Unlock()

	user, exists := us.mem[username]

	if !exists {
		return business.User{}, false
	}

	user.TeamId = teamId

	us.mem[username] = user

	return toUser(user), true

}
