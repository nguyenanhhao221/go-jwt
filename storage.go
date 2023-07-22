package main

import (
	"database/sql"
	"errors"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Storage interface {
	createAccountTable() error
	CreateAccount(*Account) (uuid.UUID, error)
	GetAccountById(accountId uuid.UUID) (*Account, error)
	// DeleteAccount(int) error
	// UpdateAccount(*Account) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error while loading env")
		return nil, err
	}

	// The URL in .env to connect to SQL database
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
		return nil, errors.New("DB_URL is not found in the environment")
	}

	/*  Open a connection to the postgres database */
	// sql.Open is used to establish a connection to the PostgreSQL database.
	// However, the sql.Open function only creates a connection object, it doesn't actually establish a connection to the database.
	// In order to use this we also need to import "github.com/lib/pq"
	sqlConnection, dbConnErr := sql.Open("postgres", dbURL)
	if dbConnErr != nil {
		log.Fatal("Cannot connect to the database: ", dbConnErr)
		return nil, dbConnErr
	}
	// sql.Open successfully returns an instance of sql.DB regardless of whether the database server is running or not.
	// To check if the connection was successful, you need to call the Ping method on the sql.DB instance.
	if err := sqlConnection.Ping(); err != nil {
		log.Fatal("Failed to ping the database, did you forget to run Docker? Error: ", err)
	}
	return &PostgresStore{
		db: sqlConnection,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.createAccountTable()
}

func (s *PostgresStore) createAccountTable() error {
	query := `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	CREATE TABLE IF NOT EXISTS ACCOUNT (
	id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	first_name VARCHAR(50),
	last_name VARCHAR(50),
	number SERIAL,
	balance SERIAL,
	created_at TIMESTAMP
	)`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) GetAccountById(accountId uuid.UUID) (*Account, error) {
	query := `
	SELECT *
	FROM account
	WHERE id = $1 
	`
	var account Account
	row := s.db.QueryRow(query, accountId)
	err := row.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
	)
	if err != nil {
		return &account, err
	}
	return &account, nil
}

func (s *PostgresStore) CreateAccount(newAccount *Account) (uuid.UUID, error) {
	query := `
	INSERT INTO ACCOUNT (first_name, last_name, number, balance, created_at)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING ID
	`

	var id uuid.UUID
	err := s.db.QueryRow(
		query,
		newAccount.FirstName,
		newAccount.LastName,
		newAccount.Number,
		newAccount.Balance,
		newAccount.CreatedAt).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
