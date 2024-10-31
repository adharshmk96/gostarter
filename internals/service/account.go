package service

import (
	"context"
	"github.com/adharshmk96/goutils/auth"
	"go.opentelemetry.io/otel/trace"
	"gostarter/infra"
	"gostarter/pkg/utils"
	"log/slog"

	"gostarter/internals/domain"
)

type accountService struct {
	logger *slog.Logger
	tracer trace.Tracer

	accountRepo domain.AccountRepository
}

func NewAccountService(container *infra.Container, accountRepo domain.AccountRepository) domain.AccountService {
	return &accountService{
		logger:      container.Logger,
		tracer:      container.Tracer,
		accountRepo: accountRepo,
	}
}

func (a *accountService) Register(ctx context.Context, account *domain.Account) error {
	ctx, span := a.tracer.Start(ctx, "AccountService.Register")
	defer span.End()

	passwdHash, err := auth.HashPassword(account.Password, auth.DefaultParams)
	if err != nil {
		return err
	}

	account.Password = passwdHash

	return a.accountRepo.CreateAccount(ctx, account)
}

func (a *accountService) Authenticate(ctx context.Context, email, password string) (*domain.Account, error) {
	ctx, span := a.tracer.Start(ctx, "AccountService.Authenticate")
	defer span.End()

	account, err := a.accountRepo.GetAccountByEmail(nil, email)
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

func (a *accountService) GetAccountByID(ctx context.Context, id int) (*domain.Account, error) {
	ctx, span := a.tracer.Start(ctx, "AccountService.GetAccountByID")
	defer span.End()

	return a.accountRepo.GetAccountByID(ctx, id)
}

func (a *accountService) GetAccountByEmail(ctx context.Context, email string) (*domain.Account, error) {
	ctx, span := a.tracer.Start(ctx, "AccountService.GetAccountByEmail")
	defer span.End()

	return a.accountRepo.GetAccountByEmail(ctx, email)
}

func (a *accountService) GetAccountByUsername(ctx context.Context, username string) (*domain.Account, error) {
	ctx, span := a.tracer.Start(ctx, "AccountService.GetAccountByUsername")
	defer span.End()

	return a.accountRepo.GetAccountByUsername(ctx, username)
}

func (a *accountService) UpdateAccount(ctx context.Context, account *domain.Account) error {
	ctx, span := a.tracer.Start(ctx, "AccountService.UpdateAccount")
	defer span.End()

	return a.accountRepo.UpdateAccount(ctx, account)
}

func (a *accountService) DeleteAccount(ctx context.Context, id int) error {
	ctx, span := a.tracer.Start(ctx, "AccountService.DeleteAccount")
	defer span.End()

	return a.accountRepo.DeleteAccount(ctx, id)
}

func (a *accountService) ListAccounts(ctx context.Context, paginationParams utils.PaginationParams) ([]*domain.Account, error) {
	ctx, span := a.tracer.Start(ctx, "AccountService.ListAccounts")
	defer span.End()

	return a.accountRepo.ListAccounts(ctx, paginationParams)
}
