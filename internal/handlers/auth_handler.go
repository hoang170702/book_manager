package handlers

import (
	"book-manager/internal/constants"
	"book-manager/internal/dto/auth"
	"book-manager/internal/dto/common"
	"book-manager/internal/services"
	"book-manager/internal/utils"
	"book-manager/internal/utils/enums/error_codes"
	"strings"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	Service services.IAuthService
}

func NewAuthHandler(service services.IAuthService) *AuthHandler {
	return &AuthHandler{Service: service}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var reqDto common.Request[auth.RegisterRequest]

	if err := c.Bind(&reqDto); err != nil {
		resp := utils.BuildResponse[any](nil, error_codes.InvalidRequest, "")
		return c.JSON(constants.StatusBadRequest, resp)
	}

	if err := c.Validate(&reqDto.Data); err != nil {
		resp := utils.BuildResponse[any](nil, error_codes.InvalidRequest, reqDto.RequestId)
		return c.JSON(constants.StatusBadRequest, resp)
	}

	resp := h.Service.Register(&reqDto)
	return c.JSON(constants.StatusOK, resp)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var reqDto common.Request[auth.LoginRequest]

	if err := c.Bind(&reqDto); err != nil {
		resp := utils.BuildResponse[any](nil, error_codes.InvalidRequest, "")
		return c.JSON(constants.StatusBadRequest, resp)
	}

	if err := c.Validate(&reqDto.Data); err != nil {
		resp := utils.BuildResponse[any](nil, error_codes.InvalidRequest, reqDto.RequestId)
		return c.JSON(constants.StatusBadRequest, resp)
	}

	resp := h.Service.Login(&reqDto)
	return c.JSON(constants.StatusOK, resp)
}

func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var reqDto common.Request[auth.RefreshRequest]

	if err := c.Bind(&reqDto); err != nil {
		resp := utils.BuildResponse[any](nil, error_codes.InvalidRequest, "")
		return c.JSON(constants.StatusBadRequest, resp)
	}

	if err := c.Validate(&reqDto.Data); err != nil {
		resp := utils.BuildResponse[any](nil, error_codes.InvalidRequest, reqDto.RequestId)
		return c.JSON(constants.StatusBadRequest, resp)
	}

	resp := h.Service.RefreshToken(&reqDto)
	return c.JSON(constants.StatusOK, resp)
}

func (h *AuthHandler) Logout(c echo.Context) error {
	var reqDto common.Request[auth.RefreshRequest]

	if err := c.Bind(&reqDto); err != nil {
		resp := utils.BuildResponse[any](nil, error_codes.InvalidRequest, "")
		return c.JSON(constants.StatusBadRequest, resp)
	}

	// Extract access token from Authorization header
	accessToken := ""
	authHeader := c.Request().Header.Get("Authorization")
	if parts := strings.SplitN(authHeader, " ", 2); len(parts) == 2 {
		accessToken = parts[1]
	}

	// Refresh token from request body
	refreshToken := reqDto.Data.RefreshToken

	resp := h.Service.Logout(accessToken, refreshToken, reqDto.RequestId)
	return c.JSON(constants.StatusOK, resp)
}
