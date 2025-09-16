# Golang OTP Authentication API

A robust Golang backend service implementing OTP-based authentication with user management, JWT tokens, rate limiting, and containerized deployment.

## 🚀 Features

- **OTP Authentication**: Secure mobile number-based authentication using One-Time Passwords
- **User Management**: Complete user registration and profile management
- **JWT Tokens**: Access and refresh token implementation for secure API access
- **Rate Limiting**: Built-in rate limiting for OTP requests to prevent abuse
- **Clean Architecture**: Modular design with separation of concerns (handler, usecase, repository layers)
- **Swagger Documentation**: Interactive API documentation
- **Docker Support**: Fully containerized application with Docker Compose
- **Multi-Environment Config**: Support for development, production, and Docker environments

## 🏗️ Architecture

The project follows Clean Architecture principles with the following layers:

```
src/
├── cmd/                    # Application entry point
├── internal/
│   ├── middlewares/        # HTTP middlewares (CORS, rate limiting)
│   └── user/
│       ├── api/
│       │   ├── dto/        # Data Transfer Objects
│       │   ├── handler/    # HTTP handlers
│       │   ├── router/     # Route definitions
│       │   └── validation/ # Input validation
│       ├── repository/     # Data access layer
│       └── usecase/        # Business logic layer
├── pkg/
│   ├── cache/             # Redis cache implementation
│   ├── config/            # Configuration management
│   ├── db/                # Database connection
│   └── helper/            # Utility functions
├── migrations/            # Database migrations
└── docs/                  # Swagger documentation
```

## 🛠️ Tech Stack

- **Language**: Go 1.23.0
- **Web Framework**: Gin
- **Database**: PostgreSQL with GORM ORM
- **Cache**: Redis
- **Authentication**: JWT tokens
- **Documentation**: Swagger/OpenAPI
- **Containerization**: Docker & Docker Compose

## 📋 Prerequisites

- Go 1.23.0 or higher
- PostgreSQL 12+
- Redis 6+
- Docker & Docker Compose (for containerized setup)

## 🚀 How to Run Locally

### 1. Clone the Repository

```bash
git clone https://github.com/alielmi98/golang-otp-auth.git
cd golang-otp-auth
```

### 2. Setup Database

**PostgreSQL:**
```bash
# Create database
createdb user_db

# Or using psql
psql -U postgres
CREATE DATABASE user_db;
```

**Redis:**
```bash
# Start Redis server
redis-server

# Or using Docker
docker run -d -p 6379:6379 redis:latest
```

### 3. Configure Environment

Update the configuration file `src/pkg/config/config-development.yml`:

```yaml
server:
  internalPort: 5005
  externalPort: 5005
  runMode: debug

postgres:
  host: localhost
  port: 5432
  user: postgres
  password: admin
  dbName: user_db
  sslMode: disable

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

otp:
  expireTime: 120  # seconds
  digits: 6
  limiter: 100     # seconds between requests

jwt:
  secret: "your-secret-key"
  refreshSecret: "your-refresh-secret"
  accessTokenExpireDuration: 60    # minutes
  refreshTokenExpireDuration: 1440 # minutes (24 hours)
```

### 4. Install Dependencies

```bash
cd src
go mod download
```

### 5. Run the Application

```bash
cd src
go run cmd/main.go
```

The API will be available at `http://localhost:5005`

### 6. Access Swagger Documentation

Open your browser and navigate to:
```
http://localhost:5005/swagger/index.html
```

## 🐳 How to Run with Docker

### 1. Using Docker Compose (Recommended)

```bash
# Navigate to docker directory
cd docker

# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

This will start:
- **Backend API** on port `5000`
- **PostgreSQL** on port `5432`
- **Redis** on port `6379`

### 2. Build and Run Manually

```bash
# Build the Docker image
cd src
docker build -t golang-otp-auth .

# Run with environment variables
docker run -d \
  -p 5000:5000 \
  -e APP_ENV=docker \
  golang-otp-auth
```

### 3. Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_ENV` | Environment mode (development/docker/production) | development |
| `PORT` | External port override | - |

## 📚 API Documentation

### Base URL
- **Local**: `http://localhost:5005/api/v1`
- **Docker**: `http://localhost:5000/api/v1`

### Authentication
Most endpoints require JWT authentication. Include the token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

### API Endpoints

#### 1. Send OTP
**POST** `/users/send-otp`

Sends an OTP to the specified mobile number.

**Request:**
```bash
curl -X POST "http://localhost:5005/api/v1/users/send-otp" \
  -H "Content-Type: application/json" \
  -d '{
    "mobile_number": "09123456789"
  }'
```

**Response:**
```json
{
  "result": null,
  "success": true,
  "resultCode": 0,
  "error": null
}
```

#### 2. Register/Login with Mobile & OTP
**POST** `/users/login-by-mobile`

Register a new user or login existing user using mobile number and OTP.

**Request:**
```bash
curl -X POST "http://localhost:5005/api/v1/users/login-by-mobile" \
  -H "Content-Type: application/json" \
  -d '{
    "mobileNumber": "09123456789",
    "otp": "123456"
  }'
```

**Response:**
```json
{
  "result": {
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "accessTokenExpireTime": 1640995200,
    "refreshTokenExpireTime": 1641081600
  },
  "success": true,
  "resultCode": 0,
  "error": null
}
```

#### 3. Get User by Mobile Number
**GET** `/users/{mobile_number}`

Retrieve user information by mobile number.

**Request:**
```bash
curl -X GET "http://localhost:5005/api/v1/users/09123456789" \
  -H "Authorization: Bearer <your-jwt-token>"
```

