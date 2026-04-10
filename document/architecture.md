# Book Manager — Architecture Overview

## Tech Stack
- **Go 1.24** + **Echo v4** (HTTP framework)
- **GORM** + **PostgreSQL** (ORM + Database)
- **golang-jwt/jwt/v5** (JWT authentication)
- **bcrypt** (`golang.org/x/crypto`) (password hashing)
- **go-playground/validator/v10** (DTO validation)
- **UUID** (request tracing)
- **godotenv** (.env configuration)

## Project Structure
```
book-manager/
├── cmd/api/main.go              ← Entry point, DI wiring, graceful shutdown
├── document/                    ← Changelog, architecture, fixes plan
├── internal/
│   ├── config/                   ← App config (Port, JWTSecret) + DB migration
│   ├── constants/                ← HTTP status constants
│   ├── dto/                      ← Data Transfer Objects (độc lập với models)
│   │   ├── auth/                 ← Register/Login/Refresh DTOs
│   │   ├── author/               ← Author request + response DTOs
│   │   ├── category/             ← Category request + response DTOs
│   │   └── common/               ← Generic Request[T], Response[T], Paginate
│   ├── handlers/                 ← HTTP handlers (Controller layer)
│   ├── mapper/                   ← DTO → Entity mapping
│   ├── middleware/               ← Recovery, RequestID, Logger, Validator, JWTAuth, RateLimit
│   ├── models/                   ← GORM entities (internal, không expose ra API)
│   │   ├── base/                 ← AbstractStatus, AbsTimestamp (embedded)
│   │   ├── book/                 ← Book entity
│   │   ├── relations/            ← N-N join tables
│   │   └── user/                 ← User entity (username, hashed password)
│   ├── repositories/             ← Data access layer (interface + impl)
│   ├── routes/                   ← Route registration
│   ├── services/                 ← Business logic interfaces
│   │   └── impl/                 ← Service implementations
│   └── utils/
│       ├── auth/                 ← JWT generate/validate + GetCurrentUser helper
│       ├── enums/error_codes/    ← Error code system (ErrorCode + AppError)
│       └── logger/               ← Custom logger + GORM logger with Request ID
└── pkg/database/                 ← DB connection (returns *gorm.DB, no globals)
```

## Layered Architecture
```
Handler (HTTP) → Service (Business Logic) → Repository (Data Access) → GORM/DB
```
- **Handler** depends on `services.IXxxService` (interface)
- **Service** depends on `repositories.IXxxRepository` (interface)
- **Repository** depends on `*gorm.DB` (injected)

## Data Flow
```
Request JSON → Handler (Bind + Validate) → Service → Repository → DB
DB → Repository (models.Entity) → Service (map → ResponseDTO) → Handler → Response JSON
```
- **Request**: Client gửi JSON → Bind vào `Request[T]` → Validate `reqDto.Data`
- **Response**: Repository trả `models.Entity` → Service map sang `XxxResponse` DTO → client chỉ thấy `id`, `name`

## Authentication Flow
```
POST /auth/login → AuthService → bcrypt verify → JWT GenerateToken
                                                  → { access_token (1h), refresh_token (7d) }

GET /categories  → JWTAuth Middleware → ValidateToken (type=access)
                 → Check blacklist (revoked_tokens) → Set claims → Handler

POST /auth/refresh → Validate refresh_token (type=refresh)
                   → Generate new { access_token, refresh_token }

POST /auth/logout  → [Protected] Extract access token from header
                   → Blacklist access_token + refresh_token vào revoked_tokens
```
- **Public routes**: `/auth/register`, `/auth/login`, `/auth/refresh`
- **Protected routes**: Category, Author, `/auth/logout` (yêu cầu `Authorization: Bearer <access_token>`)

## Dependency Injection Flow
```go
main.go
  └─ db := database.Connect()
  └─ config.RunMigrations(db)       // includes User table
  └─ middleware.RegisterValidator(e)
  └─ setupMiddleware(e)             // CORS → Recovery → RequestID → Logger
  └─ routes.RegisterRoutes(e, db)
       ├─ Public group (no JWT):
       │    └─ AuthRepo → AuthService → AuthHandler → /auth/*
       └─ Protected group (JWTAuth middleware):
            ├─ CategoryRepo → CategoryService → CategoryHandler → /categories/*
            └─ AuthorRepo → AuthorService → AuthHandler → /authors/*
```

## Key Design Patterns
- **Generic DTO**: `Request[T]` / `Response[T]` — enterprise pattern giống `ApiResponse<T>` trong Java
- **Response DTO Mapping**: Service map entity → response DTO, không expose DB fields ra API
- **Interface-based DI**: Service → Repository đều qua interface, dễ mock test
- **JWT Authentication**: Access token (1h) + Refresh token (7d), bcrypt password hashing
- **Route Protection**: Public/Protected group split, JWT middleware chỉ chấp nhận access token
- **Rate Limiting**: 10 req/min per IP cho `/auth/login` và `/auth/register`, chống brute-force
- **DTO Validation**: `go-playground/validator` tích hợp Echo, validate `reqDto.Data` sau Bind
- **Safe Error Handling**: Luôn dùng `errors.As()` thay vì type assertion, có fallback cho unexpected errors
- **Custom GORM Logger**: Chỉ log slow query >200ms, kèm Request ID cho traceability
- **HTTP Logging**: Log request/response body, ưu tiên body `request_id` thay UUID
- **Soft Delete**: Update `status = "deleted"` thay vì xóa thật
- **Base model embedded**: `AbstractStatus` + `AbsTimestamp` cho tất cả entities
- **Graceful Shutdown**: Signal handling + context timeout
