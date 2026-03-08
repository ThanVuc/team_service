# Team Service

Microservice quản lý team, group, sprint và work items sử dụng Go, gRPC, PostgreSQL và RabbitMQ.

---

## 📁 Cấu trúc thư mục

```
team_service/
├── dev.yaml                          # File cấu hình development
├── Dockerfile                        # Docker image definition
├── go.mod                            # Go module dependencies
├── main.go                           # Entry point của ứng dụng
├── README.md                         # Tài liệu hướng dẫn
│
├── cmd/                              # Command layer - khởi tạo ứng dụng
│   └── boostrap.go                   # Bootstrap logic, khởi động và shutdown service
│
├── global/                           # Global dependency injection
│   ├── di.go                         # Quản lý dependency injection toàn cục
│   └── lifecycle.go                  # Interface lifecycle cho các component
│
├── internal/                         # Business logic chính (private)
│   ├── adapter/                      # Adapter layer - Controllers và Handlers
│   │   ├── dependency.go            # Khởi tạo adapters
│   │   ├── gprc/                    # gRPC Controllers
│   │   │   ├── group.controller.go  # Controller xử lý Group gRPC requests
│   │   │   ├── sprint.controller.go # Controller xử lý Sprint gRPC requests
│   │   │   └── work.controller.go   # Controller xử lý Work gRPC requests
│   │   ├── job/                     # Job handlers (scheduled tasks)
│   │   │   └── example.handler.go   # Example job handler
│   │   └── messaging/               # Message queue handlers
│   │       └── example.handler.go   # RabbitMQ consumer handlers
│   │
│   ├── application/                 # Application layer - Business logic
│   │   ├── dependency.go           # Khởi tạo use cases
│   │   ├── common/                 # Shared application components
│   │   │   ├── interface/          # Interfaces/contracts
│   │   │   │   ├── repository/     # Repository interfaces
│   │   │   │   └── store/          # Store interfaces
│   │   │   ├── mapper/             # Data mappers (DTO ↔ Entity)
│   │   │   ├── model/              # Application models/DTOs
│   │   │   └── validation/         # Business validation logic
│   │   └── usecase/                # Use cases (business logic)
│   │       ├── facade.go           # Use case interfaces
│   │       ├── group.usecase.go    # Group business logic
│   │       ├── sprint.usecase.go   # Sprint business logic
│   │       └── work.usecase.go     # Work business logic
│   │
│   ├── domain/                      # Domain layer - Core business rules
│   │   ├── common/                 # Domain common components
│   │   │   └── apperror/           # Application error definitions
│   │   │       └── errordictionary/# Error code dictionaries
│   │   └── entity/                 # Domain entities
│   │       └── entity.example.go   # Example domain entity
│   │
│   ├── infrastructure/              # Infrastructure layer
│   │   ├── dependency.go           # Infrastructure dependency initialization
│   │   ├── logging/                # Logging infrastructure
│   │   │   └── logger.go           # Logger configuration
│   │   ├── messaging/              # Message queue infrastructure
│   │   │   └── eventbus.go         # RabbitMQ connector setup
│   │   ├── persistence/            # Database infrastructure
│   │   │   ├── connection.go       # Database connection pool
│   │   │   ├── db/                 # SQLC generated code
│   │   │   │   ├── sqlc.yaml       # SQLC configuration
│   │   │   │   ├── database/       # Generated SQLC code
│   │   │   │   └── sql/            # SQL files
│   │   │   │       ├── query/      # SQL queries cho SQLC
│   │   │   │       └── schema/     # Database migrations
│   │   │   ├── mapper/             # Data mappers (SQLC row → Domain Entity)
│   │   │   │   └── group.mapper.go # Mapper từ SQLC row to domain entity
│   │   │   ├── repository/         # Repository implementations
│   │   │   │   ├── group.repository.go
│   │   │   │   ├── sprint.repository.go
│   │   │   │   └── work.repository.go
│   │   │   └── store/              # Store pattern implementation
│   │   │       ├── repositorycontainer.go  # Repository container
│   │   │       └── store.go        # Transaction store
│   │   └── share/                  # Shared infrastructure utilities
│   │       ├── settings/           # Configuration management
│   │       └── utils/              # Utility functions
│   │
│   └── transport/                   # Transport layer - Server HTTP/gRPC
│       ├── dependency.go           # Transport layer initialization
│       └── grpc/                   # gRPC transport
│           └── grpcserver.go       # gRPC server implementation
│
├── proto/                           # Protocol Buffers definitions & generated code
│   ├── common/                     # Common proto messages
│   │   ├── common.pb.go
│   │   ├── enum.pb.go
│   │   ├── error.pb.go
│   │   ├── outbox.pb.go
│   │   ├── pagination.pb.go
│   │   └── syncdatabase*.pb.go
│   └── team_service/               # Team service specific protos
│       ├── common.team.pb.go
│       ├── group*.pb.go            # Group service proto
│       ├── sprint*.pb.go           # Sprint service proto
│       └── work*.pb.go             # Work service proto
│
└── scripts/                         # Utility scripts
    ├── build.sh                    # Script build và đẩy Docker image
    ├── gen-sqlc.sh                 # Generate SQLC code từ SQL files
    └── migrate.sh                  # Database migration management
```

