package web

import (
	"context"
	"gostarter/infra/config"
	"gostarter/internals/delivery/http/helpers"
	"gostarter/internals/domain"
	"gostarter/pkg/rendering"
	"net/http"
)

type AccountWebHandler struct {
	accountService domain.AccountService
	tokenService   domain.TokenService
	renderer       *rendering.HtmlRenderer
}

func NewAccountWebHandler(
	tokenService domain.TokenService,
	accountService domain.AccountService,
) *AccountWebHandler {
	renderer := rendering.NewHtmlRenderer(config.TEMPLATE_DIR)
	return &AccountWebHandler{
		accountService: accountService,
		tokenService:   tokenService,
		renderer:       renderer,
	}
}

func (h *AccountWebHandler) GetRegisterMember(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Register",
	}
	err := h.renderer.RenderWithLayout(
		w, "layout/main.html", "register.html", data,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *AccountWebHandler) PostRegisterMember(w http.ResponseWriter, r *http.Request) {
	// Parse the form
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the form data
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	acc := &domain.Account{
		Username: email,
		Email:    email,
		Password: password,
		Roles:    []string{domain.ROLE_USER},
	}

	// Register the member
	err = h.accountService.Register(context.Background(), acc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate JWT
	token, err := h.tokenService.GenerateJWT(acc.Id, acc.Username, acc.Roles)
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

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func (h *AccountWebHandler) GetLogin(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Login",
	}
	err := h.renderer.RenderWithLayout(
		w, "layout/main.html", "login.html", data,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *AccountWebHandler) PostLogin(w http.ResponseWriter, r *http.Request) {
	// Parse the form
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the form data
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	// Authenticate
	acc, err := h.accountService.Authenticate(context.Background(), email, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate JWT
	token, err := h.tokenService.GenerateJWT(acc.Id, acc.Username, acc.Roles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func (h *AccountWebHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	acc, err := helpers.GetAccountFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":   "Profile",
		"Account": acc,
	}
	err = h.renderer.RenderWithLayout(
		w, "layout/main.html", "profile.html", data,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *AccountWebHandler) PostLogout(w http.ResponseWriter, r *http.Request) {
	// Set token in http only cookie
	http.SetCookie(w, &http.Cookie{
		Path:     "/",
		Name:     config.AUTH_COOKIE_NAME,
		Value:    "",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
