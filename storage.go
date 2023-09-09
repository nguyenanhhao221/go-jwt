package main

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/nguyenanhhao221/go-jwt/util"
)

type Storage interface {
	createAccountTable() error
	GetAllAccounts() ([]Account, error)
	CreateAccount(*Account) (uuid.UUID, error)
	GetAccountById(accountId uuid.UUID) (*AccountResponse, error)
	GetAccountByUsername(username string) (*Account, error)
	DeleteAccountById(accountId uuid.UUID) error
	UpdateAccountById(updateAccount *Account, accountId uuid.UUID) error
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

func (s *PostgresStore) GetAllAccounts() ([]Account, error) {
	query := `
	SELECT * from ACCOUNT 
	`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allAccounts []Account
	for rows.Next() {
		var account Account
		if err := rows.Scan(
			&account.ID,
			&account.FirstName,
			&account.LastName,
			&account.Number,
			&account.Balance,
			&account.CreatedAt,
		); err != nil {
			return nil, err
		}
		allAccounts = append(allAccounts, account)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return allAccounts, err
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
	username VARCHAR(255) NOT NULL,
	password BYTEA NOT NULL,
	number INTEGER,
	balance INTEGER,
	created_at TIMESTAMP
	)`
	_, err := s.db.Exec(query)
	return err
}

// AccountResponse Use this if we don't want to include the username and password
type AccountResponse struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Number    int64     `json:"number"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
}

func (s *PostgresStore) GetAccountById(accountId uuid.UUID) (*AccountResponse, error) {
	query := `
	SELECT id, first_name, last_name, number, balance, created_at
	FROM account
	WHERE id = $1 
	`
	var account AccountResponse
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

func (s *PostgresStore) GetAccountByUsername(username string) (*Account, error) {
	query := `
	SELECT *
	FROM account
	WHERE username = $1 
	`
	var account Account
	row := s.db.QueryRow(query, username)
	err := row.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Username,
		&account.Password,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
	)
	if err != nil {
		return &account, err
	}
	return &account, nil
}

func (s *PostgresStore) DeleteAccountById(accountId uuid.UUID) error {
	query := `
	DELETE 
	FROM account
	WHERE id = $1 
	`
	if _, err := s.db.Exec(query, accountId); err != nil {
		return err
	} else {
		return nil
	}
}

// CreateAccount Create account in the database, also handle hashing the password
func (s *PostgresStore) CreateAccount(newAccount *Account) (uuid.UUID, error) {
	query := `
	INSERT INTO ACCOUNT (first_name, last_name, number, balance, created_at, username, password)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING ID
	`
	hashPassword, hashPasswordErr := util.HashPassword(newAccount.Password)
	if hashPasswordErr != nil {
		return uuid.Nil, hashPasswordErr
	}
	var id uuid.UUID
	err := s.db.QueryRow(
		query,
		newAccount.FirstName,
		newAccount.LastName,
		newAccount.Number,
		newAccount.Balance,
		newAccount.CreatedAt,
		newAccount.Username, hashPassword).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (s *PostgresStore) UpdateAccountById(updateAccount *Account, accountId uuid.UUID) error {
	query := `
	UPDATE ACCOUNT 	
	SET first_name = $2, last_name = $3, number = $4, balance = $5
	WHERE id = $1
	`
	_, err := s.db.Exec(
		query,
		accountId,
		updateAccount.FirstName,
		updateAccount.LastName,
		updateAccount.Number,
		updateAccount.Balance,
	)
	return err
}
