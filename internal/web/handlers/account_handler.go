package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/henriquesd/go-gateway-api/internal/dto"
	"github.com/henriquesd/go-gateway-api/internal/service"
)

type AccountHandler struct {
	accountService *service.AccountService
}

func NewAccountHandler(accountService *service.AccountService) *AccountHandler {
	return &AccountHandler{accountService: accountService}
}

func (accountHandler *AccountHandler) Create(writer http.ResponseWriter, request *http.Request) {
	var input dto.CreateAccountInput

	err := json.NewDecoder(request.Body).Decode(&input)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	output, err := accountHandler.accountService.CreateAccount(input)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(output)
}

func (accountHandler *AccountHandler) Get(writer http.ResponseWriter, request *http.Request) {
	apiKey := request.Header.Get("X-API-Key")

	if apiKey == "" {
		http.Error(writer, "X-API-Key header is required", http.StatusUnauthorized)
		return
	}

	output, err := accountHandler.accountService.FindByAPIKey(apiKey)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(output)
}
