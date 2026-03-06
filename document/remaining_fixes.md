# All Fixes — Chi tiết giải thích (#1 → #9)

---

## Fix #1: Global DB Variable + `pkg/` import `internal/` ✅ ĐÃ FIX

### Vấn đề gì?

**1. Biến global `var DB *gorm.DB`:**
```go
// pkg/database/db.go
var DB *gorm.DB  // ← bất kỳ file nào cũng truy cập được, không kiểm soát

func Connect() {
    DB, err = gorm.Open(...)  // ← set global
}
```

**2. `pkg/` import `internal/`:**
```go
// pkg/database/db.go
import "book-manager/internal/models"  // ← vi phạm Go convention
```

### Vì sao gây lỗi?

1. **Global variable**: Bất kỳ goroutine nào cũng đọc/ghi `database.DB` → race condition. Không control được lifecycle (ai khởi tạo, ai close)
2. **Go convention**: `pkg/` là thư mục public (có thể import bởi project khác), `internal/` là private. `pkg` import `internal` = public phụ thuộc vào private → sai kiến trúc
3. **Khó test**: Unit test không thể thay DB giả vào, vì code truy cập global trực tiếp

### Đã fix gì?

- `pkg/database/db.go` — Xóa `var DB`, `Connect()` return `*gorm.DB`
- `internal/config/migration.go` — **[NEW]** Chuyển migration logic vào internal
- `cmd/api/main.go` — `db := database.Connect()`, truyền `db` qua DI
- `internal/handlers/health_handler.go` — Nhận `*gorm.DB` qua closure
- Custom GORM Logger — Warn level + Request ID tracing cho slow query
- `internal/utils/response_config.go` — Auto WARN log khi response ≠ success

---

## Fix #2: Repository Interface ✅ ĐÃ FIX

### Vấn đề gì?

```go
type CategoryService struct {
    Repo *repositories.CategoryRepository  // ← concrete struct
}
```

### Vì sao gây lỗi?

1. **Không mock được**: Unit test service phải kết nối DB thật → chậm, fragile
2. **Vi phạm Dependency Inversion**: High-level (service) phụ thuộc low-level (concrete repo)
3. **Coupling**: Thay đổi repo implementation → phải sửa service

### Đã fix gì?

- Tạo `ICategoryRepository`, `IAuthorRepository` interfaces
- Service depend vào interface thay concrete struct
- `NewCategoryService(repo ICategoryRepository)` — nhận interface

---

## Fix #3: Validation Middleware ✅ ĐÃ FIX

### Vấn đề gì?

```go
type AddCategory struct {
    Name string `json:"name" validate:"required"`  // ← tag có nhưng không ai gọi
}

// Handler:
c.Bind(&reqDto)  // ← chỉ unmarshal JSON, KHÔNG validate
```

### Vì sao gây lỗi?

1. Client gửi `{"data": {}}` (thiếu `name`) → server **vẫn chấp nhận** → tạo record thiếu data
2. Validate tag tồn tại nhưng vô nghĩa → đánh lừa developer nghĩ "đã validate"
3. Bad data vào DB → gây lỗi cascade ở các query/report sau

### Đã fix gì?

- Tạo `middleware/validator.go` wraps `go-playground/validator`
- Register validator vào Echo instance
- Handler gọi `c.Validate(&reqDto.Data)` sau `c.Bind()`

---

## Fix #4: Error Handling Consistency ✅ ĐÃ FIX

### Vấn đề gì?

**CategoryService** dùng **type assertion** (NGUY HIỂM):
```go
if err != nil {
    appErr := err.(*error_codes.AppError)  // ← nếu err KHÔNG phải AppError → PANIC!
    return utils.BuildResponse[any](nil, ...)
}
```

**AuthorService** dùng **errors.As** (AN TOÀN):
```go
if err != nil {
    var appErr *error_codes.AppError
    if errors.As(err, &appErr) {           // ← kiểm tra trước, không panic
        return utils.BuildResponse[any](nil, ...)
    }
    return utils.BuildResponse[any](nil, error_codes.BadRequest, ...)  // fallback
}
```

