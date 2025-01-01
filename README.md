# Identity Service

This service handles user management, authentication, and authorization for the microservices architecture.

## Features

- User registration and management
- JWT-based authentication
- Role-based access control (RBAC)
- Password reset functionality
- Event-driven notifications via Kafka

## Tech Stack

- Go 1.21
- PostgreSQL
- Redis
- Kafka
- Docker

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.21 or later (for local development)

### Local Development

1. Clone the repository:
```bash
git clone https://github.com/mibrahim2344/identity-service.git
```

2. Start the services:
```bash
docker-compose up -d
```

3. The service will be available at `http://localhost:8080`

## API Documentation

### Endpoints

- POST /api/v1/register - User registration
- POST /api/v1/login - User login
- POST /api/v1/refresh - Refresh access token
- POST /api/v1/reset-password - Password reset
- GET /api/v1/me - Get current user

## Testing

Run the tests:
```bash
go test ./...
```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request
