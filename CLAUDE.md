# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go web service for an action camera platform using a classic layered architecture. The service provides user authentication, registration with email verification, and uses Redis for caching and MySQL for persistent storage.

## Building and Running

### Prerequisites
- Go 1.25.3 or higher
- MySQL database running on localhost:3306
- Redis server running on localhost:6379

### Running the Application
```bash
# Install dependencies
go mod download

# Run the application
go run cmd/main.go
```

The server will start on the port specified in `config/config.ini` (default: 3000).

### Building
```bash
go build -o action-camera cmd/main.go
./action-camera
```

## Architecture

This project follows a strict **4-layer architecture**:

```
API Layer (api/v1/) → Service Layer (service/) → DAO Layer (dao/) → Model Layer (model/)
                                    ↓
                              Cache Layer (cache/)
```

### Layer Responsibilities

1. **API Layer** (`api/v1/`): 
   - Thin HTTP handlers using Gin
   - Only responsible for request binding and response formatting
   - Delegates all logic to Service layer
   - Pattern: Bind request → Call service method → Return JSON response

2. **Service Layer** (`service/`):
   - Contains all business logic and validation
   - Orchestrates between DAO and Cache layers
   - Returns standardized responses using `serizlizer.Response`
   - Handles error codes from `pkg/e` package
   - Must validate all input before calling DAO

3. **DAO Layer** (`dao/`):
   - Data Access Objects for database operations
   - Each DAO struct embeds `*gorm.DB`
   - Factory pattern: `NewXxxDao(ctx)` creates DAO with context
   - Returns domain models, not DTOs
   - Pattern example from `user.go`:
     ```go
     type UserDao struct {
         *gorm.DB
     }
     func NewUserDao(ctx context.Context) *UserDao {
         return &UserDao{NewDBClient(ctx)}
     }
     ```

4. **Model Layer** (`model/`):
   - GORM models representing database tables
   - Contains model-specific methods (e.g., `SetPassword`, `CheckPassword`)
   - Uses `gorm` struct tags for column mapping

5. **Cache Layer** (`cache/`):
   - Redis operations wrapper
   - RDB struct embeds `*redis.Client`
   - Factory pattern: `NewRDB(ctx)` creates cache client
   - Currently used for email verification codes with 5-minute TTL
   - Key pattern: `verify:email:{email}`

## Configuration

### config/config.ini Structure
- `[service]`: Application mode and HTTP port
- `[mysql]`: Database connection details
- `[redis]`: Redis connection and pool settings
- `[email]`: SMTP settings for email verification

**Important**: Never commit real credentials. The current `config.ini` should be template-only or gitignored.

### Database Initialization
The application automatically:
1. Connects to MySQL using settings from config.ini
2. Sets up master-slave configuration (currently both point to same instance)
3. Runs migrations via `dao.Migration()` on startup
4. Configures connection pool (max 20 connections, 30s lifetime)

## Key Patterns and Conventions

### Context Propagation
All layers accept `context.Context` as first parameter:
- API layer gets context from `c.Request.Context()`
- Pass context through Service → DAO → Database
- Cache operations also use context

### DAO Creation Pattern
Always create DAO instances per request:
```go
userDao := dao.NewUserDao(ctx)
```
Never reuse DAO instances across requests.

### Service Response Pattern
Services return `serizlizer.Response` with:
- `Status`: Error code from `pkg/e`
- `Msg`: Human-readable message
- `Error`: Detailed error (for debugging)
- `Data`: Response payload (optional)

### Password Handling
- Passwords are encrypted in the Model layer using `user.SetPassword()`
- Verification uses `user.CheckPassword()`
- Never store plain text passwords in any layer

### Email Verification Flow
1. User requests verification code via `/api/v1/user/vaild-email`
2. Service generates 6-digit code, stores in Redis with 5-minute TTL
3. User submits code during registration
4. Service validates code from Redis before creating user

### Database Master-Slave Setup
The codebase supports read-write splitting via `gorm.io/plugin/dbresolver`:
- Sources: Write operations
- Replicas: Read operations with random policy
- Currently configured with same connection for both (single instance)
- To enable real master-slave: update `connRead` and `connWrite` in config initialization

## Router Structure

Routes are defined in `routes/router.go`:
- Base path: `/api/v1`
- CORS middleware applied globally via `middleware.Cors()`
- Current endpoints:
  - `POST /api/v1/user/login`: User authentication
  - `POST /api/v1/user/vaild-email`: Send verification email
  - `POST /api/v1/user/register`: User registration with email verification

### Adding New Endpoints
1. Define handler in appropriate `api/v1/*.go` file
2. Create service method in corresponding `service/*.go` file
3. Add DAO methods if new data access is needed
4. Register route in `routes/router.go`

## Database

### Connection Details
- Default database name: `action_camera`
- Charset: `utf8mb4` with `parseTime=true`
- ORM: GORM with MySQL driver
- Naming: Singular table names (SingularTable: true)

### Migrations
Migrations run automatically on startup via `dao.Migration()` in `dao/migration.go`. Check this file to see which models are auto-migrated.

## Dependencies

Key dependencies:
- `github.com/gin-gonic/gin`: Web framework
- `gorm.io/gorm` + `gorm.io/driver/mysql`: ORM and MySQL driver
- `github.com/redis/go-redis/v9`: Redis client
- `github.com/dgrijalva/jwt-go`: JWT token generation
- `golang.org/x/crypto`: Password hashing
- `gopkg.in/gomail.v2`: Email sending
- `gopkg.in/ini.v1`: INI file parsing
- `gorm.io/plugin/dbresolver`: Master-slave DB configuration

## Common Development Tasks

When adding new features, follow this checklist:
1. Define the model in `model/` if new data structure needed
2. Add migration in `dao/migration.go`
3. Create DAO with methods in `dao/`
4. Implement business logic in `service/`
5. Create API handler in `api/v1/`
6. Register route in `routes/router.go`
7. Update error codes in `pkg/e/` if new error types needed

## Code Style Notes

- Chinese comments are used throughout the codebase
- Error messages mix Chinese and English
- Variable naming uses camelCase
- Package imports are organized: standard library → external → internal
