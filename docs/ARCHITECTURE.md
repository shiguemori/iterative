# B3 Challenge - System Architecture

## Overview

The B3 Challenge API is built using modern software engineering principles, emphasizing clean architecture, high performance, and maintainability. This document provides a comprehensive overview of the system's architecture, design decisions, and implementation details.

## Architectural Principles

### 1. Clean Architecture
The system follows Uncle Bob's Clean Architecture principles:
- **Independence of Frameworks**: Business logic is not dependent on external frameworks
- **Testability**: Business rules can be tested without UI, database, or external elements
- **Independence of UI**: The UI can change without changing business rules
- **Independence of Database**: Business rules are not bound to the database
- **Independence of External Agencies**: Business rules don't know about the outside world

### 2. SOLID Principles
- **Single Responsibility**: Each module has one reason to change
- **Open/Closed**: Open for extension, closed for modification
- **Liskov Substitution**: Objects should be replaceable with instances of their subtypes
- **Interface Segregation**: Many client-specific interfaces are better than one general-purpose interface
- **Dependency Inversion**: Depend on abstractions, not concretions

### 3. Domain-Driven Design (DDD)
- Clear separation between domain logic and infrastructure
- Rich domain models with business logic encapsulation
- Repository pattern for data access abstraction
- Service layer for business operations

