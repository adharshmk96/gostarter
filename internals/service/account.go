package service

import "gostarter/internals/domain"

type accountService struct {
	accountRepo domain.AccountRepository
}

func NewAccountService(accountRepo domain.AccountRepository) domain.AccountService {
	return &accountService{
		accountRepo: accountRepo,
	}
}

func (a *accountService) Register(account *domain.Account) error {
	//TODO implement me
	panic("implement me")
}

func (a *accountService) GetAccountByID(id int) (*domain.Account, error) {
	//TODO implement me
	panic("implement me")
}

func (a *accountService) GetAccountByEmail(email string) (*domain.Account, error) {
	//TODO implement me
	panic("implement me")
}

func (a *accountService) GetAccountByUsername(username string) (*domain.Account, error) {
	//TODO implement me
	panic("implement me")
}

func (a *accountService) UpdateAccount(account *domain.Account) error {
	//TODO implement me
	panic("implement me")
}

func (a *accountService) DeleteAccount(id int) error {
	//TODO implement me
	panic("implement me")
}

func (a *accountService) ListAccounts() ([]*domain.Account, error) {
	//TODO implement me
	panic("implement me")
}