### Vì sao gây lỗi?

1. Hiện tại repository **chỉ trả `*AppError`**, nên chưa crash. Nhưng nếu sau này thêm middleware/wrapper bọc lại error → error type thay đổi → **PANIC trên production**
2. GORM cũng trả error khác type (connection timeout, deadlock) → **PANIC**
3. Go team khuyến nghị luôn dùng `errors.As()` thay vì type assertion cho error

### Fix gì?

Đổi `err.(*error_codes.AppError)` → `errors.As(err, &appErr)` trong `category_service_impl.go` (5 chỗ).

**Scope:** 1 file

---

## Fix #5: DTO Improvements 🟡 CHƯA FIX

### Vấn đề gì?

**1. Service trả entity model ra ngoài API:**
```go
GetOne(...) common.Response[models.Category]  // ← trả thẳng entity DB ra client
```
Client nhận được cả `status`, `created_by`, `updated_by` — field nội bộ DB.

**2. DTO import Model:**
```go
type Request[T any] struct {
    Paginate *models.Paginate  // ← DTO phụ thuộc Model layer
}
```

### Vì sao phải fix?

1. **Bảo mật**: Expose entity = lộ cấu trúc DB
2. **Coupling**: Thay đổi DB schema → breaking change cho client
3. **Clean Architecture**: DTO layer nên **độc lập** với Model layer

### Fix gì?

- Tạo response DTO riêng (e.g. `CategoryResponse`)
- Tạo `dto/common/paginate.go` thay `models.Paginate`
- Service map entity → response DTO

**Scope:** Tạo 3-4 file mới, sửa services + interfaces

---

## Fix #6: CORS Middleware 🟡 CHƯA FIX

### Vấn đề gì?

Không có CORS middleware → browser từ domain khác gọi API bị block.

### Vì sao phải fix?

- Frontend `localhost:3000` gọi API `localhost:8091` → browser chặn
- Production deploy frontend và API khác domain → **bắt buộc cần CORS**

### Fix gì?

Thêm Echo built-in CORS middleware vào `main.go`. 2-3 dòng code.

**Scope:** 1 file: `main.go`

---

## Fix #7: File Naming Consistency 🟢 CHƯA FIX

### Vấn đề gì?

```
handlers/AuthorHandler.go    ← PascalCase
handlers/category_handler.go ← snake_case
```

### Vì sao phải fix?

- Go convention: file name dùng **snake_case**
- Không gây lỗi code, nhưng không professional

### Fix gì?

Rename `AuthorHandler.go` → `author_handler.go`

**Scope:** Rename 1 file

---

## Fix #8: Hardcoded "Anonymous" User 🟢 CHƯA FIX

### Vấn đề gì?

```go
ok, err := c.Repo.Delete(req, "Anonymous")  // ← hardcode ở 10+ chỗ
```

### Vì sao phải fix?

- Khi thêm auth (JWT), cần lấy user từ context
- Hardcode ở 10+ chỗ → sửa tốn công, dễ bỏ sót

### Fix gì?

Tạo helper `GetCurrentUser(c echo.Context) string` → mặc định trả `"Anonymous"`. Sau này có auth chỉ sửa 1 chỗ.

**Scope:** 1 file helper mới, sửa handlers + services

---

## Fix #9: Clean Up Empty Packages 🟢 CHƯA FIX

### Vấn đề gì?

```
dto/auth/login.go      → chỉ `package auth` (rỗng)
models/user/role.go    → chỉ `package user` (rỗng)
models/user/user.go    → chỉ `package user` (rỗng)
```

### Vì sao phải fix?

- Gây confuse: không biết "chưa implement" hay "đã xóa logic"

### Fix gì?

Thêm comment TODO hoặc xóa file rỗng.

**Scope:** 3 file
