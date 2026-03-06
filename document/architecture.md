# Book Manager — Architecture Overview

## Tech Stack
- **Go 1.24** + **Echo v4** (HTTP framework)
- **GORM** + **PostgreSQL** (ORM + Database)
- **UUID** (request tracing)
- **godotenv** (.env configuration)

## Project Structure
```
book-manager/
├── cmd/api/main.go              ← Entry point, DI wiring, graceful shutdown
├── internal/
│   ├── config/                   ← App config + DB migration
│   ├── constants/                ← HTTP status constants
│   ├── dto/                      ← Data Transfer Objects
│   │   ├── auth/                 ← (reserved for auth DTOs)
│   │   ├── author/               ← Author request DTOs
│   │   ├── category/             ← Category request DTOs
│   │   └── common/               ← Generic Request[T] / Response[T]
│   ├── handlers/                 ← HTTP handlers (Controller layer)
│   ├── mapper/                   ← DTO → Entity mapping
│   ├── middleware/               ← Recovery, RequestID, Logging (req/res body)
│   ├── models/                   ← GORM entities
│   │   ├── base/                 ← AbstractStatus, AbsTimestamp (embedded)
│   │   ├── book/                 ← Book entity
│   │   ├── relations/            ← N-N join tables
│   │   └── user/                 ← (reserved for user entities)
│   ├── repositories/             ← Data access layer (interface + impl)
│   ├── routes/                   ← Route registration
│   ├── services/                 ← Business logic interfaces
│   │   └── impl/                 ← Service implementations
│   └── utils/
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

## Dependency Injection Flow
```go
main.go
  └─ db := database.Connect()
  └─ config.RunMigrations(db)
  └─ routes.RegisterRoutes(e, db)
       └─ repo := &repositories.XxxRepository{DB: db}
       └─ service := impl.NewXxxService(repo)
       └─ handler := handlers.NewXxxHandler(service)
       └─ route.register(apiGroup)
```

## Key Design Patterns
- **Generic DTO**: `Request[T]` / `Response[T]` — enterprise pattern giống `ApiResponse<T>` trong Java
- **Interface-based DI**: Service → Repository đều qua interface
- **Soft Delete**: Update `status = "deleted"` thay vì xóa thật
- **Base model embedded**: `AbstractStatus` + `AbsTimestamp` cho tất cả entities
- **Custom GORM Logger**: Chỉ log slow query >200ms, kèm Request ID cho traceability
- **DTO Validation**: `go-playground/validator` tích hợp Echo, validate `reqDto.Data` sau Bind
- **Graceful Shutdown**: Signal handling + context timeout
