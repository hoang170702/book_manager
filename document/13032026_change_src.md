# 13032026 — Change Source Log

## Tổng quan

Refactoring dự án Book Manager (Go/Echo) theo chuẩn enterprise: sửa 9 lỗi kiến trúc + triển khai Authentication & JWT.

---

## ĐÃ HOÀN THÀNH

### Fix #1: Global DB Variable + `pkg/` import `internal/` ✅
- Xóa global `var DB *gorm.DB`, `Connect()` return `*gorm.DB`
- Chuyển migration logic vào `internal/config/migration.go`
- Truyền `db` qua DI chain từ `main.go`

### Fix #2: Repository Interface ✅
- Tạo `ICategoryRepository`, `IAuthorRepository` interfaces
- Service depend vào interface thay concrete struct

### Fix #3: Validation Middleware ✅
- Thêm `go-playground/validator/v10`
- Tạo `middleware/validator.go`, gọi `c.Validate(&reqDto.Data)` sau `c.Bind()`

### Fix #4: Error Handling Consistency ✅
- Đổi `err.(*AppError)` → `errors.As(err, &appErr)` trong `category_service_impl.go`

### Fix #5: DTO Improvements ✅
- Tạo `CategoryResponse`, `AuthorResponse` — chỉ expose `id`, `name`
- Tạo `dto/common/paginate.go` thay `models.Paginate`
- Service map entity → response DTO

### Fix #6: CORS Middleware ✅
- Thêm `echomw.CORSWithConfig()` vào `main.go` (AllowOrigins: *)

### Fix #7: File Naming Consistency ✅
- Rename `AuthorHandler.go` → `author_handler.go`

### Fix #8: Hardcoded "Anonymous" User ✅
- Tạo `utils/auth/user.go` → `GetCurrentUser(c echo.Context) string`
- Handler đọc user từ context, truyền xuống Service → Repository

### Fix #9: Clean Up Empty Packages ✅
- Thêm `// TODO` comment vào 3 file rỗng (`dto/auth`, `models/user`)

### Feature: Authentication & JWT ✅
**Files tạo mới:**

| File | Mô tả |
|------|-------|
| `models/user/user.go` | User entity (password `json:"-"`, embed AbstractStatus + AbsTimestamp) |
| `dto/auth/login.go` | RegisterRequest, LoginRequest, LoginResponse, RefreshRequest |
| `utils/auth/jwt.go` | GenerateAccessToken (1h), GenerateRefreshToken (7d), ValidateToken |
| `repositories/auth_repository_interface.go` | `IAuthRepository` interface |
| `repositories/auth_repository.go` | Create (check duplicate username), FindByUsername |
| `services/auth_service.go` | `IAuthService` interface |
| `services/impl/auth_service_impl.go` | Register (bcrypt hash), Login (verify → JWT), RefreshToken (validate refresh → new pair) |
| `middleware/jwt_auth.go` | Extract Bearer token, validate access-only, set claims to context |
| `handlers/auth_handler.go` | POST register / login / refresh |
| `routes/auth_routes.go` | `/auth/register`, `/auth/login`, `/auth/refresh` (public) |

**Files đã sửa:**

| File | Thay đổi |
|------|----------|
| `config/migration.go` | Thêm `&user.User{}` vào AutoMigrate |
| `utils/auth/user.go` | Đọc username từ JWT claims thay vì return "Anonymous" |
| `routes/register_all.go` | Chia public (auth) / protected (category, author) groups |
| `error_codes/error_code.go` | Thêm `Unauthorized`, `UserAlreadyExist`, `InvalidCredential`, `InvalidToken` |
| `go.mod` | Thêm `github.com/golang-jwt/jwt/v5` |
| `document/architecture.md` | Thêm JWT tech stack, auth flow, design patterns |

**Token Flow:**
```
Login  → { access_token (1h), refresh_token (7d) }
API    → Authorization: Bearer <access_token>
Refresh → POST /auth/refresh { refresh_token } → new token pair
```

**API Endpoints:**
```
POST /book-store/api/auth/register   (public)
POST /book-store/api/auth/login      (public)
POST /book-store/api/auth/refresh    (public)

POST /book-store/api/categories/*    (protected - JWT required)
POST /book-store/api/authors/*       (protected - JWT required)
```

---

## CẦN LÀM TIẾP

### 1. Test Authentication Flow ✅ ĐÃ TEST
- User đã test bằng curl: Register → Login → API call với token → Refresh token
- Tất cả đều hoạt động tốt

### 2. Thêm JWT_SECRET vào .env ✅ ĐÃ LÀM
- Đã thêm `JWT_SECRET=bm-s3cr3t-k3y-ch4ng3-m3-1n-pr0duct10n` vào `.env`

### 3. Config struct cập nhật ✅ ĐÃ LÀM
- Thêm `JWTSecret string` vào `config/config.go`, đọc từ `JWT_SECRET` env var

### 4. Xóa file `models/user/role.go` ✅ ĐÃ LÀM
- Đã xóa file rỗng, chưa cần thiết cho giai đoạn hiện tại

### 5. Cập nhật `cmd/verification/main.go` ✅ ĐÃ LÀM
- MockService đã thêm `user string` param cho Create/Update/Delete
- Xóa import `models` không dùng

### 6. Viết Unit Tests 📋 CHƯA LÀM
- Auth Service: Register, Login (đúng/sai password), RefreshToken (đúng/sai type)
- JWT Utility: GenerateToken, ValidateToken, expired token
- JWT Middleware: no token, invalid token, refresh token (should reject)

### 7. Rate Limiting cho Auth endpoints ✅ ĐÃ LÀM
- Tạo `middleware/rate_limit.go` — in-memory sliding window, 10 req/min per IP
- Áp dụng cho `/auth/register` và `/auth/login` trong `auth_routes.go`
- `/auth/refresh` không rate limit vì đã yêu cầu refresh token hợp lệ

### 8. Logout / Token Blacklist ✅ ĐÃ LÀM
- Tạo `models/user/revoked_token.go` — lưu token đã bị vô hiệu hóa trong DB
- Thêm `RevokeToken`/`IsTokenRevoked` vào auth repository
- Thêm `Logout`/`IsTokenRevoked` vào auth service
- `POST /auth/logout` là **protected route** — bắt buộc gửi access token trong header
- JWT middleware check blacklist trước khi cho phép truy cập
