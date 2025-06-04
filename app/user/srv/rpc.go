package srv

import (
	"fmt"
	"github.com/alibaba/sentinel-golang/pkg/adapters/grpc"
	upb "minifast/api/user/v1"
	"minifast/app/pkg/options"
	"minifast/gmicro/core/trace"
	"minifast/gmicro/server/rpcserver"
)

func NewUserRPCServer(telemetry *options.TelemetryOptions, serverOpts *options.ServerOptions, userver upb.UserServer) (*rpcserver.Server, error) {
	// 初始化open-telemetry的exporter
	trace.InitAgent(trace.Options{
		Name:     telemetry.Name,
		Endpoint: telemetry.Endpoint,
		Sampler:  telemetry.Sampler,
		Batcher:  telemetry.Batcher,
	})

	rpcAddr := fmt.Sprintf("%s:%d", serverOpts.Host, serverOpts.Port)

	var opts []rpcserver.ServerOption
	opts = append(opts, rpcserver.WithAddress(rpcAddr))
	if serverOpts.EnableLimit {
		opts = append(opts, rpcserver.WithUnaryInterceptor(grpc.NewUnaryServerInterceptor()))
	}
	urpcServer := rpcserver.NewServer(opts...)

	upb.RegisterUserServer(urpcServer.Server, userver)

	return urpcServer, nil
}
