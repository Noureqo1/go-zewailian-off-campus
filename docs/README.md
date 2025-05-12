# Project Documentation

This directory contains comprehensive documentation for the Go-Zewailian Off-Campus project. The documentation is organized into several sections to help you understand the project's architecture, design decisions, and implementation details.

## 📁 Directory Structure

- `design/` - Architecture and design documentation
  - C4 model diagrams
  - System architecture
  - Component interactions
  - Database schema

- `Google_OAuth0.2/` - OAuth Implementation
  - Google OAuth 2.0 setup
  - Authentication flow
  - Security considerations

- `Testing/` - Testing Documentation
  - Unit testing guidelines
  - Integration testing
  - Test coverage reports
  - Testing strategies

- `Postman_Tests/` - API Testing
  - Postman collection
  - API endpoints documentation
  - Request/response examples

## 🏗️ Architecture Overview

The project follows a clean architecture pattern with clear separation of concerns:

1. **Frontend (Next.js)**
   - React components
   - TypeScript implementation
   - Tailwind CSS styling
   - WebSocket client integration

2. **Backend (Go)**
   - RESTful APIs
   - WebSocket server
   - PostgreSQL database
   - Redis caching
   - Google OAuth integration

## 🔍 Key Documentation Files

- `design/architecture.md` - System architecture and design patterns
- `design/database.md` - Database schema and relationships
- `Google_OAuth0.2/setup.md` - OAuth configuration guide
- `Testing/guidelines.md` - Testing standards and practices

## 🚀 Getting Started with Documentation

1. Start with the architecture overview in `design/`
2. Review the API documentation for endpoint details
3. Check testing guidelines for development
4. Refer to OAuth setup for authentication implementation

## 📝 Contributing to Documentation

When contributing to the documentation:

1. Follow the existing markdown format
2. Update diagrams when making architectural changes
3. Keep API documentation in sync with implementation
4. Add examples for new features

## 🔄 Documentation Updates

The documentation is regularly updated to reflect:

- New feature implementations
- Architecture changes
- API modifications
- Testing requirements

For any questions or suggestions about the documentation, please open an issue in the repository.

---

<img src="../docs/design/architecture-diagram.png" width="600" alt="Architecture Overview">

*Note: This documentation is maintained by the Go-Zewailian Off-Campus team and is regularly updated to reflect the current state of the project.*
