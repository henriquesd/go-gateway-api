package repository

import (
	"database/sql"
	"time"

	"github.com/henriquesd/go-gateway-api/internal/domain"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (repository *AccountRepository) Save(account *domain.Account) error {

	stmt, err := repository.db.Prepare(`
		INSERT INTO accounts (id, name, email, api_key, balance, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(
		account.ID,
		account.Name,
		account.Email,
		account.APIKey,
		account.Balance,
		account.CreatedAt,
		account.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (accountRepository *AccountRepository) FindByAPIKey(apiKey string) (*domain.Account, error) {
	var account domain.Account
	var createdAt, updatedAt time.Time

	err := accountRepository.db.QueryRow(`
		SELECT id, name, email, api_key, balance, created_at, updated_at
		FROM accounts
		WHERE api_key = $1
	`, apiKey).Scan(
		&account.ID,
		&account.Name,
		&account.Email,
		&account.APIKey,
		&account.Balance,
		&createdAt,
		&updatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrorAccountNotFound
	}
	if err != nil {
		return nil, err
	}

	account.CreatedAt = createdAt
	account.UpdatedAt = updatedAt
	return &account, nil
}

func (accountRepository *AccountRepository) FindByID(id string) (*domain.Account, error) {
	var account domain.Account
	var createdAt, updatedAt time.Time

	err := accountRepository.db.QueryRow(`
		SELECT id, name, email, api_key, balance, created_at, updated_at
		FROM accounts
		WHERE id = $1
	`, id).Scan(
		&account.ID,
		&account.Name,
		&account.Email,
		&account.APIKey,
		&account.Balance,
		&createdAt,
		&updatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrorAccountNotFound
	}
	if err != nil {
		return nil, err
	}

	account.CreatedAt = createdAt
	account.UpdatedAt = updatedAt
	return &account, nil
}

func (accountRepository *AccountRepository) UpdateBalance(account *domain.Account) error {
	transaction, err := accountRepository.db.Begin()
	if err != nil {
		return err
	}
	defer transaction.Rollback()

	var currentBalance float64
	err = transaction.QueryRow(`
		SELECT balance
		FROM accounts
		WHERE id = $1
		FOR UPDATE
		`, account.ID).Scan(&currentBalance)

	if err == sql.ErrNoRows {
		return domain.ErrorAccountNotFound
	}

	if err != nil {
		return err
	}

	_, err = transaction.Exec(`
		UPDATE accounts
		SET balance = $1, updated_at = $2
		WHERE id = $3
	`, currentBalance+account.Balance, time.Now(), account.ID)

	if err != nil {
		return err
	}

	return transaction.Commit()
}
