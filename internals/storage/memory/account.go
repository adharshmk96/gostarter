package memory

import (
	"context"
	"gostarter/infra"
	"gostarter/internals/domain"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel/trace"
)

type accountRepository struct {
	logger   *slog.Logger
	tracer   trace.Tracer
	accounts []domain.Account
}

func NewAccountRepository(container *infra.Container) domain.AccountRepository {
	logger := container.Logger.With("path", "accountRepository")
	return &accountRepository{
		logger:   logger,
		tracer:   container.Tracer,
		accounts: []domain.Account{},
	}
}

func (a *accountRepository) CreateAccount(ctx context.Context, account *domain.Account) error {
	_, span := a.tracer.Start(ctx, "AccountRepository.CreateAccount")
	defer span.End()

	id := len(a.accounts) + 1
	account.Id = id
	now := time.Now()
	account.CreatedAt = now
	account.UpdatedAt = now

	if account.Username == "" {
		account.Username = account.Email
	}

	a.accounts = append(a.accounts, *account)

	return nil
}

func (a *accountRepository) GetAccountByID(ctx context.Context, id int) (*domain.Account, error) {
	_, span := a.tracer.Start(ctx, "AccountRepository.GetAccountByID")
	defer span.End()

	var account *domain.Account

	for _, acc := range a.accounts {
		if acc.Id == id {
			account = &acc
			break
		}
	}

	if account == nil || account.Id == 0 {
		return nil, domain.ErrAccountNotFound
	}

	return account, nil
}

func (a *accountRepository) GetAccountByEmail(ctx context.Context, email string) (*domain.Account, error) {
	_, span := a.tracer.Start(ctx, "AccountRepository.GetAccountByEmail")
	defer span.End()

	var account domain.Account

	for _, acc := range a.accounts {
		if acc.Email == email {
			account = acc
			break
		}
	}

	if account.Email == "" {
		return nil, domain.ErrAccountNotFound
	}

	return &account, nil
}

func (a *accountRepository) GetAccountByUsername(ctx context.Context, username string) (*domain.Account, error) {
	_, span := a.tracer.Start(ctx, "AccountRepository.GetAccountByUsername")
	defer span.End()

	var account domain.Account

	for _, acc := range a.accounts {
		if acc.Username == username {
			account = acc
			break
		}
	}

	if account.Username == "" {
		return nil, domain.ErrAccountNotFound
	}

	return &account, nil
}

func (a *accountRepository) UpdateAccount(ctx context.Context, account *domain.Account) error {
	_, span := a.tracer.Start(ctx, "AccountRepository.UpdateAccount")
	defer span.End()

	var updatedAccount *domain.Account

	for i, acc := range a.accounts {
		if acc.Id == account.Id {
			updatedAccount = account
			a.accounts[i] = *account
			break
		}
	}

	if updatedAccount == nil || updatedAccount.Id == 0 {
		return domain.ErrAccountNotFound
	}

	return nil
}

func (a *accountRepository) DeleteAccount(ctx context.Context, id int) error {
	_, span := a.tracer.Start(ctx, "AccountRepository.DeleteAccount")
	defer span.End()

	var deletedAccount *domain.Account

	for i, acc := range a.accounts {
		if acc.Id == id {
			deletedAccount = &a.accounts[i]
			a.accounts = append(a.accounts[:i], a.accounts[i+1:]...)
			break
		}
	}

	if deletedAccount == nil || deletedAccount.Id == 0 {
		return domain.ErrAccountNotFound
	}

	return nil
}

func (a *accountRepository) ListAccounts(ctx context.Context, pagination *domain.Pagination) ([]*domain.Account, error) {
	_, span := a.tracer.Start(ctx, "AccountRepository.ListAccounts")
	defer span.End()

	offset := pagination.GetOffset()

	if offset >= len(a.accounts) {
		return []*domain.Account{}, nil
	}

	end := offset + pagination.Size

	if end > len(a.accounts) {
		end = len(a.accounts)
	}

	var result []*domain.Account

	for i := offset; i < end; i++ {
		result = append(result, &a.accounts[i])
	}

	return result, nil
}
