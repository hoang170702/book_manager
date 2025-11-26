package handlers

import (
	"book-manager/internal/dto/author"
	"book-manager/internal/dto/common"
	"book-manager/internal/services"
	"book-manager/internal/utils"
	"book-manager/internal/utils/enums/error_codes"

	"github.com/labstack/echo/v4"
)

type AuthorHandler struct {
	Service services.IAuthorService
}

func NewAuthorHandler(service services.IAuthorService) *AuthorHandler {
	return &AuthorHandler{Service: service}
}

func (h *AuthorHandler) Add(c echo.Context) error {
	var reqDto common.Request[author.AddAuthor]

	if err := c.Bind(&reqDto); err != nil {
		resp := utils.BuildResponse[any](nil, error_codes.InvalidRequest)
		return c.JSON(400, resp)
	}

	resp := h.Service.Create(&reqDto)
	return c.JSON(200, resp)
}

func (h *AuthorHandler) GetOne(c echo.Context) error {
	var reqDto common.Request[author.GetOneAuthor]

	if err := c.Bind(&reqDto); err != nil {
		resp := utils.BuildResponse[any](nil, error_codes.InvalidRequest)
		return c.JSON(400, resp)
	}

	resp := h.Service.GetOne(&reqDto)
	return c.JSON(200, resp)
}

func (h *AuthorHandler) GetAll(c echo.Context) error {
	resp := h.Service.GetAll(&common.Request[any]{})
	return c.JSON(200, resp)
}

func (h *AuthorHandler) Update(c echo.Context) error {
	var reqDto common.Request[author.UpdateAuthor]
	if err := c.Bind(&reqDto); err != nil {
		resp := utils.BuildResponse[any](nil, error_codes.InvalidRequest)
		return c.JSON(400, resp)
	}
	resp := h.Service.Update(&reqDto)
	return c.JSON(200, resp)
}
