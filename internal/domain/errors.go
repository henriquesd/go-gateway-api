package domain

import "errors"

var (
	ErrorAccountNotFound    = errors.New("account not found")
	ErrorDuplicatedAPIKey   = errors.New("api key already exists")
	ErrorInvoiceNotFound    = errors.New("invoice not found")
	ErrorUnauthorizedAccess = errors.New("unauthorized access")
	ErrorInvalidAmount      = errors.New("invalid amount")
	ErrorInvalidStatus      = errors.New("invalid status")
)
