package user

import (
	"github.com/google/wire"
	v1 "minifast/api/user/v1"
	srv1 "minifast/app/user/srv/service/v1"
)

var ProviderSet = wire.NewSet(NewUserServer)

type userServer struct {
	v1.UnimplementedUserServer
	srv srv1.UserSrv
}

// java中的ioc，控制翻转 ioc = injection of control
// 代码分层，第三方服务， rpc， redis， 等等， 带来一定的复杂度
func NewUserServer(srv srv1.UserSrv) v1.UserServer {
	return &userServer{srv: srv}
}

var _ v1.UserServer = &userServer{}
