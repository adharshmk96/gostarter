package service

import (
	"github.com/adharshmk96/goutils/token"
	"github.com/golang-jwt/jwt/v5"
	"gostarter/infra/config"
	"gostarter/internals/domain"
	"gostarter/pkg/utils"
	"time"
)

type tokenService struct {
	jwtUtil     *token.JWTUtil
	tokenExpiry int
}

func (a *tokenService) GenerateJWT(id int, email string, roles []string) (string, error) {
	claims := jwt.MapClaims{
		"userId": id,
		"email":  email,
		"roles":  roles,
		"exp":    time.Now().Add(time.Hour * time.Duration(a.tokenExpiry)).Unix(),
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
	email, ok := decodedJwt.Claims.(jwt.MapClaims)["email"].(string)
	if !ok {
		return nil, domain.ErrInvalidToken
	}

	userAccount := &domain.Account{
		Id:       userId,
		Username: email,
		Email:    email,
		Roles:    roles,
	}

	return userAccount, nil
}

func NewTokenService(cfg config.JWTConfig) domain.TokenService {

	privateKey, publicKey, err := utils.LoadECDSAKeyPair(cfg.PrivateKeyPath, cfg.PublicKeyPath)
	if err != nil {
		panic(err)
	}

	jwtUtil := token.NewJwtUtil(token.JWTConfig{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	})

	return &tokenService{
		jwtUtil:     jwtUtil,
		tokenExpiry: cfg.ExpirationHours,
	}
}
