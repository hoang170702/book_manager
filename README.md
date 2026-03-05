# Book Manager API

A RESTful API for managing books, authors, and categories built with Go.

## Tech Stack

- **Framework**: [Echo v4](https://echo.labstack.com/) - High performance HTTP framework
- **ORM**: [GORM](https://gorm.io/) - ORM library for Go
- **Database**: PostgreSQL
- **Logging**: Custom structured logger

## Project Structure

```
book-manager/
├── cmd/
│   └── api/
│       └── main.go           # Application entrypoint
├── internal/
│   ├── config/               # App configuration
│   ├── dto/                  # Data Transfer Objects
│   │   ├── common/           # Shared Request/Response structures
│   │   ├── category/         # Category DTOs
│   │   └── author/           # Author DTOs
│   ├── handlers/             # HTTP handlers (controllers)
│   ├── middleware/           # HTTP middleware stack
│   ├── mapper/               # Entity <-> DTO mappers
│   ├── models/               # Database models
│   ├── repositories/         # Database access layer
│   ├── routes/               # Route registration
│   ├── services/             # Business logic layer
│   │   └── impl/             # Service implementations
│   └── utils/                # Utilities (logger, enums, etc.)
├── pkg/
│   └── database/             # Database connection
├── .env                      # Environment variables
└── go.mod
```

## Enterprise Features

### 1. Middleware Stack

| Middleware | File | Purpose |
|------------|------|---------|
| **Recovery** | `internal/middleware/recovery.go` | Catches panics, logs stack trace, returns proper error response. Prevents server crash from unhandled exceptions. |
| **Request ID** | `internal/middleware/request_id.go` | Generates UUID or extracts `X-Request-ID` from header for request tracing across logs. |
| **Logger** | `internal/middleware/logging.go` | Logs all incoming requests and outgoing responses with timing and request ID. |

**Lý do thêm:**
- **Recovery**: Trong production, panic không được crash toàn bộ server. Middleware này bắt panic và trả về response lỗi thay vì để server chết.
- **Request ID**: Khi debug trong môi trường nhiều requests, cần có ID để trace log của một request cụ thể từ đầu đến cuối.
- **Logger**: Không thể debug hay monitor API mà không có log. Middleware log tự động mọi request/response.

### 2. Health Check Endpoint

```
GET /health
```

Response:
```json
{
  "status": "healthy",
  "database": "connected"
}
```

**Lý do thêm:**
- Load balancer và Kubernetes cần endpoint để kiểm tra app còn sống không
- Monitoring tools dùng để alert khi service down
- Kiểm tra database connection trước khi nhận traffic

### 3. Graceful Shutdown

**Lý do thêm:**
- Khi deploy mới, cần đợi requests đang xử lý hoàn thành trước khi tắt server
- Đóng database connection sạch sẽ để tránh connection leak
- Xử lý SIGINT (Ctrl+C) và SIGTERM (kill) properly

### 4. Response ID Tracking

Mỗi response trả về `response_id` match với `request_id` từ client:

```json
{
  "response_id": "abc-123",
  "response_code": "00",
  "response_msg": "Success",
  "data": {...}
}
```

**Lý do thêm:**
- Client có thể trace response về request tương ứng
- Dễ debug khi client báo lỗi, chỉ cần gửi request_id
- Support idempotency checking trong tương lai

## Getting Started

### Prerequisites

- Go 1.24+
- PostgreSQL

### Installation

```bash
# Clone repository
git clone <repository-url>
cd book-manager

# Copy environment file
cp .example.env .env

# Edit .env with your database credentials
# DB_HOST=localhost
# DB_PORT=5432
# DB_USER=postgres
# DB_PASSWORD=your_password
# DB_NAME=book_manager

# Install dependencies
go mod download

# Run the server
go run cmd/api/main.go
```

### API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| POST | `/book-store/api/categories/add` | Create category |
| POST | `/book-store/api/categories/get-one` | Get category by ID |
| POST | `/book-store/api/categories/get-all` | Get all categories |
| POST | `/book-store/api/categories/update` | Update category |
| POST | `/book-store/api/categories/delete` | Delete category |
| POST | `/book-store/api/authors/add` | Create author |
| POST | `/book-store/api/authors/get-one` | Get author by ID |
| POST | `/book-store/api/authors/get-all` | Get all authors |
| POST | `/book-store/api/authors/update` | Update author |
| POST | `/book-store/api/authors/delete` | Delete author |

### Example Request

```bash
curl -X POST http://localhost:8080/book-store/api/categories/get-all \
  -H "Content-Type: application/json" \
  -H "X-Request-ID: my-trace-id" \
  -d '{"request_id": "my-trace-id", "data": {}}'
```

## Running Tests

```bash
go test -v ./...
```

## Future Improvements

- [ ] Structured JSON logging (zap/zerolog)
- [ ] Repository interfaces for better testability
- [ ] Context propagation through all layers
- [ ] Docker & docker-compose support
- [ ] Swagger/OpenAPI documentation
- [ ] Rate limiting middleware
- [ ] Authentication & Authorization (JWT)
