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
	logger := container.Logger.With("path", "accountService")
	return &accountService{
		logger:      logger,
		tracer:      container.Tracer,
		accountRepo: accountRepo,
	}
}

func (a *accountService) Register(ctx context.Context, account *domain.Account) error {
	ctx, stopSpan := utils.TraceSpan(ctx, a.tracer, "AccountService.Register")
	defer stopSpan()

	passwdHash, err := auth.HashPassword(account.Password, auth.DefaultParams)
	if err != nil {
		return err
	}

	account.Password = passwdHash

	return a.accountRepo.CreateAccount(ctx, account)
}

func (a *accountService) Authenticate(ctx context.Context, email, password string) (*domain.Account, error) {
	ctx, stopSpan := utils.TraceSpan(ctx, a.tracer, "AccountService.Authenticate")
	defer stopSpan()

	account, err := a.accountRepo.GetAccountByEmail(ctx, email)
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
	ctx, stopSpan := utils.TraceSpan(ctx, a.tracer, "AccountService.GetAccountByID")
	defer stopSpan()

	return a.accountRepo.GetAccountByID(ctx, id)
}

func (a *accountService) GetAccountByEmail(ctx context.Context, email string) (*domain.Account, error) {
	ctx, stopSpan := utils.TraceSpan(ctx, a.tracer, "AccountService.GetAccountByEmail")
	defer stopSpan()

	return a.accountRepo.GetAccountByEmail(ctx, email)
}

func (a *accountService) UpdateAccount(ctx context.Context, account *domain.Account) error {
	ctx, stopSpan := utils.TraceSpan(ctx, a.tracer, "AccountService.UpdateAccount")
	defer stopSpan()

	return a.accountRepo.UpdateAccount(ctx, account)
}

func (a *accountService) DeleteAccount(ctx context.Context, id int) error {
	ctx, stopSpan := utils.TraceSpan(ctx, a.tracer, "AccountService.DeleteAccount")
	defer stopSpan()

	return a.accountRepo.DeleteAccount(ctx, id)
}

func (a *accountService) ListAccounts(ctx context.Context, paginationParams utils.PaginationParams) ([]*domain.Account, error) {
	ctx, stopSpan := utils.TraceSpan(ctx, a.tracer, "AccountService.ListAccounts")
	defer stopSpan()

	return a.accountRepo.ListAccounts(ctx, paginationParams)
}
