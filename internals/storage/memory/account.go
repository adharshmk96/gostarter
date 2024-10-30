package memory

import "gostarter/internals/domain"

type accountRepository struct {
	accounts []domain.Account
}

func NewAccountRepository() domain.AccountRepository {
	return &accountRepository{
		accounts: []domain.Account{},
	}
}

func (a *accountRepository) CreateAccount(account *domain.Account) error {
	//TODO implement me
	panic("implement me")
}

func (a *accountRepository) GetAccountByID(id int) (*domain.Account, error) {
	//TODO implement me
	panic("implement me")
}

func (a *accountRepository) GetAccountByEmail(email string) (*domain.Account, error) {
	//TODO implement me
	panic("implement me")
}

func (a *accountRepository) GetAccountByUsername(username string) (*domain.Account, error) {
	//TODO implement me
	panic("implement me")
}

func (a *accountRepository) UpdateAccount(account *domain.Account) error {
	//TODO implement me
	panic("implement me")
}

func (a *accountRepository) DeleteAccount(id int) error {
	//TODO implement me
	panic("implement me")
}

func (a *accountRepository) ListAccounts() ([]*domain.Account, error) {
	//TODO implement me
	panic("implement me")
}
