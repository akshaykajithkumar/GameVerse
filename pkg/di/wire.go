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
	wire.Build(db.ConnectDatabase, http.NewServerHTTP, repository.NewUserRepository, usecase.NewUserUseCase, handler.NewUserHandler, repository.NewOtpRepository, usecase.NewOtpUseCase, handler.NewOtpHandler, repository.NewAdminRepository, usecase.NewAdminUseCase, handler.NewAdminHandler, repository.NewCategoryRepository, usecase.NewCategoryUseCase, handler.NewCategoryHandler, repository.NewVideoRepository, usecase.NewVideoUseCase, handler.NewVideoHandler, repository.NewsubscriptionRepository, usecase.NewSubscriptionUseCase, handler.NewSubscriptionHandler)
	return &http.ServerHTTP{}, nil
}
