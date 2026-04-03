package services

import (
	"book-manager/internal/dto/auth"
	"book-manager/internal/dto/common"
)

type IAuthService interface {
	Register(req *common.Request[auth.RegisterRequest]) common.Response[any]
	Login(req *common.Request[auth.LoginRequest]) common.Response[auth.LoginResponse]
	RefreshToken(req *common.Request[auth.RefreshRequest]) common.Response[auth.LoginResponse]
}
