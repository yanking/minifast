package serverinterceptors

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

func UnaryTimeoutInterceptor(timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var resp interface{}
		var err error
		var lock sync.Mutex
		done := make(chan struct{})
		panicChan := make(chan interface{}, 1)
		go func() {
			defer func() {
				if p := recover(); p != nil {
					// attach call stack to avoid missing in different goroutine
					panicChan <- fmt.Sprintf("%+v\n\n%s", p, strings.TrimSpace(string(debug.Stack())))
				}
			}()

			lock.Lock()
			defer lock.Unlock()
			resp, err = handler(ctx, req)
			close(done)
		}()

		select {
		case p := <-panicChan:
			// panic
			panic(p)
		case <-done:
			lock.Lock()
			defer lock.Unlock()
			return resp, err
		case <-ctx.Done():
			err = ctx.Err()
			if errors.Is(err, context.Canceled) {
				err = status.Error(codes.Canceled, err.Error())
			} else if errors.Is(err, context.DeadlineExceeded) {
				err = status.Error(codes.DeadlineExceeded, err.Error())
			}
			// timeout
			return nil, err
		}
	}
}
