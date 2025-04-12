package middleware

import (
	"net/http"

	"github.com/henriquesd/go-gateway-api/internal/domain"
	"github.com/henriquesd/go-gateway-api/internal/service"
)

type AuthMiddleware struct {
	accountService *service.AccountService
}

func NewAuthMiddleware(accountService *service.AccountService) *AuthMiddleware {
	return &AuthMiddleware{
		accountService: accountService,
	}
}

func (authMiddleware *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		apiKey := request.Header.Get("X-API-KEY")
		if apiKey == "" {
			http.Error(writer, "X-API-KEY is required", http.StatusUnauthorized)
			return
		}

		_, err := authMiddleware.accountService.FindByAPIKey(apiKey)
		if err != nil {
			if err == domain.ErrorAccountNotFound {
				http.Error(writer, err.Error(), http.StatusUnauthorized)
				return
			}

			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(writer, request)
	})
}
