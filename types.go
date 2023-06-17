package main

import (
	"github.com/google/uuid"
)

type Account struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Number    int64     `json:"number"`
	Balance   int64     `json:"balance"`
}

func NewAccount(firstName, lastName string) *Account {
	id := uuid.New()
	return &Account{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Number:    0,
		Balance:   0,
	}
}
