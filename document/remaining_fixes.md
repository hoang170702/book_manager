# Remaining Fixes — Plan

## Fix #3: Validation Middleware 🟡
**Vấn đề:** DTO có `validate:"required"` tag nhưng không có validator nào được gọi. `c.Bind()` chỉ unmarshal JSON, **không validate**.

**Plan:**
- Cài thêm `go-playground/validator/v10`
- Tạo `internal/middleware/validator.go` — custom validator cho Echo
- Gắn validator vào Echo instance trong `main.go`
- Handler gọi `c.Validate()` sau `c.Bind()` để validate DTO

**Files thay đổi:** `go.mod`, `middleware/validator.go` (NEW), `main.go`, tất cả handlers

---

## Fix #4: Error Handling Consistency 🟡
**Vấn đề:** `CategoryService` dùng type assertion `err.(*error_codes.AppError)` → **panic nếu error khác type**. `AuthorService` dùng `errors.As()` → an toàn hơn.

**Plan:**
- Thống nhất tất cả service dùng `errors.As(err, &appErr)` thay vì type assertion
- Chỉ sửa `category_service_impl.go` (Author đã đúng)

**Files thay đổi:** `services/impl/category_service_impl.go`

---

## Fix #5: DTO Improvements 🟡
**Vấn đề:**
- Service trả `Response[models.Category]` → expose entity ra ngoài API
- `Request` import `models.Paginate` → DTO phụ thuộc Model layer

**Plan:**
- Tạo response DTO riêng cho từng domain (e.g. `dto/category/category_response.go`)
- Tạo `dto/common/paginate.go` thay thế `models.Paginate` trong Request
- Service trả response DTO thay vì entity model

**Files thay đổi:** `dto/` (NEW files), services, handlers

---

## Fix #6: CORS Middleware 🟡
**Vấn đề:** Thiếu CORS middleware → frontend từ domain khác không gọi được API.

**Plan:**
- Dùng Echo built-in `middleware.CORS()` hoặc custom config
- Thêm vào `setupMiddleware()` trong `main.go`

**Files thay đổi:** `cmd/api/main.go`

---

## Fix #7: File Naming Consistency 🟢
**Vấn đề:** `AuthorHandler.go` (PascalCase) vs `category_handler.go` (snake_case).

**Plan:**
- Đổi tên `AuthorHandler.go` → `author_handler.go` (Go convention: snake_case)

**Files thay đổi:** rename 1 file

---

## Fix #8: Hardcoded "Anonymous" User 🟢
**Vấn đề:** `Create`/`Update`/`Delete` hardcode `"Anonymous"` → khi có auth sẽ cần lấy từ context.

**Plan:**
- Tạo helper `GetCurrentUser(c echo.Context) string` → trả "Anonymous" mặc định
- Handler truyền user xuống service
- Sau này khi có auth middleware, chỉ cần sửa helper

**Files thay đổi:** `utils/` (NEW helper), handlers, services

---

## Fix #9: Clean Up Empty Packages 🟢
**Vấn đề:** `constants/`, `dto/auth/login.go`, `models/user/` chứa file rỗng.

**Plan:**
- Xóa file rỗng hoặc thêm TODO comment
- Giữ lại folder structure nếu có kế hoạch dùng
