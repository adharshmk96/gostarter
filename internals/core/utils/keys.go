package utils

import (
	"crypto/ecdsa"
	"github.com/adharshmk96/goutils/token"
)

func LoadECDSAKeyPair(
	privateKeyPath string,
	publicKeyPath string,
) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	privateKey, err := token.LoadPrivateKeyFromPath(privateKeyPath)
	if err != nil {
		return nil, nil, err
	}
	publicKey, err := token.LoadPublicKeyFromPath(publicKeyPath)
	if err != nil {
		return nil, nil, err
	}

	return privateKey, publicKey, nil
}
