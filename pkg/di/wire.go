//go:build wireinject
// +build wireinject

package di

import (
	http "main/pkg/api"
	handler "main/pkg/api/handler"
	config "main/pkg/config"
	db "main/pkg/db"
	repository "main/pkg/repository"
	usecase "main/pkg/usecase"

	"github.com/google/wire"
)

func InitializeAPI(cfg config.Config) (*http.ServerHTTP, error) {
	wire.Build(db.ConnectDatabase, http.NewServerHTTP, repository.NewUserRepository, usecase.NewUserUseCase, handler.NewUserHandler)
	return &http.ServerHTTP{}, nil
}
