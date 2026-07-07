package service

import (
	"action-camera/cache"
	"action-camera/dao"
	"action-camera/model"
	"action-camera/pkg/e"
	"action-camera/pkg/util"
	"action-camera/serizlizer"
	"context"
)

type UserService struct {
	Nickname  string `json:"nick_name" form:"nick_name"`
	Username  string `json:"user_name" form:"user_name"`
	Password  string `json:"password" form:"password"`
	Email     string `json:"email" form:"email"`
	EmailCode string `json:"code" form:"code"`
}

func (service *UserService) Register(ctx context.Context) serizlizer.Response {
	var user model.User
	code := e.Success
	if service.Username == "" {
		code = e.ErrorUsernameIsEmpty
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  e.GetMsg(code),
		}
	}
	if service.Password == "" {
		code = e.ErrorPasswordIsEmpty
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  e.GetMsg(code),
		}
	}
	if service.Email == "" {
		code = e.ErrorEmailIsEmpty
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  e.GetMsg(code),
		}
	}
	rdb := cache.NewRDB(ctx)
	rdbCode, err := rdb.GetVerificationCode(service.Email)
	if err != nil {
		code = e.ErrorRedisGetkey
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	if rdbCode != service.EmailCode {
		code = e.ErrorRedisGetkey
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  "验证码不匹配",
		}
	}
	userDao := dao.NewUserDao(ctx)
	_, exist, err := userDao.ExistOrNotByUserName(service.Username)
	if err != nil {
		code = e.Error
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	if exist {
		code = e.ErrorExistUser
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  "该用户已经存在",
		}
	}
	user = model.User{
		Nickname:      service.Nickname,
		Username:      service.Username,
		Email:         service.Email,
		EmailVerified: true,
	}
	//密码加密
	if err := user.SetPassword(service.Password); err != nil {
		code = e.ErrorFailEncryption
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	if err := userDao.CreateUser(&user); err != nil {
		code = e.Error
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serizlizer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

func (service *UserService) Login(ctx context.Context) serizlizer.Response {
	code := e.Success
	userDao := dao.NewUserDao(ctx)
	user, exist, err := userDao.ExistOrNotByUserName(service.Username)
	if err != nil {
		code = e.Error
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	if !exist {
		code = e.ErrorNotExistUser
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  e.GetMsg(code),
		}
	}
	if !user.CheckPassword(service.Password) {
		code = e.ErrorPassword
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  e.GetMsg(code),
		}
	}
	token, err := util.GenerateToken(user.ID, user.Username, 0)
	if err != nil {
		code = e.ErrorGenerateToken
		return serizlizer.Response{
			Status: code,
		}
	}
	return serizlizer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data: serizlizer.TokenData{
			Token: token,
			User:  serizlizer.BuildUser(user),
		},
	}
}
