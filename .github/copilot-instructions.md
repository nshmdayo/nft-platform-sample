# Academic Paper NFT Platform Specification

## Project Overview

A decentralized NFT platform that revolutionizes the publication, peer review, and distribution of academic papers. The platform aims to enable researchers to publish papers as NFTs, transparentize the peer review process, and properly manage intellectual property rights.

## Technology Stack

- **Language**: Go 1.21+
- **Framework**: Gin (Web API)
- **Database**: PostgreSQL
- **Blockchain**: Ethereum (Goerli testnet)
- **NFT Library**: go-ethereum
- **Authentication**: JWT
- **Storage**: IPFS (academic paper file storage)
- **Containerization**: Docker & Docker Compose

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   Backend       │    │   Blockchain    │
│   (Future)      │<-->│   (Go API)      │<-->│   (Ethereum)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              │
                       ┌─────────────────┐
                       │   PostgreSQL    │
                       │   Database      │
                       └─────────────────┘
                              │
                              │
                       ┌─────────────────┐
                       │     IPFS        │
                       │ (File Storage)  │
                       └─────────────────┘
```

## Directory Structure

```
.
├── cmd/
│   └── server/
│       └── main.go                    # Application entry point
├── internal/
│   ├── api/
│   │   ├── handlers/                  # HTTP handlers
│   │   │   ├── auth.go
│   │   │   ├── papers.go
│   │   │   ├── reviews.go
│   │   │   └── nft.go
│   │   ├── middleware/                # Middleware
│   │   │   ├── auth.go
│   │   │   ├── cors.go
│   │   │   └── logging.go
│   │   └── routes/
│   │       └── routes.go              # Routing definition
│   ├── blockchain/
│   │   ├── client.go                  # Ethereum client
│   │   ├── contracts/                 # Smart contracts
│   │   │   ├── paper_nft.sol
│   │   │   └── review_nft.sol
│   │   └── service.go                 # Blockchain service
│   ├── config/
│   │   └── config.go                  # Configuration management
│   ├── database/
│   │   ├── migrations/                # DB migrations
│   │   └── connection.go              # DB connection
│   ├── models/
│   │   ├── user.go
│   │   ├── paper.go
│   │   ├── review.go
│   │   └── nft.go
│   ├── repository/
│   │   ├── user_repository.go
│   │   ├── paper_repository.go
│   │   └── review_repository.go
│   ├── service/
│   │   ├── auth_service.go
│   │   ├── paper_service.go
│   │   ├── review_service.go
│   │   └── nft_service.go
│   └── utils/
│       ├── jwt.go
│       ├── hash.go
│       └── validation.go
├── pkg/
│   ├── ipfs/
│   │   └── client.go                  # IPFS client
│   └── logger/
│       └── logger.go                  # Logger
├── scripts/
│   ├── deploy-contracts.go            # Contract deployment
│   └── setup-db.go                    # DB initialization
├── docker-compose.yml                 # Development environment
├── Dockerfile                         # App container
├── go.mod
├── go.sum
└── README.md
```

## Data Models

### User
```go
type User struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    Email       string    `json:"email" gorm:"unique;not null"`
    Name        string    `json:"name" gorm:"not null"`
    WalletAddr  string    `json:"wallet_address" gorm:"unique"`
    Role        string    `json:"role"`    // researcher, reviewer, admin
    Institution string    `json:"institution"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### Paper
```go
type Paper struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    Title       string    `json:"title" gorm:"not null"`
    Abstract    string    `json:"abstract"`
    Authors     []string  `json:"authors" gorm:"type:json"`
    Keywords    []string  `json:"keywords" gorm:"type:json"`
    Category    string    `json:"category"`
    IPFSHash    string    `json:"ipfs_hash"`    // IPFS file hash
    NFTTokenID  *uint     `json:"nft_token_id"` // Token ID when NFT is minted
    OwnerID     uint      `json:"owner_id"`
    Status      string    `json:"status"`       // draft, submitted, under_review, published
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    
    Owner   User     `json:"owner" gorm:"foreignKey:OwnerID"`
    Reviews []Review `json:"reviews"`
}
```

### Review
```go
type Review struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    PaperID     uint      `json:"paper_id"`
    ReviewerID  uint      `json:"reviewer_id"`
    Score       int       `json:"score"`        // 1-10
    Comments    string    `json:"comments"`
    Status      string    `json:"status"`       // pending, completed, rejected
    NFTTokenID  *uint     `json:"nft_token_id"` // NFT of review results
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    
    Paper    Paper `json:"paper" gorm:"foreignKey:PaperID"`
    Reviewer User  `json:"reviewer" gorm:"foreignKey:ReviewerID"`
}
```

### NFTMetadata
```go
type NFTMetadata struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    TokenID     uint      `json:"token_id" gorm:"unique"`
    Type        string    `json:"type"`         // paper, review
    ReferenceID uint      `json:"reference_id"` // Paper ID or Review ID
    MetadataURI string    `json:"metadata_uri"` // IPFS URI
    TxHash      string    `json:"tx_hash"`      // Transaction hash
    CreatedAt   time.Time `json:"created_at"`
}
```

## API Specification

### Authentication Related
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - Login
- `POST /api/auth/refresh` - Token refresh

### Paper Related
- `GET /api/papers` - Get paper list
- `GET /api/papers/:id` - Get paper details
- `POST /api/papers` - Submit paper
- `PUT /api/papers/:id` - Update paper
- `POST /api/papers/:id/mint` - Mint paper as NFT
- `POST /api/papers/:id/submit` - Submit for review

### Review Related
- `GET /api/reviews` - Get review list
- `GET /api/reviews/:id` - Get review details
- `POST /api/reviews` - Submit review
- `PUT /api/reviews/:id` - Update review
- `POST /api/reviews/:id/mint` - Mint review result as NFT

### NFT Related
- `GET /api/nfts` - Get NFT list
- `GET /api/nfts/:tokenId` - Get NFT details
- `POST /api/nfts/transfer` - Transfer NFT

## Smart Contract Specification

### PaperNFT Contract
```solidity
contract PaperNFT is ERC721 {
    struct PaperMetadata {
        string title;
        string ipfsHash;
        address author;
        uint256 timestamp;
    }
    
    mapping(uint256 => PaperMetadata) public papers;
    
    function mintPaper(
        address to,
        string memory title,
        string memory ipfsHash
    ) public returns (uint256);
}
```

### ReviewNFT Contract
```solidity
contract ReviewNFT is ERC721 {
    struct ReviewMetadata {
        uint256 paperId;
        address reviewer;
        uint8 score;
        string ipfsHash;
        uint256 timestamp;
    }
    
    mapping(uint256 => ReviewMetadata) public reviews;
    
    function mintReview(
        address to,
        uint256 paperId,
        uint8 score,
        string memory ipfsHash
    ) public returns (uint256);
}
```

## Development Guidelines

### Coding Standards
- Follow Go's standard naming conventions
- Functions should follow the single responsibility principle
- Implement proper error handling
- Target test coverage of 80% or higher

### File Naming Conventions
- Use `snake_case`
- Test files use `_test.go` suffix
- Interfaces use `I` prefix (e.g., `IUserRepository`)

### Database Design
- Migration files use `YYYYMMDD_HHMMSS_description.sql` format
- Set foreign key constraints appropriately
- Utilize indexes effectively

### Security
- Proper JWT token management
- Password hashing (using bcrypt)
- SQL injection prevention
- Proper CORS configuration management

## Environment Variables

```bash
# Database
DATABASE_URL=postgres://user:password@localhost:5432/nft_platform

