package domain

import (
	"errors"
	"gostarter/internals/core/utils"
	"net/http"
)

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
	Authenticate(email, password string) (*Account, error)

	GetAccountByID(id int) (*Account, error)
	GetAccountByEmail(email string) (*Account, error)
	GetAccountByUsername(username string) (*Account, error)
	UpdateAccount(account *Account) error
	DeleteAccount(id int) error

	ListAccounts(utils.PaginationParams) ([]*Account, error)
}

type TokenService interface {
	GenerateJWT(id int, username string, roles []string) (string, error)
	VerifyJWT(token string) (bool, error)
	ExtractAccount(token string) (*Account, error)
}

var (
	ErrLoadingKey   = errors.New("error loading key")
	ErrInvalidToken = errors.New("invalid token")
)

type AccountRepository interface {
	CreateAccount(account *Account) error
	GetAccountByID(id int) (*Account, error)
	GetAccountByEmail(email string) (*Account, error)
	GetAccountByUsername(username string) (*Account, error)
	UpdateAccount(account *Account) error
	DeleteAccount(id int) error

	ListAccounts(utils.PaginationParams) ([]*Account, error)
}

// Errors
var (
	ErrAccountNotFound = errors.New("account not found")
)
