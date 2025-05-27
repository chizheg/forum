# Forum Project

A microservices-based forum application with real-time chat functionality.

## Architecture

The project consists of two microservices:
1. **Auth Service** - Handles user authentication and authorization
2. **Forum Service** - Manages forum functionality and chat

### Key Features
- Real-time chat using WebSocket
- PostgreSQL database without ORM
- gRPC communication between services
- Database migrations using golang-migrate
- Clean Architecture principles
- Comprehensive test coverage (80%+)
- API documentation using Swagger
- Logging using Zap

## Prerequisites

- Go 1.21+
- PostgreSQL
- Docker (optional)

## Project Structure

```
.
├── cmd/
│   ├── auth-service/
│   └── forum-service/
├── internal/
│   ├── auth/
│   └── forum/
├── pkg/
│   ├── logger/
│   └── database/
├── migrations/
├── proto/
└── docs/
```

## Setup

1. Clone the repository:
```bash
git clone https://github.com/yourusername/forum.git
```

2. Install dependencies:
```bash
go mod download
```

3. Set up PostgreSQL databases:
```bash
# Create databases
createdb forum_auth
createdb forum_main
```

4. Run migrations:
```bash
make migrate-up
```

5. Start the services:
```bash
# Start auth service
make run-auth

# Start forum service
make run-forum
```

## Testing

Run tests with:
```bash
make test
```

## API Documentation

After starting the services, Swagger documentation is available at:
- Auth Service: http://localhost:8081/swagger/index.html
- Forum Service: http://localhost:8082/swagger/index.html 