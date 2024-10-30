package service

import (
	"github.com/adharshmk96/goutils/auth"
	"go.opentelemetry.io/otel/trace"
	"gostarter/infra"
	"log/slog"

	"gostarter/internals/core/utils"
	"gostarter/internals/domain"
)

type accountService struct {
	logger *slog.Logger
	tracer *trace.Tracer

	accountRepo domain.AccountRepository
}

func NewAccountService(container *infra.Container, accountRepo domain.AccountRepository) domain.AccountService {
	return &accountService{
		logger:      container.Logger,
		tracer:      container.Tracer,
		accountRepo: accountRepo,
	}
}

func (a *accountService) Register(account *domain.Account) error {
	passwdHash, err := auth.HashPassword(account.Password, auth.DefaultParams)
	if err != nil {
		return err
	}

	account.Password = passwdHash

	return a.accountRepo.CreateAccount(account)
}

func (a *accountService) Authenticate(email, password string) (*domain.Account, error) {
	account, err := a.accountRepo.GetAccountByEmail(email)
	if err != nil {
		return nil, err
	}

	if account == nil {
		return nil, domain.ErrAccountNotFound
	}

	match, err := auth.VerifyPasswordHash(password, account.Password)
	if err != nil {
		return nil, err
	}

	if !match {
		return nil, domain.ErrAccountNotFound
	}

	return account, nil
}

func (a *accountService) GetAccountByID(id int) (*domain.Account, error) {
	return a.accountRepo.GetAccountByID(id)
}

func (a *accountService) GetAccountByEmail(email string) (*domain.Account, error) {
	return a.accountRepo.GetAccountByEmail(email)
}

func (a *accountService) GetAccountByUsername(username string) (*domain.Account, error) {
	return a.accountRepo.GetAccountByUsername(username)
}

func (a *accountService) UpdateAccount(account *domain.Account) error {
	return a.accountRepo.UpdateAccount(account)
}

func (a *accountService) DeleteAccount(id int) error {
	return a.accountRepo.DeleteAccount(id)
}

func (a *accountService) ListAccounts(paginationParams utils.PaginationParams) ([]*domain.Account, error) {
	return a.accountRepo.ListAccounts(paginationParams)
}
