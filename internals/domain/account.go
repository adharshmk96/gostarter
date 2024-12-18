package domain

import (
	"context"
	"errors"
	"net/http"
	"time"
)

const (
	ROLE_USER  = "user"
	ROLE_ADMIN = "admin"
)

type Account struct {
	Id int `json:"id"`

	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`

	Roles []string `json:"roles"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AccountHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	Profile(w http.ResponseWriter, r *http.Request)
}

var (
	ErrGettingAccountInfo = errors.New("error getting account info")
)

type AccountService interface {
	Register(ctx context.Context, account *Account) error
	Authenticate(ctx context.Context, email, password string) (*Account, error)

	GetAccountByID(ctx context.Context, id int) (*Account, error)
	GetAccountByEmail(ctx context.Context, email string) (*Account, error)
	UpdateAccount(ctx context.Context, account *Account) error
	DeleteAccount(ctx context.Context, id int) error

	ListAccounts(context.Context, *Pagination) ([]*Account, error)
}

var (
	ErrLoadingKey   = errors.New("error loading key")
	ErrInvalidToken = errors.New("invalid token")
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, account *Account) error
	GetAccountByID(ctx context.Context, id int) (*Account, error)
	GetAccountByEmail(ctx context.Context, email string) (*Account, error)
	UpdateAccount(ctx context.Context, account *Account) error
	DeleteAccount(ctx context.Context, id int) error

	ListAccounts(context.Context, *Pagination) ([]*Account, error)
}

// Errors
var (
	ErrAccountNotFound = errors.New("account not found")
)
