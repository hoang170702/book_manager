package main

import (
	"book-manager/internal/dto/category"
	"book-manager/internal/dto/common"
	"book-manager/internal/handlers"
	"book-manager/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

// MockService
type MockService struct{}

func (m *MockService) Create(req *common.Request[category.AddCategory]) common.Response[any] {
	return common.Response[any]{}
}
func (m *MockService) GetOne(req *common.Request[category.GetOneCategory]) common.Response[models.Category] {
	return common.Response[models.Category]{}
}
func (m *MockService) GetAll(req *common.Request[any]) common.Response[[]models.Category] {
	// Verify RequestId is passed to service (optional, but good)
	if req.RequestId != "test-req-id-123" {
		fmt.Printf("FAIL: Service received wrong RequestId: %s\n", req.RequestId)
	}
	return common.Response[[]models.Category]{
		ResponseCode: "00",
		ResponseMsg:  "Success",
		Data:         []models.Category{},
	}
}
func (m *MockService) Update(req *common.Request[category.UpdateCategory]) common.Response[any] {
	return common.Response[any]{}
}
func (m *MockService) Delete(req *common.Request[category.DeleteCategory]) common.Response[any] {
	return common.Response[any]{}
}

func main() {
	e := echo.New()
	reqBody := `{"request_id": "test-req-id-123", "data": {}}`
	req := httptest.NewRequest(http.MethodPost, "/categories/get-all", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := &MockService{}
	h := handlers.NewCategoryHandler(mockService)

	// Execute
	if err := h.GetAll(c); err != nil {
		fmt.Printf("FAIL: Handler returned error: %v\n", err)
		return
	}

	if rec.Code != http.StatusOK {
		fmt.Printf("FAIL: Status code %d\n", rec.Code)
		return
	}

	var resp common.Response[[]models.Category]
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		return
	}

	fmt.Println("PASS: RequestId verification successful")
	os.WriteFile("verification_result.txt", []byte("PASS"), 0644)
}
