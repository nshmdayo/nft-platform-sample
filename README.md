# NFT Platform for Academic Papers - Backend

Backend API for NFT platform for academic papers

## Overview

This project is a backend for an NFT platform for academic papers implemented in Go. It provides a system where researchers can submit papers, undergo peer review, and publish them as NFTs.

## Features

- **User Authentication**: User registration and login with JWT authentication
- **Paper Management**: Creation, editing, deletion, and search of papers
- **Peer Review System**: Peer review and scoring functionality
- **NFT Integration**: NFT conversion after paper review completion (in preparation)

## Technology Stack

- **Language**: Go 1.22+
- **Web Framework**: Gin
- **Database**: PostgreSQL with GORM
- **Authentication**: JWT
- **Logging**: slog (JSON format)

## Architecture

Adopting Clean Architecture pattern:
- `models/`: Data models
- `repository/`: Data access layer
- `service/`: Business logic layer
- `handlers/`: HTTP handlers
- `middleware/`: Middleware
- `router/`: Routing configuration

## Setup

### Environment Variables

Create a `.env` file and configure the following environment variables:

```env
# Database
DATABASE_URL=postgres://user:password@localhost:5432/nft_platform?sslmode=disable

# Server
PORT=8080
ENVIRONMENT=development
GIN_MODE=debug

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRES_IN=24h

# IPFS (for future implementation)
IPFS_API_URL=http://localhost:5001

# Ethereum (for future implementation)
ETHEREUM_RPC_URL=http://localhost:8545
ETHEREUM_PRIVATE_KEY=your-private-key
PAPER_CONTRACT_ADDRESS=0x...
REVIEW_CONTRACT_ADDRESS=0x...
```

### Database Setup

Start PostgreSQL and create the database:

```bash
createdb nft_platform
```

### Application Execution

```bash
# Install dependencies
go mod tidy

# Run application
go run cmd/server/main.go
```

## API Endpoints

### Authentication

- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - Login
- `GET /api/v1/auth/profile` - Get profile (authentication required)

### Papers

- `POST /api/v1/papers` - Create paper (authentication required)
- `GET /api/v1/papers` - Get paper list
- `GET /api/v1/papers/my` - Get my papers (authentication required)
- `GET /api/v1/papers/search?q=query` - Search papers
- `GET /api/v1/papers/:id` - Get paper details
- `PUT /api/v1/papers/:id` - Update paper (authentication required)
- `DELETE /api/v1/papers/:id` - Delete paper (authentication required)
- `POST /api/v1/papers/:id/submit` - Submit for review (authentication required)

### Reviews

- `POST /api/v1/reviews` - Create review (authentication required)
- `GET /api/v1/reviews/my` - Get my reviews (authentication required)
- `GET /api/v1/reviews/pending` - Get pending review papers (authentication required)
- `GET /api/v1/reviews/:id` - Get review details (authentication required)
- `PUT /api/v1/reviews/:id` - Update review (authentication required)
- `DELETE /api/v1/reviews/:id` - Delete review (authentication required)
- `GET /api/v1/papers/:paper_id/reviews` - Get paper reviews (authentication required)
- `GET /api/v1/papers/:paper_id/score` - Get paper score (authentication required)

## Project Structure

```
cmd/
  server/
    main.go           # Application entry point
internal/
  config/
    config.go         # Configuration management
  database/
    connection.go     # Database connection
  handlers/
    auth_handler.go   # Authentication handler
    paper_handler.go  # Paper handler
    review_handler.go # Review handler
  middleware/
    auth.go          # Authentication middleware
    middleware.go    # Other middleware
  models/
    user.go          # User model
    paper.go         # Paper model
    review.go        # Review model
    nft_metadata.go  # NFT metadata model
  repository/
    user_repository.go    # User repository
    paper_repository.go   # Paper repository
    review_repository.go  # Review repository
  router/
    router.go        # Routing configuration
  service/
    auth_service.go  # Authentication service
    paper_service.go # Paper service
    review_service.go # Review service
  utils/
    jwt.go           # JWT utility
    password.go      # Password hash
pkg/
  logger/
    logger.go        # Log configuration
```

## Future Implementation Plans

- [ ] IPFS integration for file storage
- [ ] Ethereum smart contract integration
- [ ] NFT minting functionality
- [ ] File upload functionality
- [ ] Notification system
- [ ] Detailed access control

## Development

### Running Tests

```bash
go test ./...
```

### Build

```bash
go build -o bin/server cmd/server/main.go
```

## License

MIT License