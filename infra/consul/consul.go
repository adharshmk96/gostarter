package consul

import (
	consulapi "github.com/hashicorp/consul/api"
)

type ConsulService struct {
	client *consulapi.Client
}

func NewConsulService(address string) (*ConsulService, error) {
	config := consulapi.DefaultConfig()
	config.Address = address

	client, err := consulapi.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &ConsulService{
		client: client,
	}, nil
}

func (c *ConsulService) GetKVPair(key string) (string, error) {
	kv := c.client.KV()

	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return "", err
	}

	if pair == nil {
		return "", nil
	}

	return string(pair.Value), nil
}
