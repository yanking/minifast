package srv

import (
	"github.com/google/wire"
	"minifast/app/pkg/options"
	"minifast/app/user/srv/config"
	gapp "minifast/gmicro/app"
	"minifast/gmicro/server/rpcserver"
	"minifast/pkg/app"
	"minifast/pkg/log"
)

var ProviderSet = wire.NewSet(NewUserApp, NewUserRPCServer)

func NewApp(basename string) *app.App {
	cfg := config.New()
	appl := app.NewApp(
		"user",
		basename,
		app.WithOptions(cfg),
		app.WithRunFunc(run(cfg)),
		//app.WithNoConfig(), //设置不读取配置文件
	)

	return appl
}

func NewUserApp(logOpts *log.Options, serverOpts *options.ServerOptions, rpcServer *rpcserver.Server) (*gapp.App, error) {
	//初始化log
	log.Init(logOpts)
	defer log.Flush()

	return gapp.New(
		gapp.WithName(serverOpts.Name),
		gapp.WithRPCServer(rpcServer),
	), nil
}

func run(cfg *config.Config) app.RunFunc {
	return func(baseName string) error {
		userApp, err := initApp(cfg.Log, cfg.Server, cfg.Telemetry, cfg.MySQLOptions)
		if err != nil {
			return err
		}

		//启动
		if err = userApp.Run(); err != nil {
			log.Errorf("run user app error: %s", err)
			return err
		}
		return nil
	}
}
