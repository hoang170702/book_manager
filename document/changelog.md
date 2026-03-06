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
