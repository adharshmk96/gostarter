package routing

import (
	"github.com/go-chi/chi/v5"
	"gostarter/internals/delivery/http/api"
	"gostarter/internals/service"
	"gostarter/internals/storage/memory"
	"net/http"
)

func SetupRoutes(routing chi.Router) {
	routing.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	accountRepo := memory.NewAccountRepository()
	accountService := service.NewAccountService(accountRepo)
	accountHandler := api.NewAccountHandler(accountService)

	routing.Post("/v1/auth/register", accountHandler.Register)
}
