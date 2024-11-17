package resolver

import (
	"gostarter/infra"
	"gostarter/internals/domain"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Container      *infra.Container
	AccountService domain.AccountService
}
