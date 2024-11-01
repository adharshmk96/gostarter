package vault

import (
	vaultapi "github.com/hashicorp/vault/api"
)

type VaultService struct {
	client *vaultapi.Client
}

func NewVaultService(address, token string) (*VaultService, error) {
	config := vaultapi.DefaultConfig()
	config.Address = address

	client, err := vaultapi.NewClient(config)
	if err != nil {
		return nil, err
	}

	client.SetToken(token)

	return &VaultService{
		client: client,
	}, nil
}

func (v *VaultService) ReadSecret(path string) (map[string]interface{}, error) {
	secret, err := v.client.Logical().Read(path)
	if err != nil {
		return nil, err
	}

	return secret.Data, nil
}
