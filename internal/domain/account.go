package domain

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID        string
	Name      string
	Email     string
	ApiKey    string
	Balance   float64
	mutex     sync.RWMutex // It allows multiple readers to access the data simultaneously but ensures that only one writer can modify the data at a time (using lock). This avoids race conditions.
	CreatedAt time.Time
	UpdatedAt time.Time
}

func generateApiKey() string {
	keyBytes := make([]byte, 16)
	rand.Read(keyBytes)
	return hex.EncodeToString(keyBytes)
}

func NewAccount(name, email string) *Account {
	account := &Account{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		ApiKey:    generateApiKey(),
		Balance:   0.0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return account
}

func (account *Account) AddBalance(amount float64) {
	account.mutex.Lock()

	// Using defer, it will be unlocked at the end of the function, even if an error occurs.
	defer account.mutex.Unlock()

	account.Balance += amount
	account.UpdatedAt = time.Now()
}
