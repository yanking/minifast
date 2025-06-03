package srv

import (
	"minifast/app/user/srv/config"
	"minifast/pkg/app"
	"minifast/pkg/log"
)

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
