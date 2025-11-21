package handlers

import (
	"book-manager/internal/dto/category"
	"book-manager/internal/models/common"
	"book-manager/internal/services"
	"book-manager/internal/utils"
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
		resp := utils.BuildResponse[any](nil, error_codes.InvalidRequest)
		return c.JSON(400, resp)
	}

	resp := h.Service.Create(&reqDto)
	return c.JSON(200, resp)
}
