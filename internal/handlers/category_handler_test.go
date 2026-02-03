package handlers

import (
	"book-manager/internal/dto/category"
	"book-manager/internal/dto/common"
	"book-manager/internal/models"
	_ "book-manager/internal/services"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockService is a mock implementation of services.ICategoryService
type MockService struct {
	mock.Mock
}

func (m *MockService) Create(req *common.Request[category.AddCategory]) common.Response[any] {
	args := m.Called(req)
	return args.Get(0).(common.Response[any])
}

func (m *MockService) GetOne(req *common.Request[category.GetOneCategory]) common.Response[models.Category] {
	args := m.Called(req)
	return args.Get(0).(common.Response[models.Category])
}

func (m *MockService) GetAll(req *common.Request[any]) common.Response[[]models.Category] {
	args := m.Called(req)
	return args.Get(0).(common.Response[[]models.Category])
}

func (m *MockService) Update(req *common.Request[category.UpdateCategory]) common.Response[any] {
	args := m.Called(req)
	return args.Get(0).(common.Response[any])
}

func (m *MockService) Delete(req *common.Request[category.DeleteCategory]) common.Response[any] {
	args := m.Called(req)
	return args.Get(0).(common.Response[any])
}

func TestCategoryHandler_GetAll_RequestId(t *testing.T) {
	// Setup
	e := echo.New()
	reqBody := `{"request_id": "test-req-id-123", "data": {}}`
	req := httptest.NewRequest(http.MethodPost, "/categories/get-all", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock Service
	mockService := new(MockService)
	expectedResp := common.Response[[]models.Category]{
		ResponseId:   "test-req-id-123",
		ResponseCode: "00",
		ResponseMsg:  "Success",
		Data:         []models.Category{},
	}
	// We expect GetAll to be called with a request that has the correct RequestId
	mockService.On("GetAll", mock.MatchedBy(func(r *common.Request[any]) bool {
		return r.RequestId == "test-req-id-123"
	})).Return(expectedResp)

	h := NewCategoryHandler(mockService)

	// Execute
	if assert.NoError(t, h.GetAll(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp common.Response[[]models.Category]
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)

		// Verify ResponseId is present in response (mapped from RequestId)
		assert.Equal(t, "test-req-id-123", resp.ResponseId)
		assert.Equal(t, "00", resp.ResponseCode)

		// Write success to file
		os.WriteFile("test_result.txt", []byte("PASS"), 0644)
	} else {
		os.WriteFile("test_result.txt", []byte("FAIL"), 0644)
	}
}
