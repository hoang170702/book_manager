package impl

import (
	"book-manager/internal/dto/auth"
	"book-manager/internal/dto/common"
	"book-manager/internal/models/user"
	"book-manager/internal/repositories"
	"book-manager/internal/services"
	"book-manager/internal/utils"
	authutil "book-manager/internal/utils/auth"
	"book-manager/internal/utils/enums"
	"book-manager/internal/utils/enums/error_codes"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Repo repositories.IAuthRepository
}

func (s AuthService) Register(req *common.Request[auth.RegisterRequest]) common.Response[any] {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Data.Password), bcrypt.DefaultCost)
	if err != nil {
		return utils.BuildResponse[any](nil, error_codes.BadRequest, req.RequestId)
	}

	u := &user.User{
		Username: req.Data.Username,
		Password: string(hashedPassword),
	}
	u.Status = enums.StatusActive
	u.CreatedBy = req.Data.Username
	u.CreatedDate = time.Now()

	err = s.Repo.Create(u, req.RequestId)
	if err != nil {
		var appErr *error_codes.AppError
		if errors.As(err, &appErr) {
			return utils.BuildResponse[any](nil, error_codes.ErrorCode{
				Code: appErr.Code,
				Msg:  appErr.Message,
			}, req.RequestId)
		}
		return utils.BuildResponse[any](nil, error_codes.BadRequest, req.RequestId)
	}

	return utils.BuildResponse[any](nil, error_codes.Success, req.RequestId)
}

func (s AuthService) Login(req *common.Request[auth.LoginRequest]) common.Response[auth.LoginResponse] {
	// Find user by username
	u, err := s.Repo.FindByUsername(req.Data.Username, req.RequestId)
	if err != nil {
		var appErr *error_codes.AppError
		if errors.As(err, &appErr) {
			return utils.BuildResponse[auth.LoginResponse](auth.LoginResponse{}, error_codes.ErrorCode{
				Code: appErr.Code,
				Msg:  appErr.Message,
			}, req.RequestId)
		}
		return utils.BuildResponse[auth.LoginResponse](auth.LoginResponse{}, error_codes.BadRequest, req.RequestId)
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Data.Password))
	if err != nil {
		return utils.BuildResponse[auth.LoginResponse](auth.LoginResponse{}, error_codes.InvalidCredential, req.RequestId)
	}

	// Generate tokens
	accessToken, err := authutil.GenerateAccessToken(u.Id, u.Username)
	if err != nil {
		return utils.BuildResponse[auth.LoginResponse](auth.LoginResponse{}, error_codes.BadRequest, req.RequestId)
	}

	refreshToken, err := authutil.GenerateRefreshToken(u.Id, u.Username)
	if err != nil {
		return utils.BuildResponse[auth.LoginResponse](auth.LoginResponse{}, error_codes.BadRequest, req.RequestId)
	}

	return utils.BuildResponse[auth.LoginResponse](auth.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, error_codes.Success, req.RequestId)
}

func (s AuthService) RefreshToken(req *common.Request[auth.RefreshRequest]) common.Response[auth.LoginResponse] {
	// Validate refresh token
	claims, err := authutil.ValidateToken(req.Data.RefreshToken)
	if err != nil {
		return utils.BuildResponse[auth.LoginResponse](auth.LoginResponse{}, error_codes.InvalidToken, req.RequestId)
	}

	// Ensure it's a refresh token, not an access token
	if claims.Type != "refresh" {
		return utils.BuildResponse[auth.LoginResponse](auth.LoginResponse{}, error_codes.InvalidToken, req.RequestId)
	}

	// Generate new token pair
	accessToken, err := authutil.GenerateAccessToken(claims.UserId, claims.Username)
	if err != nil {
		return utils.BuildResponse[auth.LoginResponse](auth.LoginResponse{}, error_codes.BadRequest, req.RequestId)
	}

	refreshToken, err := authutil.GenerateRefreshToken(claims.UserId, claims.Username)
	if err != nil {
		return utils.BuildResponse[auth.LoginResponse](auth.LoginResponse{}, error_codes.BadRequest, req.RequestId)
	}

	return utils.BuildResponse[auth.LoginResponse](auth.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, error_codes.Success, req.RequestId)
}

func (s AuthService) Logout(accessToken string, refreshToken string, requestId string) common.Response[any] {
	// Blacklist access token
	if accessToken != "" {
		claims, err := authutil.ValidateToken(accessToken)
		if err == nil {
			_ = s.Repo.RevokeToken(accessToken, claims.ExpiresAt.Time, requestId)
		}
	}

	// Blacklist refresh token
	if refreshToken != "" {
		claims, err := authutil.ValidateToken(refreshToken)
		if err == nil {
			_ = s.Repo.RevokeToken(refreshToken, claims.ExpiresAt.Time, requestId)
		}
	}

	return utils.BuildResponse[any](nil, error_codes.Success, requestId)
}

func (s AuthService) IsTokenRevoked(token string) bool {
	return s.Repo.IsTokenRevoked(token)
}

func NewAuthService(repo repositories.IAuthRepository) services.IAuthService {
	return &AuthService{Repo: repo}
}
