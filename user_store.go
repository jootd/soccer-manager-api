package main

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/jootd/soccer-manager/business"
)

type dbUser struct {
	Username string
	Password string
	TeamId   int
}

type UserStore struct {
	mem   map[string]dbUser
	mutex sync.RWMutex
}

func toUser(db dbUser) business.User {
	return business.User{
		Username: db.Username,
		Password: db.Password,
		TeamId:   db.TeamId,
	}
}

func fromUser(user business.User) dbUser {
	return dbUser{
		Username: user.Username,
		Password: user.Password,
		TeamId:   user.TeamId,
	}
}

func (us *UserStore) GetUser(ctx context.Context, username string) (business.User, error) {
	us.mutex.RLock()
	defer us.mutex.RUnlock()

	dbUser, exists := us.mem[username]
	if !exists {
		return business.User{}, errors.New("store:GetUser:not_found")
	}

	return toUser(dbUser), nil
}

func (us *UserStore) CreateUser(ctx context.Context, username string, passHash string) (business.User, error) {
	us.mutex.Lock()
	defer us.mutex.Unlock()
	_, exists := us.mem[username]
	if exists {
		return business.User{}, errors.New("store:CreateUser:duplicate_username")
	}

	newDbUser := dbUser{
		Username: username,
		Password: passHash,
	}
	us.mem[username] = newDbUser

	return toUser(newDbUser), nil
}

func (us *UserStore) UpdateUser(ctx context.Context, username string, teamId int) (business.User, error) {
	user, err := us.GetUser(ctx, username)
	if err != nil {
		return business.User{}, fmt.Errorf("store:UpdateUser:%w", err)
	}

	us.mutex.Lock()
	defer us.mutex.Unlock()
	user.TeamId = teamId
	us.mem[username] = fromUser(user)

	return user, nil

}
