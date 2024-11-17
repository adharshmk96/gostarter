package resolver

import "gostarter/internals/domain"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	AccountService domain.AccountService
}
