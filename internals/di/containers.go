package di

import (
	"gostarter/infra"
	"gostarter/internals/delivery/http/api"
	"gostarter/internals/delivery/http/web"
	"gostarter/internals/domain"
	"gostarter/internals/service"
	"gostarter/internals/storage/pgstorage"
)

type RepoContainer struct {
	AccountRepo domain.AccountRepository
}

func NewRepoContainer(container *infra.Container) *RepoContainer {
	return &RepoContainer{
		AccountRepo: pgstorage.NewAccountRepository(container),
	}
}

type ServiceContainer struct {
	TokenService   domain.TokenService
	AccountService domain.AccountService
}

func NewServiceContainer(container *infra.Container, repoContainer *RepoContainer) *ServiceContainer {
	return &ServiceContainer{
		TokenService:   service.NewTokenService(container.Cfg.JWT),
		AccountService: service.NewAccountService(container, repoContainer.AccountRepo),
	}
}

type HandlerContainer struct {
	AccountHandler    domain.AccountHandler
	AccountWebHandler *web.AccountWebHandler
}

func NewHandlerContainer(container *infra.Container, serviceContainer *ServiceContainer) *HandlerContainer {
	return &HandlerContainer{
		AccountHandler:    api.NewAccountHandler(container, serviceContainer.AccountService, serviceContainer.TokenService),
		AccountWebHandler: web.NewAccountWebHandler(serviceContainer.TokenService, serviceContainer.AccountService),
	}
}