**Response:**
```json
{
  "id": 1,
  "mobile_number": "09123456789",
  "registered_at": "2024-01-15T10:30:00Z"
}
```

#### 4. Get Users (Paginated)
**GET** `/users`

Retrieve a paginated list of users with optional filtering.

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `page_size` (optional): Items per page (default: 10)
- `mobile_number` (optional): Filter by mobile number

**Request:**
```bash
curl -X GET "http://localhost:5005/api/v1/users?page=1&page_size=10&mobile_number=0912" \
  -H "Authorization: Bearer <your-jwt-token>"
```

**Response:**
```json
{
  "users": [
    {
      "id": 1,
      "mobile_number": "09123456789",
      "registered_at": "2024-01-15T10:30:00Z"
    }
  ],
  "total": 1,
  "page": 1,
  "page_size": 10
}
```

### Error Responses

All endpoints return consistent error responses:

```json
{
  "result": null,
  "success": false,
  "resultCode": 40001,
  "error": {
    "message": "Validation error",
    "details": "Invalid mobile number format"
  }
}
```

**Common Result Codes:**
- `0`: Success
- `40001`: Validation Error
- `40101`: Authentication Error
- `40401`: Not Found Error
- `42901`: Rate Limiter Error
- `42902`: OTP Rate Limiter Error
- `50001`: Internal Server Error

## 🗄️ Database Choice Justification

### PostgreSQL (Primary Database)
**Why PostgreSQL was chosen:**

1. **ACID Compliance**: Ensures data consistency and reliability for user accounts and authentication data
2. **Mature Ecosystem**: Excellent Go support through GORM, extensive documentation, and community
3. **Scalability**: Handles concurrent connections well, supports read replicas for scaling
4. **Data Integrity**: Strong typing system and constraints prevent data corruption
5. **JSON Support**: Native JSON columns for flexible data storage when needed
6. **Performance**: Excellent query optimization and indexing capabilities
7. **Reliability**: Battle-tested in production environments, excellent backup and recovery tools

**Use Cases in this project:**
- User account storage (ID, mobile number, registration date)
- Authentication logs and audit trails
- User profile data and preferences
- Persistent data that requires ACID guarantees

### Redis (Cache & Session Store)
**Why Redis was chosen:**

1. **Speed**: In-memory storage provides microsecond latency for OTP operations
2. **TTL Support**: Built-in expiration perfect for temporary OTP codes (120 seconds)
3. **Atomic Operations**: Ensures thread-safe OTP generation and validation
4. **Rate Limiting**: Efficient sliding window rate limiting implementation
5. **Session Management**: Fast JWT token blacklisting and session storage
6. **Pub/Sub**: Future extensibility for real-time notifications
7. **Memory Efficiency**: Optimized data structures for caching scenarios

**Use Cases in this project:**
- OTP code storage with automatic expiration
- Rate limiting counters (prevent OTP spam)
- JWT token blacklisting for logout functionality
- Caching frequently accessed user data
- Session management and temporary data

### Database Architecture Benefits

This **dual-database approach** provides:

- **Performance**: Fast reads from Redis, reliable writes to PostgreSQL
- **Scalability**: Independent scaling of cache and persistent storage
- **Reliability**: Redis failures don't affect core user data
- **Cost Efficiency**: Use expensive SSD storage only for persistent data
- **Flexibility**: Different optimization strategies for different data types

## 🔧 Configuration

The application supports multiple environments through YAML configuration files:

- `config-development.yml`: Local development
- `config-docker.yml`: Docker containerized environment  
- `config-production.yml`: Production deployment

Key configuration sections:

### Server Configuration
```yaml
server:
  internalPort: 5005      # Internal application port
  externalPort: 5005      # External exposed port
  runMode: debug          # Gin mode: debug/release
```

### Database Configuration
```yaml
postgres:
  host: localhost
  port: 5432
  user: postgres
  password: admin
  dbName: user_db
  sslMode: disable
  maxIdleConns: 15        # Connection pool settings
  maxOpenConns: 100
  connMaxLifetime: 5      # Minutes
```

### OTP Configuration
```yaml
otp:
  expireTime: 120         # OTP expiration in seconds
  digits: 6               # OTP length
  limiter: 100            # Rate limit window in seconds
```

## 🧪 Testing

### Manual Testing with curl

1. **Send OTP:**
```bash
curl -X POST "http://localhost:5005/api/v1/users/send-otp" \
  -H "Content-Type: application/json" \
  -d '{"mobile_number": "09123456789"}'
```

2. **Login with OTP:**
```bash
curl -X POST "http://localhost:5005/api/v1/users/login-by-mobile" \
  -H "Content-Type: application/json" \
  -d '{"mobileNumber": "09123456789", "otp": "123456"}'
```

3. **Get User Info:**
```bash
curl -X GET "http://localhost:5005/api/v1/users/09123456789" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Using Swagger UI

Navigate to `http://localhost:5005/swagger/index.html` for interactive API testing.

## 🚀 Deployment

### Production Deployment

1. **Update production config** in `config-production.yml`
2. **Set environment variables:**
```bash
export APP_ENV=production
export PORT=8080
```

3. **Build and deploy:**
```bash
go build -o auth-api cmd/main.go
./auth-api
```

### Docker Production Deployment

```bash
# Build production image
docker build -t golang-otp-auth:prod .

# Run with production config
docker run -d \
  -p 8080:5000 \
  -e APP_ENV=production \
  -e PORT=5000 \
  golang-otp-auth:prod
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📞 Support

For support and questions:
- Create an issue on GitHub
- Contact: [alielmi98](https://github.com/alielmi98)

---

**Happy Coding! 🚀**
