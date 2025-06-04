package v1

import (
	"context"
	"minifast/app/pkg/code"
	dv1 "minifast/app/user/srv/data/v1"
	"minifast/pkg/errors"
)

type UserDTO struct {
	dv1.UserDO
}

type UserSrv interface {
	Create(ctx context.Context, user *UserDTO) error
	//List(ctx context.Context, orderby []string, opts metav1.ListMeta) (*UserDTOList, error)
	//Update(ctx context.Context, user *UserDTO) error
	//GetByID(ctx context.Context, ID uint64) (*UserDTO, error)
	//GetByMobile(ctx context.Context, mobile string) (*UserDTO, error)
}

type userService struct {
	userStrore dv1.UserStore
}

func NewUserService(us dv1.UserStore) UserSrv {
	return &userService{
		userStrore: us,
	}
}

var _ UserSrv = &userService{}

func (u *userService) Create(ctx context.Context, user *UserDTO) error {
	//先判断用户是否存在
	_, err := u.userStrore.GetByMobile(ctx, user.Mobile)
	if err != nil && errors.IsCode(err, code.ErrUserNotFound) {
		return u.userStrore.Create(ctx, &user.UserDO)
	}

	//这里应该区别到底是什么错误，用户已经存在？ 数据访问错误？
	return errors.WithCode(code.ErrUserAlreadyExists, "用户已经存在")
}
