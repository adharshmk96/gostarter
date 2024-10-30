package domain

import "net/http"

type Account struct {
	ID int `json:"id"`

	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`

	Roles []string `json:"roles"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type AccountHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	Profile(w http.ResponseWriter, r *http.Request)
	ChangePassword(w http.ResponseWriter, r *http.Request)
}

type AccountService interface {
	Register(account *Account) error
	GetAccountByID(id int) (*Account, error)
	GetAccountByEmail(email string) (*Account, error)
	GetAccountByUsername(username string) (*Account, error)
	UpdateAccount(account *Account) error
	DeleteAccount(id int) error

	ListAccounts() ([]*Account, error)
}

type AccountRepository interface {
	CreateAccount(account *Account) error
	GetAccountByID(id int) (*Account, error)
	GetAccountByEmail(email string) (*Account, error)
	GetAccountByUsername(username string) (*Account, error)
	UpdateAccount(account *Account) error
	DeleteAccount(id int) error

	ListAccounts() ([]*Account, error)
}
