package resolver

import (
	"gostarter/infra"
	"gostarter/internals/di"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Container *infra.Container
	ServiceDi *di.ServiceContainer
}
