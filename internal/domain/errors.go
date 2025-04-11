package domain

import "errors"

var (
	ErrorAccountNotFound    = errors.New("account not found")
	ErrorDuplicatedApiKey   = errors.New("api key already exists")
	ErrorInvoiceNotFound    = errors.New("invoice not found")
	ErrorUnauthorizedAccess = errors.New("unauthorized access")
)