---

## 📖 Giải thích các thư mục chính

### `cmd/`
Command layer - điểm khởi đầu của ứng dụng. Chứa logic bootstrap để khởi động service và graceful shutdown.

### `global/`
Quản lý dependency injection toàn cục và lifecycle của tất cả các component (infrastructure, application, adapter, transport).

### `internal/adapter/`
**Adapter Layer** - Layer này chứa các controllers và handlers:
- **grpc/**: Controllers xử lý gRPC requests, chuyển đổi proto messages thành calls tới use cases
- **messaging/**: Handlers xử lý messages từ RabbitMQ
- **job/**: Handlers cho scheduled jobs/cron tasks

### `internal/application/`
**Application Layer** - Business logic layer:
- **usecase/**: Chứa business logic chính, orchestrate giữa repositories và domain
- **common/interface/**: Định nghĩa interfaces cho repositories và stores
- **common/mapper/**: Convert giữa DTOs và domain entities
- **common/validation/**: Business validation rules

### `internal/domain/`
**Domain Layer** - Core business layer:
- **entity/**: Domain entities - đại diện cho core business objects
- **apperror/**: Application-specific error definitions và error dictionaries

### `internal/infrastructure/`
**Infrastructure Layer** - Technical implementations:
- **logging/**: Logger setup và configuration
- **messaging/**: RabbitMQ connection và event bus
- **persistence/**: Database connections, repositories, mappers, stores, và SQLC generated code
- **share/**: Configuration loading và utility functions

### `internal/transport/`
**Transport Layer** - Entry points cho external requests:
- **grpc/**: gRPC server implementation, register services

### `proto/`
Protocol Buffers definitions và generated Go code cho gRPC services.

### `scripts/`
Helper scripts cho development và deployment.

---

## 🛠️ Hướng dẫn sử dụng Scripts

### 1. `scripts/build.sh`
Build Docker image và push lên Docker Hub.

**Cách sử dụng:**
```bash
./scripts/build.sh
```

**Chức năng:**
- Build Docker image với tag `latest`
- Login vào Docker Hub
- Push image lên registry

**Lưu ý:** Cần sửa biến `DOCKER_USERNAME` và `DOCKER_REPO` trong script cho phù hợp với project của bạn.

---

### 2. `scripts/gen-sqlc.sh`
Generate Go code từ SQL queries và schema sử dụng SQLC.

**Cách sử dụng:**
```bash
./scripts/gen-sqlc.sh
```

**Chức năng:**
- Xóa thư mục generated code cũ
- Generate Go code từ SQL files trong `internal/infrastructure/persistence/db/sql/`
- Output vào `internal/infrastructure/persistence/db/database/`

**Khi nào cần chạy:**
- Sau khi thêm/sửa SQL queries trong `sql/query/`
- Sau khi thay đổi database schema

---

### 3. `scripts/migrate.sh`
Quản lý database migrations sử dụng goose.

**Cách sử dụng:**

**Tạo migration mới:**
```bash
./scripts/migrate.sh create tên_migration
```
**Ghi chú:** Lệnh này tạo file migration mới trong `internal/infrastructure/persistence/db/sql/schema/`. Tên file sẽ có format `YYYYMMDDHHmmss_tên_migration.sql`

**Xem trạng thái migrations:**
```bash
./scripts/migrate.sh status
```

**Chạy tất cả migrations (up):**
```bash
./scripts/migrate.sh up
```

**Rollback migration gần nhất:**
```bash
./scripts/migrate.sh down
```

**Reset tất cả migrations:**
```bash
./scripts/migrate.sh reset
```

**Rollback đến version cụ thể:**
```bash
./scripts/migrate.sh down-to 20260307130627
```

**Migrate lên version cụ thể:**
```bash
./scripts/migrate.sh up-to 20260307130627
```

**Chạy service:**
```bash
./scripts/migrate.sh run
```

**Lưu ý:** 
- Cần cấu hình database connection string trong biến `GOOSE_DBSTRING`
- Migration files nằm trong `internal/infrastructure/persistence/db/sql/schema/`

---

## 🔄 Workflow xử lý gRPC Request

### Kiến trúc Clean Architecture Flow

```
Client gRPC Request
        ↓
[Transport Layer] - grpcserver.go
        ↓
[Adapter Layer] - controller.go (gRPC Controller)
        ↓
[Application Layer] - usecase.go (Business Logic)
        ↓
[Infrastructure Layer] - repository.go (Database Access)
        ↓
PostgreSQL Database
```

### Checklist triển khai một gRPC endpoint mới

#### ✅ Bước 1: Định nghĩa Protocol Buffer
**File:** `proto/team_service/group.proto` (ví dụ) hoặc repo `schedule_proto`

- [ ] Định nghĩa message request
- [ ] Định nghĩa message response
- [ ] Định nghĩa service method trong service interface
- [ ] Generate Go code từ proto (nếu sử dụng `schedule_proto` repo, sinh code theo instruction trong repo đó)
- [ ] Hoặc generate thủ công: `protoc --go_out=. --go-grpc_out=. proto/team_service/*.proto`

---

#### ✅ Bước 2: Tạo/Cập nhật Database Schema & Queries

**File:** `internal/infrastructure/persistence/db/sql/schema/YYYYMMDD_*.sql`

- [ ] Chạy script để tạo file migration mới: `./scripts/migrate.sh create tên_migration`
- [ ] Viết migration script/SQL nếu cần thay đổi schema
- [ ] Chạy migration: `./scripts/migrate.sh up`

**File:** `internal/infrastructure/persistence/db/sql/query/*.sql`

- [ ] Viết SQL queries với SQLC annotation (-- name:, -- :one, -- :many, etc.)
- [ ] Generate SQLC code: `./scripts/gen-sqlc.sh`

---

#### ✅ Bước 3: Mapper (SQLC Row → Domain Entity)

**File:** `internal/infrastructure/persistence/mapper/[domain].mapper.go`

- [ ] Tạo mapper struct
- [ ] Implement method để convert từ SQLC row struct sang domain entity
- [ ] Xử lý type conversion và null checking
- [ ] Reuse trong repository methods

---

#### ✅ Bước 4: Repository Interface & Implementation

**File:** `internal/application/common/interface/repository/group.repository.go`

- [ ] Định nghĩa repository interface với methods cần thiết
- [ ] Đặt tên methods theo business logic (VD: `CreateGroup`, `GetGroupByID`)

**File:** `internal/infrastructure/persistence/repository/group.repository.go`

- [ ] Implement repository interface
- [ ] Inject SQLC Queries struct
- [ ] Gọi SQLC generated methods
- [ ] Handle errors và convert sang domain entities

---

#### ✅ Bước 5: Use Case (Business Logic)

**File:** `internal/application/usecase/facade.go`

- [ ] Định nghĩa use case interface

**File:** `internal/application/usecase/group.usecase.go`

- [ ] Implement use case interface
- [ ] Inject Store (để có transaction support)
- [ ] Viết business logic
- [ ] Sử dụng `store.ExecTx()` nếu cần transaction
- [ ] Handle business validation
- [ ] Return errors sử dụng error dictionary

**File:** `internal/application/dependency.go`

- [ ] Khởi tạo use case trong `NewDependency`
- [ ] Wire dependencies (store, repositories)

---

#### ✅ Bước 6: gRPC Controller

**File:** `internal/adapter/gprc/group.controller.go`

- [ ] Implement gRPC service interface từ generated proto code
- [ ] Inject use case vào controller
- [ ] Parse request proto message
- [ ] Call use case method
- [ ] Convert response sang proto message
- [ ] Handle errors và convert sang gRPC status codes
- [ ] **Important**: Bắt các panic errors (lỗi crash hệ thống) và log ra mà không làm hệ thống crash (khôi phục gracefully)

**File:** `internal/adapter/dependency.go`

- [ ] Khởi tạo controller trong `NewDependency`
- [ ] Wire use case dependency

---

#### ✅ Bước 7: Register Service với gRPC Server

**File:** `internal/transport/grpc/grpcserver.go`

- [ ] Import generated proto service package
- [ ] Register service trong `Start()` method: `team_service.RegisterXxxServiceServer()`
- [ ] Pass controller instance từ adapter

---

#### ✅ Bước 8: Error Dictionary (Optional nhưng recommended)

**File:** `internal/domain/common/apperror/errordictionary/group.errordict.go`

- [ ] Define error codes và messages cho domain này
- [ ] Sử dụng consistent error codes

---

### Ví dụ flow xử lý CreateGroup request:

1. **Client** gửi `CreateGroupRequest` tới gRPC server
2. **Transport Layer** (`grpcserver.go`) receive request, route tới `GroupController`
3. **Adapter** (`group.controller.go`) parse request, call `groupUseCase.CreateGroup()`
4. **Application** (`group.usecase.go`) thực thi business logic:
   - Validate input
   - Call `store.ExecTx()` để bắt đầu transaction
   - Call `groupRepository.CreateGroup()` để lưu vào database
   - Commit transaction
5. **Infrastructure** (`group.repository.go`) thực thi SQL query
6. **Database** lưu data và return
7. Response flow ngược lại: Repository → UseCase → Controller → Client

---

## 📨 Workflow xử lý Message Queue (RabbitMQ Consumer)

### Kiến trúc Event-Driven Flow

```
RabbitMQ Message → Handler → UseCase → Repository → Database
```

### Các bước triển khai Consumer

#### ✅ Bước 1: Tạo Message Handler
**File:** `internal/adapter/messaging/[tên_handler].handler.go`

- [ ] Tạo handler struct với use case dependency
- [ ] Parse message payload thành Go struct
- [ ] Call use case để xử lý business logic
- [ ] ACK/NACK message dựa trên kết quả
- [ ] **Important**: Bắt panic errors để tránh crash

#### ✅ Bước 2: Implement Business Logic
**File:** `internal/application/usecase/[domain].usecase.go`

- [ ] Thêm method xử lý logic trong use case
- [ ] Sử dụng transaction nếu cần (`store.ExecTx()`)

#### ✅ Bước 3: Register Consumer
**File:** `internal/adapter/dependency.go`

- [ ] Khởi tạo handler
- [ ] Register consumer với RabbitMQ connector
- [ ] Config queue name, exchange, routing key

#### ✅ Bước 4: Cấu hình
**File:** `dev.yaml`

- [ ] RabbitMQ connection (host, port, user, password)
- [ ] Consumer settings (queue name, prefetch count)

**Lưu ý quan trọng:**
- Messages có thể được gửi lại nhiều lần → implement idempotency
- ACK khi success, NACK khi fail (có thể requeue hoặc gửi DLQ)
- Sử dụng transaction để đảm bảo data consistency

---

## ⏰ Workflow xử lý Scheduled Jobs

### Kiến trúc Background Job Flow

```
Scheduler/Cron → Job Handler → UseCase → Repository → Database
```

### Các bước triển khai Job

#### ✅ Bước 1: Tạo Job Handler
**File:** `internal/adapter/job/[tên_job].handler.go`

- [ ] Tạo handler struct với use case dependency
- [ ] Implement job execution method
- [ ] Call use case để xử lý logic
- [ ] Log job execution status (start, success, error)
- [ ] **Important**: Bắt panic errors để tránh crash

#### ✅ Bước 2: Implement Business Logic
**File:** `internal/application/usecase/[domain].usecase.go`

- [ ] Thêm method xử lý job logic trong use case
- [ ] Sử dụng transaction nếu cần (`store.ExecTx()`)
- [ ] Handle batch processing nếu xử lý nhiều records

#### ✅ Bước 3: Register Job với Scheduler
**File:** `internal/adapter/dependency.go` hoặc scheduler setup

- [ ] Khởi tạo job handler
- [ ] Config cron expression/schedule interval
- [ ] Register job với scheduler (cron, periodic timer, etc.)

#### ✅ Bước 4: Cấu hình
**File:** `dev.yaml` hoặc code

- [ ] Job schedule/cron expression
- [ ] Timeout settings
- [ ] Concurrency settings (nếu cần)

**Lưu ý quan trọng:**
- Jobs nên idempotent - chạy nhiều lần không gây lỗi
- Implement proper locking nếu job chỉ được chạy 1 instance tại 1 thời điểm
- Log chi tiết để dễ debug và monitor

---

## 🚀 Khởi chạy Service

**Prerequisites:**
- Go 1.24.3+
- PostgreSQL
- RabbitMQ
- SQLC (cho code generation)
- Goose (cho migrations)

**Steps:**

1. **Cấu hình `dev.yaml`**: Update database và RabbitMQ connection strings

2. **Run migrations**:
   ```bash
   ./scripts/migrate.sh up
   ```

3. **Generate SQLC code** (nếu có thay đổi queries):
   ```bash
   ./scripts/gen-sqlc.sh
   ```

4. **Run service**:
   ```bash
   go run main.go
   ```

5. **Hoặc build và run**:
   ```bash
   go build -o team_service
   ./team_service
   ```

---

## 📝 Notes

- Service sử dụng Clean Architecture với separation of concerns rõ ràng
- Dependency injection được quản lý centrally trong `global/` package
- SQLC được sử dụng để generate type-safe Go code từ SQL
- Goose được sử dụng cho database migrations
- gRPC được sử dụng cho synchronous communication
- RabbitMQ được sử dụng cho asynchronous event-driven communication
- Transaction support thông qua Store pattern

---

## 🤝 Contributing

Khi contribute code mới, hãy tuân theo Clean Architecture principles và checklist ở trên để đảm bảo consistency.
