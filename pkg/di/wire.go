//go:build wireinject
// +build wireinject

package di

import (
	http "./pkg/api"
	config "./pkg/config"
)

func InitializeAPI(cfg config.Config) (*http.ServerHTTP, error) {

	return &http.ServerHTTP{}, nil
}
