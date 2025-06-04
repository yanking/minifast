//go:build wireinject
// +build wireinject

package srv

import (
	"github.com/google/wire"
	"minifast/app/pkg/options"
	"minifast/app/user/srv/controller/user"
	"minifast/app/user/srv/data/v1/db"
	v1 "minifast/app/user/srv/service/v1"
	gapp "minifast/gmicro/app"
	"minifast/pkg/log"
)

func initApp(*log.Options, *options.ServerOptions, *options.TelemetryOptions, *options.MySQLOptions) (*gapp.App, error) {
	wire.Build(ProviderSet, v1.ProviderSet, db.ProviderSet, user.ProviderSet)
	return &gapp.App{}, nil
}
