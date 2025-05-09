package service

import (
	"github.com/henriquesd/go-gateway-api/internal/domain"
	"github.com/henriquesd/go-gateway-api/internal/dto"
)

type AccountService struct {
	repository domain.AccountRepository
}

func NewAccountService(repository domain.AccountRepository) *AccountService {
	return &AccountService{repository: repository}
}

func (service *AccountService) CreateAccount(input dto.CreateAccountInput) (*dto.AccountOutput, error) {
	account := dto.ToAccount(input)

	existingAccount, err := service.repository.FindByAPIKey(account.APIKey)
	if err != nil && err != domain.ErrorAccountNotFound {
		return nil, err
	}
	if existingAccount != nil {
		return nil, domain.ErrorDuplicatedAPIKey
	}

	err = service.repository.Save(account)
	if err != nil {
		return nil, err
	}

	output := dto.FromAccount(account)
	return &output, nil
}

func (service *AccountService) UpdateBalance(apiKey string, amount float64) (*dto.AccountOutput, error) {
	account, err := service.repository.FindByAPIKey(apiKey)
	if err != nil {
		return nil, err
	}

	account.AddBalance(amount)

	err = service.repository.UpdateBalance(account)
	if err != nil {
		return nil, err
	}

	output := dto.FromAccount(account)
	return &output, nil
}

func (service *AccountService) FindByAPIKey(apiKey string) (*dto.AccountOutput, error) {
	account, err := service.repository.FindByAPIKey(apiKey)
	if err != nil {
		return nil, err
	}

	output := dto.FromAccount(account)
	return &output, nil
}

func (service *AccountService) FindByID(id string) (*dto.AccountOutput, error) {
	account, err := service.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	output := dto.FromAccount(account)
	return &output, nil
}
