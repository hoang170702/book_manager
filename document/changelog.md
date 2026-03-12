# Changelog — Book Manager Project

## Completed Fixes

### Fix #1: Global DB Variable + `pkg/` import `internal/` ✅
**Vấn đề:** `pkg/database/db.go` dùng biến global `var DB *gorm.DB` và import `internal/models` → vi phạm Go convention.

**Đã sửa:**
- `pkg/database/db.go` — Xóa global `DB`, `Connect()` return `*gorm.DB`. Xóa toàn bộ import `internal/`.
- `internal/config/migration.go` — **[NEW]** Chuyển logic migration từ `pkg/database` vào đây, nhận `*gorm.DB` qua parameter.
- `cmd/api/main.go` — `db := database.Connect()`, truyền `db` xuống dependency chain.
- `internal/handlers/health_handler.go` — Dùng closure nhận `*gorm.DB`, không dùng global.

### GORM Logger + Request ID Tracing ✅
**Vấn đề:** GORM log SQL mặc định rất ồn, không có request ID để trace.

**Đã sửa:**
- `internal/utils/logger/gorm_logger.go` — **[NEW]** Custom GORM logger: chỉ log **slow query >200ms** (Warn), kèm Request ID. Không auto-log errors (để repository tự quyết).
- `pkg/database/db.go` — Dùng custom `GormLogger` thay GORM default.
- `internal/repositories/category_repository.go` — Thêm `dbCtx()` inject Request ID vào mỗi query.
- `internal/repositories/author_repository.go` — Tương tự.
- `internal/utils/response_config.go` — `BuildResponse` tự log WARN khi response code ≠ success.

### HTTP Logging Middleware ✅
**Vấn đề:** Log chỉ hiện middleware UUID, không hiện business `request_id` từ body. Không log request/response body.

**Đã sửa:**
- `internal/middleware/logging.go` — Đọc request body, extract `request_id` từ JSON body, capture response body. Ưu tiên body `request_id`, fallback UUID khi không có.

**Format log:**
```
[INFO] [HTTP] [1212-12121-12121-12121] --> POST /api/v1/categories/add | body: {...}
[INFO] [HTTP] [1212-12121-12121-12121] <-- POST /api/v1/categories/add | 200 | 5ms | body: {...}
```

### HTTP Status Constants ✅
**Đã sửa:**
- `internal/constants/http_status.go` — **[NEW]** HTTP status constants (`StatusOK`, `StatusBadRequest`, ...).
- Tất cả handlers + `recovery.go` đã đổi từ hardcoded → constants.

### Fix #2: Repository Interface ✅
**Vấn đề:** Service depend vào concrete struct `*repositories.CategoryRepository` → không mock được.

**Đã sửa:**
- `internal/repositories/category_repository_interface.go` — **[NEW]** `ICategoryRepository` interface.
- `internal/repositories/author_repository_interface.go` — **[NEW]** `IAuthorRepository` interface.
- `internal/services/impl/category_service_impl.go` — `Repo` field đổi sang `repositories.ICategoryRepository`.
- `internal/services/impl/author_service_impl.go` — `Repo` field đổi sang `repositories.IAuthorRepository`.

### Fix #3: Validation Middleware ✅
**Vấn đề:** DTO có `validate:"required"` tag nhưng không có validator nào được gọi. `c.Bind()` chỉ unmarshal JSON, **không validate**.

**Đã sửa:**
- `go.mod` — Thêm dependency `github.com/go-playground/validator/v10`
- `internal/middleware/validator.go` — **[NEW]** Custom validator wraps `go-playground/validator` cho Echo.
- `cmd/api/main.go` — Gọi `middleware.RegisterValidator(e)` trước khi setup routes.
- `internal/handlers/category_handler.go` — Thêm `c.Validate(&reqDto.Data)` sau `c.Bind()` ở Create, GetOne, Update, Delete.
- `internal/handlers/AuthorHandler.go` — Tương tự cho tất cả methods.

### Fix #4: Error Handling Consistency ✅
**Vấn đề:** `category_service_impl.go` dùng `err.(*AppError)` (type assertion) → panic nếu error khác type.

**Đã sửa:**
- `internal/services/impl/category_service_impl.go` — Đổi tất cả 5 chỗ type assertion thành `errors.As(err, &appErr)` + fallback `BadRequest` cho unexpected errors. Giờ thống nhất với `author_service_impl.go`.

### Fix #5: DTO Improvements ✅
**Vấn đề:** Service trả `models.Category`/`models.Author` trực tiếp → expose DB fields (status, created_by...). `Request` import `models.Paginate` → DTO phụ thuộc Model.

**Đã sửa:**
- `dto/common/paginate.go` — **[NEW]** Copy Paginate từ models vào dto layer.
- `dto/common/request.go` — Dùng `common.Paginate` thay `models.Paginate`, xóa import models.
- `dto/category/category_response.go` — **[NEW]** Response DTO chỉ expose `id`, `name`.
- `dto/author/author_response.go` — **[NEW]** Response DTO chỉ expose `id`, `name`.
- `services/category_service.go` — `GetOne`/`GetAll` return `CategoryResponse` thay `models.Category`.
- `services/author_service.go` — Tương tự cho `AuthorResponse`.
- `services/impl/category_service_impl.go` — Map entity → response DTO.
- `services/impl/author_service_impl.go` — Tương tự.

### Fix #6: CORS Middleware ✅
**Vấn đề:** Không có CORS middleware → browser từ domain khác gọi API bị block.

**Đã sửa:**
- `cmd/api/main.go` — Thêm `echomw.CORSWithConfig()` vào `setupMiddleware()`. Cho phép `AllowOrigins: *`, các method GET/POST/PUT/DELETE/OPTIONS, và headers Origin/Content-Type/Accept/Authorization. CORS đặt đầu tiên để xử lý preflight OPTIONS trước các middleware khác.

### Fix #7: File Naming Consistency ✅
**Vấn đề:** File `AuthorHandler.go` viết hoa PascalCase không đúng chuẩn Go.

**Đã sửa:**
- Đổi tên file `internal/handlers/AuthorHandler.go` thành `internal/handlers/author_handler.go` để đồng bộ chuẩn snake_case.

### Fix #8: Hardcoded "Anonymous" User ✅
**Vấn đề:** Hardcode `"Anonymous"` ở 10+ chỗ khi gọi Repository `Create/Update/Delete`.

**Đã sửa:**
- `internal/utils/auth/user.go` — **[NEW]** Tạo helper `GetCurrentUser(c echo.Context) string` trả về user (hiện tại tạm gán "Anonymous", sau này lấy từ JWT).
- `internal/handlers/author_handler.go`, `category_handler.go` — Đọc user từ context và truyền xuống Service layer.
- Thay đổi `Create/Update/Delete` trong `Service` interface và implementation để nhận parameter `user string` thay vì hardcode. Dữ liệu này được truyền trọn vẹn xuống Repository.

### Fix #9: Clean Up Empty Packages ✅
**Vấn đề:** Các file `dto/auth/login.go`, `models/user/role.go`, `models/user/user.go` chỉ chứa duy nhất dòng khai báo `package` gây khó hiểu, không rõ là đã xóa logic hay chưa implement.

**Đã sửa:**
- Thêm comment kiểu `// TODO: Implement...` vào trong 3 file rỗng đó để đánh dấu rõ ràng đây là các file chờ được implement logic sau này thay vì xóa đi cấu trúc thư mục.
