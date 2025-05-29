package app

import (
	"context"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"minifast/gmicro/registry"
	gs "minifast/gmicro/server"
	"minifast/pkg/log"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type App struct {
	opts options

	lk       sync.Mutex
	instance *registry.ServiceInstance

	cancel func()
}

func New(opts ...Option) *App {
	o := options{
		sigs:             []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		registrarTimeout: time.Second * 10,
		stopTimeout:      time.Second * 10,
	}

	o.id = uuid.NewString()

	for _, opt := range opts {
		opt(&o)
	}

	return &App{
		opts: o,
	}
}

func (a *App) Run() error {
	//注册的信息
	instance, err := a.buildInstance()
	if err != nil {
		return err
	}

	//这个变量可能被其他的goroutine访问到
	a.lk.Lock()
	a.instance = instance
	a.lk.Unlock()

	//if a.opts.rpcServer != nil {
	//	// 启动rpc服务， 如果我想要给这个rpc服务设置port 我们想要给这个rpc服务register我们自定义的interceptor
	//	a.opts.rpcServer.Serve()
	//}

	//重点， 写的很简单， http服务要启动
	//if a.opts.rpcServer != nil {
	//	err := a.opts.rpcServer.Start()
	//	if err != nil {
	//		return err
	//	}
	//}

	//现在启动了两个server，一个是restserver，一个是rpcserver
	/*
		这两个server是否必须同时启动成功？
		如果有一个启动失败，那么我们就要停止另外一个server
		如果启动了多个， 如果其中一个启动失败，其他的应该被取消
			如果剩余的server的状态：
				1. 还没有开始调用start
					stop
				2. start进行中
					调用进行中的cancel
				3. start已经完成
					调用stop
		如果我们的服务启动了然后这个时候用户立马进行了访问
	*/

	var servers []gs.Server
	if a.opts.restServer != nil {
		servers = append(servers, a.opts.restServer)
	}
	if a.opts.rpcServer != nil {
		servers = append(servers, a.opts.rpcServer)
	}

	ctx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel
	eg, ctx := errgroup.WithContext(ctx)
	wg := sync.WaitGroup{}
	for _, srv := range servers {
		// 启动server
		// 在启动一个goroutine去监听是否有err产生
		eg.Go(func() error {
			<-ctx.Done() //wait for stop signal
			//不可能无休止的等待stop
			sctx, scancel := context.WithTimeout(context.Background(), a.opts.stopTimeout)
			defer scancel()
			return srv.Stop(sctx)
		})

		wg.Add(1)
		eg.Go(func() error {
			wg.Done()
			log.Info("start rest server")
			return srv.Start(ctx)
		})
	}

	wg.Wait()

	// 注册服务
	if a.opts.registrar != nil {
		rctx, rcancel := context.WithTimeout(context.Background(), a.opts.registrarTimeout)
		defer rcancel()
		err = a.opts.registrar.Register(rctx, instance)
		if err != nil {
			log.Errorf("register service error: %v", err)
			return err
		}
	}

	// 监听退出信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, a.opts.sigs...)
	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-c:
			log.Info("received signal, stopping...")
			return a.Stop()
		}
	})
	if err = eg.Wait(); err != nil {
		log.Errorf("app error: %v", err)
		return err
	}

	return nil
}

/*
Stop 停止服务
http basic 认证
cache： 1. redis 2. memcache 3. local cache
jwt
*/
// 停止服务
func (a *App) Stop() error {
	a.lk.Lock()
	instance := a.instance
	a.lk.Unlock()

	log.Info("start deregister service")
	if a.opts.registrar != nil {
		rctx, rcancel := context.WithTimeout(context.Background(), a.opts.registrarTimeout)
		defer rcancel()
		if err := a.opts.registrar.Deregister(rctx, instance); err != nil {
			log.Errorf("deregister service error: %v", err)
			return err
		}
	}

	if a.cancel != nil {
		a.cancel()
	}

	return nil
}

// 创建服务注册结构体
func (a *App) buildInstance() (*registry.ServiceInstance, error) {
	endpoints := make([]string, 0)
	for _, e := range a.opts.endpoints {
		endpoints = append(endpoints, e.String())
	}

	// 从rpcserver,restserver 去获取信息
	if a.opts.rpcServer != nil {
		if a.opts.rpcServer.Endpoint() != "" {
			endpoints = append(endpoints, a.opts.rpcServer.Endpoint().String())
		} else {
			u := &url.URL{
				Scheme: "grpc",
				Host:   a.opts.rpcServer.Address(),
			}
			endpoints = append(endpoints, u.String())
		}
	}

	return &registry.ServiceInstance{
		ID:        a.opts.id,
		Name:      a.opts.name,
		Endpoints: endpoints,
	}, nil
}
