package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/henriquesd/go-gateway-api/internal/service"
	"github.com/henriquesd/go-gateway-api/internal/web/handlers"
	"github.com/henriquesd/go-gateway-api/internal/web/middleware"
)

type Server struct {
	router         *chi.Mux
	server         *http.Server
	accountService *service.AccountService
	invoiceService *service.InvoiceService
	port           string
}

func NewServer(
	accountService *service.AccountService,
	invoiceService *service.InvoiceService,
	port string) *Server {
	router := chi.NewRouter()

	return &Server{
		router:         router,
		accountService: accountService,
		invoiceService: invoiceService,
		port:           port,
	}
}

func (server *Server) ConfigureRoutes() {
	accountHandler := handlers.NewAccountHandler(server.accountService)
	invoiceHandler := handlers.NewInvoiceHandler(server.invoiceService)

	authMiddleware := middleware.NewAuthMiddleware(server.accountService)

	server.router.Post("/accounts", accountHandler.Create)
	server.router.Get("/accounts", accountHandler.Get)

	server.router.Group(func(router chi.Router) {
		router.Use(authMiddleware.Authenticate)
		server.router.Post("/invoice", invoiceHandler.Create)
		server.router.Get("/invoice/{id}", invoiceHandler.GetByID)
		server.router.Get("/invoice", invoiceHandler.ListByAccount)
	})
}

func (server *Server) Start() error {
	server.server = &http.Server{
		Addr:    ":" + server.port,
		Handler: server.router,
	}

	return server.server.ListenAndServe()
}
