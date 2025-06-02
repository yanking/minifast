package restserver

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/penglongli/gin-metrics/ginmetrics"
	mws "minifast/gmicro/server/restserver/middlewares"
	"minifast/gmicro/server/restserver/pprof"
	"minifast/gmicro/server/restserver/validation"
	"minifast/pkg/log"
	"net/http"
	"time"
)

type JwtInfo struct {
	//defaults to "JWT"
	Realm string
	// defaults to empty
	Key string
	// defaults to 7 days
	Timeout time.Duration
	// defaults to 7 days
	NaxRefresh time.Duration
}

type Server struct {
	*gin.Engine

	// 端口号，默认值8080
	port int

	// 开发模式，默认debug
	mode string

	//是否开启健康检查接口， 默认开启， 如果开启会自动添加 /health 接口
	healthz bool

	//是否开启pprof接口， 默认开启， 如果开启会自动添加 /debug/pprof 接口
	enableProfiling bool

	//是否开启metrics接口， 默认开启， 如果开启会自动添加 /metrics 接口
	enableMetrics bool

	enableTracing bool

	//中间件
	middlewares []string

	//jwt配置信息
	jwt *JwtInfo

	//翻译器, 默认值 zh
	transName string
	trans     ut.Translator

	server *http.Server

	serviceName string
}

func New(opts ...ServerOption) *Server {
	srv := &Server{
		port:            8080,
		mode:            gin.DebugMode,
		healthz:         true,
		enableProfiling: true,
		enableMetrics:   true,
		jwt: &JwtInfo{
			Realm:      "JWT",
			Key:        "sdhHUNYsdYT2457JuliHIUGbnbv",
			Timeout:    7 * 24 * time.Hour,
			NaxRefresh: 7 * 24 * time.Hour,
		},
		Engine:      gin.Default(),
		transName:   "zh",
		serviceName: "minifast",
	}

	for _, opt := range opts {
		opt(srv)
	}

	for _, m := range srv.middlewares {
		mw, ok := mws.Middlewares[m]
		if !ok {
			log.Warnf("can not find middleware: %s", m)
			continue
			//panic(errors.Errorf("can not find middleware: %s", m))
		}

		log.Infof("intall middleware: %s", m)
		srv.Use(mw)
	}

	return srv
}

func (s *Server) Translator() ut.Translator {
	return s.trans
}

func (s *Server) Start(ctx context.Context) error {
	// 设置开发模式，打印路由信息
	if s.mode != gin.DebugMode && s.mode != gin.ReleaseMode && s.mode != gin.TestMode {
		return errors.New("mode must be one of debug/release/test")
	}

	// 设置开发模式，打印路由信息
	gin.SetMode(s.mode)
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Infof("%-6s %-25s --> %s (%d handlers)\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}

	// 初始化翻译器
	err := s.initTrans(s.transName)
	if err != nil {
		log.Errorf("initTrans error %s", err.Error())
		return err
	}

	//注册mobile验证码
	validation.RegisterMobile(s.trans)

	//根据配置初始化pprof路由
	if s.enableProfiling {
		pprof.Register(s.Engine)
	}

	if s.enableMetrics {
		// get global Monitor object
		m := ginmetrics.GetMonitor()
		// +optional set metric path, default /debug/metrics
		m.SetMetricPath("/metrics")
		// +optional set slow time, default 5s
		// +optional set request duration, default {0.1, 0.3, 1.2, 5, 10}
		// used to p95, p99
		m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})
		m.Use(s)
	}

	if s.enableTracing {
		s.Use(mws.TracingHandler(s.serviceName))
	}

	log.Infof("rest server is running on port %d", s.port)
	address := fmt.Sprintf(":%d", s.port)
	s.server = &http.Server{
		Addr:    address,
		Handler: s.Engine,
	}
	_ = s.SetTrustedProxies(nil)
	if err = s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	log.Infof("rest server is stopping")
	if s.server == nil {
		return nil
	}
	if err := s.server.Shutdown(ctx); err != nil {
		log.Errorf("rest server shutdown error: %s", err.Error())
		return err
	}
	log.Info("rest server stopped")
	return nil
}
