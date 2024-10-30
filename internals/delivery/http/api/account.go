package api

import (
	"encoding/json"
	"go.opentelemetry.io/otel/trace"
	"gostarter/infra"
	"gostarter/infra/config"
	"gostarter/internals/delivery/http/helpers"
	"gostarter/internals/domain"
	"log/slog"
	"net/http"
)

type AccountHandler struct {
	logger *slog.Logger
	tracer *trace.Tracer

	accountService domain.AccountService
	tokenService   domain.TokenService
}

func NewAccountHandler(
	container *infra.Container,
	accountService domain.AccountService,
	tokenService domain.TokenService,
) *AccountHandler {
	return &AccountHandler{
		logger:         container.Logger,
		tracer:         container.Tracer,
		accountService: accountService,
		tokenService:   tokenService,
	}
}

type RegisterAccountRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterAccountResponse struct {
	Message string `json:"message"`
}

// @Router /v1/auth/register [post]
// @Tags Account
// @Summary Register a new account
// @Description Register a new account
// @Accept json
// @Produce json
// @Param account body RegisterAccountRequest true "Account to register"
// @Success 200 {object} RegisterAccountResponse
// @Failure 400 {object} helpers.GeneralResponse
// @Failure 500 {object} helpers.GeneralResponse
func (a *AccountHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Parse request
	req, err := helpers.UnmarshalData[RegisterAccountRequest](r.Body)
	if err != nil {
		errorResponse := helpers.GeneralResponse{
			Message: "invalid request",
			Errors: []string{
				err.Error(),
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	acc := &domain.Account{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	// Register account
	err = a.accountService.Register(acc)
	if err != nil {
		errorResponse := helpers.GeneralResponse{
			Message: "failed to register account",
			Errors: []string{
				err.Error(),
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Generate JWT
	token, err := a.tokenService.GenerateJWT(acc.ID, acc.Username, acc.Roles)
	if err != nil {
		errorResponse := helpers.GeneralResponse{
			Message: "failed to generate token",
			Errors: []string{
				err.Error(),
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Set token in http only cookie
	http.SetCookie(w, &http.Cookie{
		Path:     "/",
		Name:     config.AUTH_COOKIE_NAME,
		Value:    token,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	// Response
	resp := RegisterAccountResponse{
		Message: "account registered successfully",
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// @Router /v1/auth/login [post]
// @Tags Account
// @Summary Login an account
// @Description Login an account
// @Accept json
// @Produce json
// @Param account body LoginRequest true "Account to login"
// @Success 200 {object} helpers.GeneralResponse
// @Failure 400 {object} helpers.GeneralResponse
// @Failure 500 {object} helpers.GeneralResponse
func (a *AccountHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Parse request
	req, err := helpers.UnmarshalData[LoginRequest](r.Body)
	if err != nil {
		errorResponse := helpers.GeneralResponse{
			Message: "invalid request",
			Errors: []string{
				err.Error(),
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Authenticate account
	acc, err := a.accountService.Authenticate(req.Username, req.Password)
	if err != nil {
		errorResponse := helpers.GeneralResponse{
			Message: "invalid credentials",
			Errors: []string{
				err.Error(),
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Generate JWT
	token, err := a.tokenService.GenerateJWT(acc.ID, acc.Username, acc.Roles)
	if err != nil {
		errorResponse := helpers.GeneralResponse{
			Message: "failed to generate token",
			Errors: []string{
				err.Error(),
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Set token in http only cookie
	http.SetCookie(w, &http.Cookie{
		Path:     "/",
		Name:     config.AUTH_COOKIE_NAME,
		Value:    token,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	// Response
	resp := helpers.GeneralResponse{
		Message: "login successful",
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// @Router /v1/auth/logout [post]
// @Tags Account
// @Summary Logout an account
// @Description Logout an account
// @Accept json
// @Produce json
// @Success 200 {object} helpers.GeneralResponse
// @Failure 500 {object} helpers.GeneralResponse
func (a *AccountHandler) Logout(w http.ResponseWriter, _ *http.Request) {
	// Set token in http only cookie
	http.SetCookie(w, &http.Cookie{
		Path:     "/",
		Name:     config.AUTH_COOKIE_NAME,
		Value:    "",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	// Response
	resp := helpers.GeneralResponse{
		Message: "logout successful",
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

type ProfileResponse struct {
	ID       int      `json:"id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}

// @Router /v1/auth/profile [get]
// @Tags Account
// @Summary Get account profile
// @Description Get account profile
// @Accept json
// @Produce json
// @Success 200 {object} ProfileResponse
// @Failure 500 {object} helpers.GeneralResponse
func (a *AccountHandler) Profile(w http.ResponseWriter, r *http.Request) {
	// Get account from context
	acc, ok := r.Context().Value("account").(*domain.Account)
	if !ok {
		errorResponse := helpers.GeneralResponse{
			Message: "invalid account",
			Errors: []string{
				"account not found",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Response
	resp := ProfileResponse{
		ID:       acc.ID,
		Username: acc.Username,
		Roles:    acc.Roles,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
