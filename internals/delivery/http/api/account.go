package api

import (
	"encoding/json"
	"gostarter/internals/delivery/http/common"
	"gostarter/internals/domain"
	"net/http"
)

type AccountHandler struct {
	accountService domain.AccountService
}

func NewAccountHandler(accountService domain.AccountService) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
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
// @Failure 400 {object} common.GeneralResponse
// @Failure 500 {object} common.GeneralResponse
func (a *AccountHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Parse request
	req, err := common.UnmarshalData[RegisterAccountRequest](r.Body)
	if err != nil {
		errorResponse := common.GeneralResponse{
			Message: "Invalid request",
			Errors: []string{
				err.Error(),
			},
		}
		jsonData, _ := json.Marshal(errorResponse)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(jsonData)
		return
	}

	// Register account
	err = a.accountService.Register(
		&domain.Account{
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
		},
	)
	if err != nil {
		errorResponse := common.GeneralResponse{
			Message: "Failed to register account",
			Errors: []string{
				err.Error(),
			},
		}
		jsonData, _ := json.Marshal(errorResponse)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(jsonData)
		return
	}

	// Response
	resp := RegisterAccountResponse{
		Message: "Account registered successfully",
	}

	jsonData, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonData)
}