# Ethereum
ETHEREUM_RPC_URL=https://goerli.infura.io/v3/YOUR_PROJECT_ID
PRIVATE_KEY=your_ethereum_private_key
CONTRACT_ADDRESS_PAPER=0x...
CONTRACT_ADDRESS_REVIEW=0x...

# IPFS
IPFS_API_URL=http://localhost:5001

# JWT
JWT_SECRET=your_jwt_secret
JWT_EXPIRES_IN=24h

# Server
PORT=8080
GIN_MODE=debug
```

## Implementation Order

1. **Phase 1**: Basic Infrastructure Setup
   - Project initialization
   - Database configuration
   - Basic API structure

2. **Phase 2**: Authentication System
   - User registration/login
   - JWT authentication
   - Middleware implementation

3. **Phase 3**: Paper Management
   - Paper CRUD operations
   - IPFS file upload
   - Paper search functionality

4. **Phase 4**: NFT Integration
   - Smart contract development
   - Blockchain integration
   - NFT minting functionality

5. **Phase 5**: Review System
   - Review process implementation
   - Reviewer matching
   - Review result NFT minting

## Test Strategy

- Unit tests: Independent testing of each layer
- Integration tests: API endpoint testing
- E2E tests: Scenario-based testing
- Blockchain tests: Using Ganache

## Deployment

- Containerized deployment using Docker
- CI/CD pipeline setup
- Testing in staging environment
- Gradual deployment to production
