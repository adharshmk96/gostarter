package domain

type TokenService interface {
	GenerateJWT(id int, username string, roles []string) (string, error)
	VerifyJWT(token string) (bool, error)
	ExtractAccount(token string) (*Account, error)
}
