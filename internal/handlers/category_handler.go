package handlers

import (
	"book-manager/internal/constants"
	"book-manager/internal/dto/category"
	"book-manager/internal/dto/common"
	"book-manager/internal/services"
	"book-manager/internal/utils"
	"book-manager/internal/utils/auth"
	"book-manager/internal/utils/enums/error_codes"

	"github.com/labstack/echo/v4"
)

type CategoryHandler struct {
	Service services.ICategoryService
}

func NewCategoryHandler(service services.ICategoryService) *CategoryHandler {
	return &CategoryHandler{Service: service}
}

func (h *CategoryHandler) Create(c echo.Context) error {
	var reqDto common.Request[category.AddCategory]

	if err := c.Bind(&reqDto); err != nil {
		resp := utils.BuildResponse[any](nil, error_codes.InvalidRequest, "")
		return c.JSON(constants.StatusBadRequest, resp)
	}

	if err := c.Validate(&reqDto.Data); err != nil {
		resp := utils.BuildResponse[any](nil, error_codes.InvalidRequest, reqDto.RequestId)
		return c.JSON(constants.StatusBadRequest, resp)
	}

	user := auth.GetCurrentUser(c)
	resp := h.Service.Create(&reqDto, user)
	return c.JSON(constants.StatusOK, resp)
}

func (h *CategoryHandler) GetOne(c echo.Context) error {
	var reqDto common.Request[category.GetOneCategory]

	if err := c.Bind(&reqDto); err != nil {
		resp := utils.BuildResponse[any](nil, error_codes.InvalidRequest, "")
		return c.JSON(constants.StatusBadRequest, resp)
	}

	if err := c.Validate(&reqDto.Data); err != nil {
		resp := utils.BuildResponse[any](nil, error_codes.InvalidRequest, reqDto.RequestId)
		return c.JSON(constants.StatusBadRequest, resp)
	}

	resp := h.Service.GetOne(&reqDto)
	return c.JSON(constants.StatusOK, resp)
}

func (h *CategoryHandler) GetAll(c echo.Context) error {
	var reqDto common.Request[any]
	if err := c.Bind(&reqDto); err != nil {
		resp := utils.BuildResponse[any](nil, error_codes.InvalidRequest, "")
		return c.JSON(constants.StatusBadRequest, resp)
	}
	resp := h.Service.GetAll(&reqDto)
	return c.JSON(constants.StatusOK, resp)
}

func (h *CategoryHandler) Update(c echo.Context) error {
	var reqDto common.Request[category.UpdateCategory]

	if err := c.Bind(&reqDto); err != nil {
		resp := utils.BuildResponse[any](nil, error_codes.InvalidRequest, "")
		return c.JSON(constants.StatusBadRequest, resp)
	}

	if err := c.Validate(&reqDto.Data); err != nil {
		resp := utils.BuildResponse[any](nil, error_codes.InvalidRequest, reqDto.RequestId)
		return c.JSON(constants.StatusBadRequest, resp)
	}

	user := auth.GetCurrentUser(c)
	resp := h.Service.Update(&reqDto, user)
	return c.JSON(constants.StatusOK, resp)
}

func (h *CategoryHandler) Delete(c echo.Context) error {
	var reqDto common.Request[category.DeleteCategory]

	if err := c.Bind(&reqDto); err != nil {
		resp := utils.BuildResponse[any](nil, error_codes.InvalidRequest, "")
		return c.JSON(constants.StatusBadRequest, resp)
	}

	if err := c.Validate(&reqDto.Data); err != nil {
		resp := utils.BuildResponse[any](nil, error_codes.InvalidRequest, reqDto.RequestId)
		return c.JSON(constants.StatusBadRequest, resp)
	}

	user := auth.GetCurrentUser(c)
	resp := h.Service.Delete(&reqDto, user)
	return c.JSON(constants.StatusOK, resp)
}