## System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        External Systems                         │
├─────────────────┬─────────────────┬─────────────────────────────┤
│   Load Balancer │   Monitoring    │        Client Apps          │
│   (Nginx/HAProxy│   (Prometheus)  │   (Web, Mobile, Desktop)    │
└─────────┬───────┴─────────┬───────┴─────────────┬───────────────┘
          │                 │                     │
          └─────────────────┼─────────────────────┘
                            │
┌───────────────────────────▼───────────────────────────────────────┐
│                      API Gateway Layer                            │
│  ┌─────────────────────────────────────────────────────────────┐  │
│  │                 Gin HTTP Server                             │  │
│  │  ┌─────────────┬─────────────┬─────────────┬─────────────┐  │  │
│  │  │ CORS        │ Rate        │ Logging     │ Error       │  │  │
│  │  │ Middleware  │ Limiting    │ Middleware  │ Handling    │  │  │
│  │  └─────────────┴─────────────┴─────────────┴─────────────┘  │  │
│  └─────────────────────────────────────────────────────────────┘  │
└─────────────────────────────┬─────────────────────────────────────┘
                              │
┌─────────────────────────────▼─────────────────────────────────────┐
│                    Presentation Layer                             │
│  ┌─────────────────────────────────────────────────────────────┐  │
│  │                    API Handlers                             │  │
│  │  ┌─────────────┬─────────────┬─────────────┬─────────────┐  │  │
│  │  │ Health      │ Ticker      │ Aggregation │ Trade       │  │  │
│  │  │ Handler     │ Handler     │ Handler     │ Handler     │  │  │
│  │  └─────────────┴─────────────┴─────────────┴─────────────┘  │  │
│  └─────────────────────────────────────────────────────────────┘  │
└─────────────────────────────┬─────────────────────────────────────┘
                              │
┌─────────────────────────────▼─────────────────────────────────────┐
│                    Application Layer                              │
│  ┌─────────────────────────────────────────────────────────────┐  │
│  │                  Business Services                          │  │
│  │  ┌─────────────┬─────────────┬─────────────┬─────────────┐  │  │
│  │  │ Trade       │ Data Loader │ Mock Data   │ Validation  │  │  │
│  │  │ Service     │ Service     │ Service     │ Service     │  │  │
│  │  └─────────────┴─────────────┴─────────────┴─────────────┘  │  │
│  └─────────────────────────────────────────────────────────────┘  │
└─────────────────────────────┬─────────────────────────────────────┘
                              │
┌─────────────────────────────▼─────────────────────────────────────┐
│                     Domain Layer                                  │
│  ┌─────────────────────────────────────────────────────────────┐  │
│  │                   Domain Models                             │  │
│  │  ┌─────────────┬─────────────┬─────────────┬─────────────┐  │  │
│  │  │ Trade       │ Trade       │ Business    │ Domain      │  │  │
│  │  │ Entity      │ Aggregation │ Rules       │ Interfaces  │  │  │
│  │  └─────────────┴─────────────┴─────────────┴─────────────┘  │  │
│  └─────────────────────────────────────────────────────────────┘  │
└─────────────────────────────┬─────────────────────────────────────┘
                              │
┌─────────────────────────────▼─────────────────────────────────────┐
│                  Infrastructure Layer                             │
│  ┌─────────────────────────────────────────────────────────────┐  │
│  │                 Data Access Layer                           │  │
│  │  ┌─────────────┬─────────────┬─────────────┬─────────────┐  │  │
│  │  │ Trade       │ Database    │ HTTP        │ Config      │  │  │
│  │  │ Repository  │ Connection  │ Client      │ Manager     │  │  │
│  │  └─────────────┴─────────────┴─────────────┴─────────────┘  │  │
│  └─────────────────────────────────────────────────────────────┘  │
└─────────────────────────────┬─────────────────────────────────────┘
                              │
┌─────────────────────────────▼─────────────────────────────────────┐
│                     Data Storage Layer                            │
│  ┌─────────────────────────────────────────────────────────────┐  │
│  │                   PostgreSQL Database                       │  │
│  │  ┌─────────────┬─────────────┬─────────────┬─────────────┐  │  │
│  │  │ Trades      │ Indexes     │ Connection  │ Backup &    │  │  │
│  │  │ Table       │ &Constraints│ Pool        │ Recovery    │  │  │
│  │  └─────────────┴─────────────┴─────────────┴─────────────┘  │  │
│  └─────────────────────────────────────────────────────────────┘  │
└───────────────────────────────────────────────────────────────────┘
```

## Layer Responsibilities

### 1. API Gateway Layer
**Purpose**: Entry point for all HTTP requests, handling cross-cutting concerns.

**Components**:
- **Gin HTTP Server**: High-performance HTTP router and middleware engine
- **CORS Middleware**: Cross-Origin Resource Sharing configuration
- **Rate Limiting**: Request throttling and abuse prevention
- **Logging Middleware**: Request/response logging for monitoring
- **Error Handling**: Centralized error processing and response formatting

**Key Features**:
- Request routing and method handling
- Middleware chain execution
- Response formatting standardization
- Security headers management

### 2. Presentation Layer
**Purpose**: HTTP request handling and response formatting.

**Components**:
- **Health Handler**: System health and status endpoints
- **Ticker Handler**: Stock ticker information endpoints
- **Aggregation Handler**: Data aggregation and statistics endpoints
- **Trade Handler**: Individual trade record endpoints

**Responsibilities**:
- HTTP request parsing and validation
- Input parameter extraction and validation
- Response serialization to JSON
- HTTP status code management
- Error response formatting

### 3. Application Layer
**Purpose**: Business logic orchestration and use case implementation.

**Components**:
- **Trade Service**: Core business logic for trade operations
- **Data Loader Service**: External data collection and processing
- **Mock Data Service**: Test data generation for development
- **Validation Service**: Business rule validation

**Key Features**:
- Business rule enforcement
- Transaction coordination
- External service integration
- Data transformation and aggregation
- Caching strategy implementation

### 4. Domain Layer
**Purpose**: Core business entities and domain logic.

**Components**:
- **Trade Entity**: Core trade data structure and business rules
- **Trade Aggregation**: Statistical data models
- **Business Rules**: Domain-specific validation and logic
- **Domain Interfaces**: Contracts for external dependencies

**Characteristics**:
- Framework-independent
- Database-agnostic
- Pure business logic
- Rich domain models
- Encapsulated business rules

### 5. Infrastructure Layer
**Purpose**: External system integration and technical implementation.

**Components**:
- **Trade Repository**: Data persistence abstraction
- **Database Connection**: PostgreSQL connection management
- **HTTP Client**: External API communication
- **Configuration Manager**: Environment and settings management

**Responsibilities**:
- Database query execution
- External API communication
- File system operations
- Configuration loading
- Logging implementation

### 6. Data Storage Layer
**Purpose**: Persistent data storage and management.

**Components**:
- **PostgreSQL Database**: Primary data store
- **Trades Table**: Core trading data storage
- **Indexes**: Performance optimization structures
- **Connection Pool**: Efficient connection management

## Design Patterns

### 1. Repository Pattern
**Implementation**: `TradeRepository` provides data access abstraction.

```go
type TradeRepositoryInterface interface {
    CreateTrades(ctx context.Context, trades []models.Trade) error
    GetTradeAggregation(ctx context.Context, ticker string, fromDate *time.Time) (*models.TradeAggregation, error)
    GetTradesByDateRange(ctx context.Context, fromDate, toDate time.Time) ([]models.Trade, error)
    DeleteOldTrades(ctx context.Context, beforeDate time.Time) error
    GetDistinctTickers(ctx context.Context) ([]string, error)
}
```

**Benefits**:
- Database implementation abstraction
- Easy testing with mock repositories
- Centralized data access logic
- Consistent error handling

### 2. Service Layer Pattern
**Implementation**: Business logic encapsulation in service classes.

```go
type TradeServiceInterface interface {
    GetTradeAggregation(ctx context.Context, ticker string) (*models.TradeAggregation, error)
    GetAllTickersAggregation(ctx context.Context) ([]models.TradeAggregation, error)
    ValidateTicker(ctx context.Context, ticker string) (bool, error)
}
```

**Benefits**:
- Business logic centralization
- Transaction boundary management
- Reusable business operations
- Clear separation of concerns

### 3. Dependency Injection
**Implementation**: Constructor-based dependency injection.

```go
func NewTradeService(tradeRepo repository.TradeRepositoryInterface) *TradeService {
    return &TradeService{
        tradeRepo: tradeRepo,
    }
}
```

**Benefits**:
- Loose coupling between components
- Easy unit testing with mocks
- Flexible configuration
- Improved maintainability

### 4. Factory Pattern
**Implementation**: Service and repository creation.

```go
func SetupRouter(tradeService service.TradeServiceInterface) *gin.Engine {
    router := gin.New()
    handlers := NewHandlers(tradeService)
    // Configure routes...
    return router
}
```

**Benefits**:
- Centralized object creation
- Configuration encapsulation
- Consistent initialization
- Easy testing setup

## Data Flow

### 1. Request Processing Flow

```
Client Request
      ↓
API Gateway (Gin Router)
      ↓
Middleware Chain (CORS, Logging, Rate Limiting)
      ↓
Route Handler (Presentation Layer)
      ↓
Input Validation & Parameter Extraction
      ↓
Service Layer (Business Logic)
      ↓
Repository Layer (Data Access)
      ↓
Database Query Execution
      ↓
Result Processing & Aggregation
      ↓
Response Serialization
      ↓
HTTP Response to Client
```

### 2. Data Loading Flow

```
Data Loader Command
      ↓
Configuration Loading
      ↓
External API Client (BRAPI/B3)
      ↓
Data Fetching & Rate Limiting
      ↓
Data Validation & Transformation
      ↓
Batch Processing
      ↓
Repository Layer
      ↓
Database Insertion
      ↓
Success/Error Logging
```

### 3. Error Handling Flow

```
Error Occurrence
      ↓
Error Capture (defer/recover)
      ↓
Error Classification
      ↓
Logging & Monitoring
      ↓
Error Response Formatting
      ↓
HTTP Status Code Assignment
      ↓
JSON Error Response
      ↓
Client Error Handling
```

## Database Design

### Schema Architecture

```sql
-- Primary trades table
CREATE TABLE trades (
    id BIGSERIAL PRIMARY KEY,
    closing_time VARCHAR(8) NOT NULL,
    trade_date TIMESTAMP NOT NULL,
    instrument_code VARCHAR(10) NOT NULL,
    trade_price DECIMAL(10,2) NOT NULL,
    traded_quantity INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Performance optimization indexes
CREATE INDEX idx_trades_trade_date ON trades(trade_date);
CREATE INDEX idx_trades_instrument_code ON trades(instrument_code);
CREATE INDEX idx_trades_composite ON trades(instrument_code, trade_date);
```

### Indexing Strategy

1. **Primary Index**: `id` (automatic with PRIMARY KEY)
2. **Date Index**: `trade_date` for time-based queries
3. **Ticker Index**: `instrument_code` for ticker-specific queries
4. **Composite Index**: `(instrument_code, trade_date)` for combined queries

### Query Optimization

**Aggregation Queries**:
```sql
-- Optimized for max price lookup
SELECT MAX(trade_price) 
FROM trades 
WHERE instrument_code = ? AND trade_date >= ?;

-- Optimized for volume aggregation
SELECT MAX(daily_volume) 
FROM (
    SELECT SUM(traded_quantity) as daily_volume 
    FROM trades 
    WHERE instrument_code = ? AND trade_date >= ? 
    GROUP BY trade_date
) as daily_volumes;
```

**Performance Characteristics**:
- Index-only scans for aggregation queries
- Efficient range scans for date filtering
- Optimized GROUP BY operations
- Minimal table scans

## Performance Architecture

### 1. Database Performance

**Connection Pooling**:
```go
config := &gorm.Config{
    ConnPool: &sql.DB{
        MaxOpenConns:    25,
        MaxIdleConns:    5,
        ConnMaxLifetime: 5 * time.Minute,
    },
}
```

**Query Optimization**:
- Strategic indexing for common query patterns
- Efficient aggregation queries
- Proper use of LIMIT and OFFSET
- Connection reuse and pooling

### 2. Application Performance

**Memory Management**:
- Efficient struct packing
- Minimal memory allocations
- Proper garbage collection tuning
- Resource cleanup with defer statements

**Concurrency**:
- Goroutine-based request handling
- Context-based cancellation
- Channel communication patterns
- Mutex-free data structures where possible

### 3. Caching Strategy

**Current Implementation**:
- In-memory result caching for frequently accessed data
- HTTP response caching headers
- Database query result caching

**Future Enhancements**:
- Redis distributed caching
- Application-level cache warming
- Cache invalidation strategies
- CDN integration for static content

## Security Architecture

### 1. Current Security Measures

**Input Validation**:
- Parameter type validation
- Date format validation
- Range validation for queries
- SQL injection prevention through ORM

**HTTP Security**:
- CORS configuration
- Rate limiting implementation
- Request size limitations
- Secure headers configuration

### 2. Production Security Recommendations

**Authentication & Authorization**:
```go
// JWT middleware implementation
func JWTMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := extractToken(c.GetHeader("Authorization"))
        if !validateToken(token) {
            c.JSON(401, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

**Additional Security Measures**:
- HTTPS/TLS encryption
- API key authentication
- Request signing
- IP whitelisting
- Audit logging

## Monitoring and Observability

### 1. Logging Architecture

**Structured Logging**:
```go
log.WithFields(log.Fields{
    "method":      c.Request.Method,
    "path":        c.Request.URL.Path,
    "status":      c.Writer.Status(),
    "duration":    time.Since(start),
    "client_ip":   c.ClientIP(),
}).Info("Request processed")
```

**Log Levels**:
- DEBUG: Detailed debugging information
- INFO: General operational messages
- WARN: Warning conditions
- ERROR: Error conditions requiring attention

### 2. Metrics Collection

**Key Metrics**:
- Request rate (requests per second)
- Response time percentiles (50th, 95th, 99th)
- Error rate by endpoint
- Database query performance
- Memory and CPU usage

**Monitoring Integration**:
- Prometheus metrics exposition
- Grafana dashboard integration
- Alert manager configuration
- Health check endpoints

### 3. Distributed Tracing

**Implementation Strategy**:
- Request ID generation and propagation
- Span creation for major operations
- Context propagation across layers
- Performance bottleneck identification

## Scalability Architecture

### 1. Horizontal Scaling

**Load Balancing**:
```nginx
upstream b3_challenge_api {
    server api1:8080;
    server api2:8080;
    server api3:8080;
}

server {
    listen 80;
    location / {
        proxy_pass http://b3_challenge_api;
    }
}
```

**Stateless Design**:
- No server-side session storage
- Database-backed state management
- Shared cache for distributed systems
- Load balancer session affinity not required

### 2. Vertical Scaling

**Resource Optimization**:
- CPU-intensive operations optimization
- Memory usage profiling and optimization
- Database connection pool tuning
- Garbage collection optimization

**Performance Tuning**:
```go
// Runtime optimization
runtime.GOMAXPROCS(runtime.NumCPU())
debug.SetGCPercent(100)
```

### 3. Database Scaling

**Read Replicas**:
- Master-slave replication setup
- Read query distribution
- Write operation centralization
- Eventual consistency handling

**Sharding Strategy**:
- Horizontal partitioning by date
- Ticker-based sharding
- Consistent hashing implementation
- Cross-shard query optimization

## Deployment Architecture

### 1. Container Strategy

**Docker Configuration**:
```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o api cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/api .
CMD ["./api"]
```

**Benefits**:
- Consistent deployment environments
- Resource isolation
- Easy scaling and orchestration
- Simplified dependency management

### 2. Orchestration

**Kubernetes Deployment**:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: b3-challenge-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: b3-challenge-api
  template:
    spec:
      containers:
      - name: api
        image: b3-challenge:latest
        ports:
        - containerPort: 8080
        env:
        - name: DB_HOST
          value: "postgres-service"
```

**Features**:
- Automatic scaling based on metrics
- Rolling updates with zero downtime
- Health check integration
- Service discovery and load balancing

### 3. CI/CD Pipeline

**Pipeline Stages**:
1. Code checkout and dependency installation
2. Unit test execution and coverage reporting
3. Integration test execution
4. Performance benchmark execution
5. Docker image building and scanning
6. Deployment to staging environment
7. Automated testing in staging
8. Production deployment approval
9. Production deployment and monitoring

## Future Architecture Considerations

### 1. Microservices Evolution

**Service Decomposition**:
- Trade Data Service
- Aggregation Service
- User Management Service
- Notification Service
- Analytics Service

**Benefits**:
- Independent scaling and deployment
- Technology diversity
- Fault isolation
- Team autonomy

### 2. Event-Driven Architecture

**Event Streaming**:
- Apache Kafka for event streaming
- Event sourcing for audit trails
- CQRS for read/write separation
- Real-time data processing

### 3. Advanced Caching

**Multi-Level Caching**:
- CDN for static content
- Redis for application caching
- Database query result caching
- Client-side caching strategies

### 4. Machine Learning Integration

**Predictive Analytics**:
- Price prediction models
- Volume forecasting
- Anomaly detection
- Market trend analysis

---

This architecture provides a solid foundation for a high-performance, scalable trading data API while maintaining flexibility for future enhancements and requirements.

