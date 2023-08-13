package main

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Number    int64     `json:"number"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewAccount(firstName, lastName string) *Account {
	id := uuid.New()
	return &Account{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Number:    int64(rand.Intn(1000000)),
		Balance:   0,
		CreatedAt: time.Now().UTC(),
	}
}

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
