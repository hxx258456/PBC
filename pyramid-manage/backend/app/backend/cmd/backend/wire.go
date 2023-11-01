//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"pyramid/pyramid-manage/backend/app/backend/internal/biz"
	"pyramid/pyramid-manage/backend/app/backend/internal/conf"
	"pyramid/pyramid-manage/backend/app/backend/internal/data"
	"pyramid/pyramid-manage/backend/app/backend/internal/server"
	"pyramid/pyramid-manage/backend/app/backend/internal/service"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
