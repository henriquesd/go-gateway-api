package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/henriquesd/go-gateway-api/internal/domain"
	"github.com/henriquesd/go-gateway-api/internal/dto"
	"github.com/henriquesd/go-gateway-api/internal/service"
)

type InvoiceHandler struct {
	service *service.InvoiceService
}

func NewInvoiceHandler(service *service.InvoiceService) *InvoiceHandler {
	return &InvoiceHandler{
		service: service,
	}
}

// Request requires authentication with X-API-KEY header
// Endpoint: /invoice
// Method: POST
func (invoiceHandler *InvoiceHandler) Create(writer http.ResponseWriter, request *http.Request) {
	var input dto.CreateInvoiceInput
	err := json.NewDecoder(request.Body).Decode(&input)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	input.APIKey = request.Header.Get("X-API-KEY")

	output, err := invoiceHandler.service.Create(input)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(output)
}

// Endpoint: /invoice/{id}
// Method: GET
func (invoiceHandler *InvoiceHandler) GetByID(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")
	if id == "" {
		http.Error(writer, "ID is required", http.StatusBadRequest)
		return
	}

	apiKey := request.Header.Get("X-API-KEY")
	if apiKey == "" {
		http.Error(writer, "X-API-KEY is required", http.StatusBadRequest)
		return
	}

	output, err := invoiceHandler.service.GetByID(id, apiKey)
	if err != nil {
		switch err {
		case domain.ErrorInvoiceNotFound:
			http.Error(writer, err.Error(), http.StatusNotFound)
			return
		case domain.ErrorAccountNotFound:
			http.Error(writer, err.Error(), http.StatusUnauthorized)
			return
		case domain.ErrorUnauthorizedAccess:
			http.Error(writer, err.Error(), http.StatusForbidden)
			return
		default:
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(output)
}

// Endpoint: /invoice
// Method: GET
func (invoiceHandler *InvoiceHandler) ListByAccount(writer http.ResponseWriter, request *http.Request) {
	apiKey := request.Header.Get("X-API-KEY")
	if apiKey == "" {
		http.Error(writer, "X-API-KEY is required", http.StatusUnauthorized)
		return
	}

	output, err := invoiceHandler.service.ListByAccountAPIKey(apiKey)
	if err != nil {
		switch err {
		case domain.ErrorAccountNotFound:
			http.Error(writer, err.Error(), http.StatusUnauthorized)
			return
		default:
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(output)
}
