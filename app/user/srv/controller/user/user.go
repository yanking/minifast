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

func NewUserServer(srv srv1.UserSrv) v1.UserServer {
	return &userServer{srv: srv}
}

var _ v1.UserServer = &userServer{}
