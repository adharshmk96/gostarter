package api

import (
	"gostarter/infra"
	"gostarter/infra/config"
	"gostarter/internals/delivery/http/helpers"
	"gostarter/internals/domain"
	"log/slog"
	"net/http"

	"go.opentelemetry.io/otel/trace"
)

type AccountHandler struct {
	logger *slog.Logger
	tracer trace.Tracer

	accountService domain.AccountService
	tokenService   domain.TokenService
}

func NewAccountHandler(
	container *infra.Container,
	accountService domain.AccountService,
	tokenService domain.TokenService,
) domain.AccountHandler {
	logger := container.Logger.With("path", "AccountHandler")
	return &AccountHandler{
		logger:         logger,
		tracer:         container.Tracer,
		accountService: accountService,
		tokenService:   tokenService,
	}
}

type RegisterAccountRequest struct {
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
	ctx, span := a.tracer.Start(r.Context(), "AccountHandler.Register")
	defer span.End()

	// Parse request
	req, err := helpers.ParseRequest[RegisterAccountRequest](r.Body)
	if err != nil {
		errorResponse := helpers.GeneralResponse{
			Message: "invalid request",
			Errors: []string{
				err.Error(),
			},
		}
		_ = helpers.WriteResponse(w, http.StatusBadRequest, errorResponse)
		return
	}

	acc := &domain.Account{
		Username: req.Email,
		Email:    req.Email,
		Password: req.Password,
		Roles:    []string{domain.ROLE_USER},
	}

	// Register account
	err = a.accountService.Register(ctx, acc)
	if err != nil {
		errorResponse := helpers.GeneralResponse{
			Message: "failed to register account",
			Errors: []string{
				err.Error(),
			},
		}
		_ = helpers.WriteResponse(w, http.StatusInternalServerError, errorResponse)
		return
	}

	// Generate JWT
	token, err := a.tokenService.GenerateJWT(acc.Id, acc.Email, acc.Roles)
	if err != nil {
		errorResponse := helpers.GeneralResponse{
			Message: "failed to generate token",
			Errors: []string{
				err.Error(),
			},
		}
		_ = helpers.WriteResponse(w, http.StatusInternalServerError, errorResponse)
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

	_ = helpers.WriteResponse(w, http.StatusOK, resp)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// @Router /v1/auth/login [post]
// @Tags Account
// @Summary Login an account
// @Description Login an account
// @Accept json
// @Produce json
// @Param account body LoginRequest true "Login Details"
// @Success 200 {object} helpers.GeneralResponse
// @Failure 400 {object} helpers.GeneralResponse
// @Failure 500 {object} helpers.GeneralResponse
func (a *AccountHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx, span := a.tracer.Start(r.Context(), "AccountHandler.Login")
	defer span.End()

	// Parse request
	req, err := helpers.ParseRequest[LoginRequest](r.Body)
	if err != nil {
		errorResponse := helpers.GeneralResponse{
			Message: "invalid request",
			Errors: []string{
				err.Error(),
			},
		}
		_ = helpers.WriteResponse(w, http.StatusBadRequest, errorResponse)
		return
	}

	// Authenticate account
	acc, err := a.accountService.Authenticate(ctx, req.Email, req.Password)
	if err != nil {
		errorResponse := helpers.GeneralResponse{
			Message: "invalid credentials",
			Errors: []string{
				err.Error(),
			},
		}
		_ = helpers.WriteResponse(w, http.StatusInternalServerError, errorResponse)
		return
	}

	// Generate JWT
	token, err := a.tokenService.GenerateJWT(acc.Id, acc.Email, acc.Roles)
	if err != nil {
		errorResponse := helpers.GeneralResponse{
			Message: "failed to generate token",
			Errors: []string{
				err.Error(),
			},
		}
		_ = helpers.WriteResponse(w, http.StatusInternalServerError, errorResponse)
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

	_ = helpers.WriteResponse(w, http.StatusOK, resp)
}

// @Router /v1/auth/logout [post]
// @Tags Account
// @Summary Logout an account
// @Description Logout an account
// @Accept json
// @Produce json
// @Success 200 {object} helpers.GeneralResponse
// @Failure 500 {object} helpers.GeneralResponse
func (a *AccountHandler) Logout(w http.ResponseWriter, r *http.Request) {
	_, span := a.tracer.Start(r.Context(), "AccountHandler.Logout")
	defer span.End()

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

	_ = helpers.WriteResponse(w, http.StatusOK, resp)
}

type ProfileResponse struct {
	ID    int      `json:"id"`
	Email string   `json:"email"`
	Roles []string `json:"roles"`
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
	_, span := a.tracer.Start(r.Context(), "AccountHandler.Profile")
	defer span.End()

	// Get account from context
	acc, err := helpers.GetAccountFromContext(r.Context())
	if err != nil {
		errorResponse := helpers.GeneralResponse{
			Message: "invalid account",
			Errors: []string{
				"account not found",
			},
		}
		_ = helpers.WriteResponse(w, http.StatusInternalServerError, errorResponse)
		return
	}

	// Response
	resp := ProfileResponse{
		ID:    acc.Id,
		Email: acc.Email,
		Roles: acc.Roles,
	}

	_ = helpers.WriteResponse(w, http.StatusOK, resp)
}
