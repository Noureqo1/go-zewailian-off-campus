# Zewailian Off Campus

[![Go Version](https://img.shields.io/badge/Go-1.19+-00ADD8?style=flat&logo=go)](https://golang.org/doc/go1.19)
[![Next.js Version](https://img.shields.io/badge/Next.js-13.0+-000000?style=flat&logo=next.js)](https://nextjs.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A real-time chat and blogging platform designed specifically for Zewail City for Science and Technology students. The platform enables students to create chat rooms, share achievements, upload educational materials, make announcements, and participate in Q&A discussions.

## 🚀 Features

- **Real-time Chat System**
  - Create and join chat rooms
  - Real-time messaging with WebSockets
  - Message history and persistence
  - User presence indicators

- **User Authentication**
  - OAuth2 integration with Google
  - JWT-based authentication
  - Role-based access control

- **Microservices Architecture**
  - Scalable and resilient design
  - Service isolation and independent deployment
  - API Gateway for service orchestration

- **Database & Caching**
  - PostgreSQL for persistent storage
  - Redis for caching and real-time features

## 🛠️ Technology Stack

### Backend
- **Go** - Core backend language
- **Gin** - Web framework
- **Gorilla WebSockets** - WebSocket implementation
- **PostgreSQL** - Primary database
- **Redis** - Caching and pub/sub
- **JWT** - Authentication

### Frontend
- **Next.js** - React framework
- **TypeScript** - Type-safe JavaScript
- **Tailwind CSS** - Utility-first CSS framework

### DevOps
- **Docker** - Containerization
- **Docker Compose** - Multi-container orchestration
- **OpenSearch** - Centralized logging

## 📋 Prerequisites

- Go 1.19+
- Node.js 18+
- PostgreSQL 14+
- Redis 6+
- Docker & Docker Compose

## 🔧 Installation

### Clone the repository
```bash
git clone https://github.com/Noureqo1/go-zewailian-off-campus.git
cd go-zewailian-off-campus
```

### Backend Setup
```bash
cd server

# Set up environment variables
cp .env.example .env
# Edit .env with your configuration

# Install dependencies
go mod download

# Run the server
go run cmd/main.go
```

### Frontend Setup
```bash
cd client

# Install dependencies
npm install

# Run development server
npm run dev
```

### Using Docker
```bash
# Build and start all services
docker-compose up -d
```

## 🏗️ Project Structure

```
├── client/                 # Frontend Next.js application
│   ├── app/                # Next.js app directory
│   └── public/             # Static assets
├── docs/                   # Documentation
│   ├── design/             # System design documents
│   │   ├── ADRs/           # Architecture Decision Records
│   │   ├── C4/             # C4 model diagrams
│   │   └── api-spec.yaml   # API specifications
├── server/                 # Backend Go application
│   ├── cmd/                # Application entry points
│   ├── db/                 # Database connection and migrations
│   ├── internal/           # Internal packages
│   │   ├── oauth/          # OAuth authentication
│   │   ├── user/           # User management
│   │   └── ws/             # WebSocket implementation
│   ├── router/             # API routes
│   └── util/               # Utility functions
└── docker-compose.yml      # Docker compose configuration
```

## 🔒 Authentication

The application uses OAuth 2.0 with Google for authentication. JWT tokens are used for maintaining user sessions and securing API endpoints.

## 🌐 API Documentation

API documentation is available in OpenAPI format at `docs/design/api-spec.yaml`. You can view it using tools like Swagger UI or Redoc.

## 🧪 Testing

```bash
# Run backend tests
cd server
go test ./...

# Run frontend tests
cd client
npm test
```

## 📚 Documentation

Detailed documentation is available in the `docs` directory:

- System requirements and specifications
- Architecture diagrams (C4 model)
- Architecture Decision Records (ADRs)
- API specifications
- Resilience strategies

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 📧 Contact

Project Link: [https://github.com/Noureqo1/go-zewailian-off-campus](https://github.com/Noureqo1/go-zewailian-off-campus)

