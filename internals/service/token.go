package service

import (
	"github.com/adharshmk96/goutils/token"
	"github.com/golang-jwt/jwt/v5"
	"gostarter/internals/domain"
	"gostarter/pkg/utils"
	"time"
)

type tokenService struct {
	tokenExpiry int
	jwtUtil     *token.JWTUtil
}

func (a *tokenService) GenerateJWT(id int, username string, roles []string) (string, error) {
	if len(roles) == 0 {
		roles = append(roles, "user")
	}

	claims := jwt.MapClaims{
		"userId":   id,
		"username": username,
		"roles":    roles,
		"exp":      time.Now().Add(time.Hour * time.Duration(a.tokenExpiry)).Unix(),
	}

	return a.jwtUtil.EncodeJWT(claims)
}

func (a *tokenService) VerifyJWT(userJWT string) (bool, error) {
	_, err := a.jwtUtil.DecodeJWT(userJWT)
	if err != nil {
		return false, err
	}

	return true, nil
}

func getList[T any](key string, decodedJwt *jwt.Token) ([]T, bool) {
	data, ok := decodedJwt.Claims.(jwt.MapClaims)[key].([]interface{})
	if !ok {
		return nil, false
	}

	items := make([]T, len(data))
	for i, v := range data {
		items[i] = v.(T)
	}

	return items, true
}

func (a *tokenService) ExtractAccount(userJWT string) (*domain.Account, error) {
	decodedJwt, err := a.jwtUtil.DecodeJWT(userJWT)
	if err != nil {
		return nil, err
	}

	userIdFloat, ok := decodedJwt.Claims.(jwt.MapClaims)["userId"].(float64)
	userId := int(userIdFloat)
	if !ok {
		return nil, domain.ErrInvalidToken
	}
	roles, ok := getList[string]("roles", decodedJwt)
	if !ok {
		return nil, domain.ErrInvalidToken
	}
	username, ok := decodedJwt.Claims.(jwt.MapClaims)["username"].(string)
	if !ok {
		return nil, domain.ErrInvalidToken
	}

	userAccount := &domain.Account{
		ID:       userId,
		Username: username,
		Roles:    roles,
	}

	return userAccount, nil
}

func NewTokenService(privateKeyPath, publicKeyPath string, expiry int) (domain.TokenService, error) {

	privateKey, publicKey, err := utils.LoadECDSAKeyPair(privateKeyPath, publicKeyPath)
	if err != nil {
		return nil, domain.ErrLoadingKey
	}

	jwtUtil := token.NewJwtUtil(token.JWTConfig{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	})

	return &tokenService{
		tokenExpiry: expiry,
		jwtUtil:     jwtUtil,
	}, nil
}
