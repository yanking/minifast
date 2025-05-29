package serverinterceptors

import (
	"context"
	"google.golang.org/grpc"
	"minifast/pkg/log"
	"runtime/debug"
)

func UnaryCrashInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	defer handleCrash(func(r interface{}) {
		log.Errorf("%+v\n \n %s", r, debug.Stack())
	})
	return handler(ctx, req)
}

func handleCrash(f func(interface{})) {
	if r := recover(); r != nil {
		f(r)
	}
}

func StreamCrashInterceptor(svr interface{}, stream grpc.ServerStream, _ *grpc.StreamServerInfo,
	handler grpc.StreamHandler) (err error) {
	defer handleCrash(func(r interface{}) {
		log.Errorf("%+v\n \n %s", r, debug.Stack())
	})

	return handler(svr, stream)
}
