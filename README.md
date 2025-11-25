# GoFi

> A high-performance REST API built with Go and Gin framework

## 📋 Overview

GoFi is a production-ready REST API server built with Go, designed for high performance and cost-effectiveness. The codebase follows clean architecture principles with dependency injection as its core design pattern.

### Key Features

- **High Performance**: Built on the Gin web framework with optimized middleware
- **Clean Architecture**: Dependency injection and separation of concerns
- **Database Migrations**: Automated schema management with golang-migrate
- **Security**: JWT authentication, Helmet middleware, CORS support
- **Developer Experience**: Hot reload support, comprehensive Make commands
- **Production Ready**: Docker support, cross-compilation, CI/CD ready

### Tech Stack

- **Framework**: [Fiber](https://github.com/gofiber/fiber/v2) v2.52.9
- **Database**: PostgreSQL with [lib/pq](https://github.com/lib/pq)
- **Migrations**: [golang-migrate](https://github.com/golang-migrate/migrate)
- **Authentication**: JWT with [golang-jwt](https://github.com/golang-jwt/jwt)
- **Email**: [Resend](https://resend.com) integration
- **Middleware**: CORS, Gzip, Rate Limiting, Request ID, Helmet

## 🚀 Getting Started

### Prerequisites

- Go 1.24.5 or higher
- PostgreSQL 18 or higher (using uuid v7)
- Make (optional but recommended)

### Installation

1. **Clone the repository**

```bash
git clone https://github.com/masb0ymas/gofi.git
cd gofi
```

2. **Set up environment variables**

```bash
cp .envrc.example .envrc
```

Edit `.envrc` and configure your environment variables:

```bash
# Required configurations
export PORT=8080
export ENV=development
export DEBUG=true
export APP_NAME=gofi

# Database connection
export DB_DSN=postgres://postgres:postgres@localhost:5432/gofi?sslmode=disable

# JWT secret (generate a secure random string, with `openssl rand -base64 32`)
export JWT_SECRET=your-secret-key-here

# Application URLs
export CLIENT_URL=http://localhost:3000
export SERVER_URL=http://localhost:8080

# Optional: Email service (Resend)
export RESEND_API_KEY=your-resend-api-key
export RESEND_FROM_EMAIL=noreply@yourdomain.com
```

3. **Set up the database**

Create your PostgreSQL database, then update the database name in `/migrations/000001_initial-database.up.sql` (search for `dev_gintama` and replace with your database name).

4. **Run migrations and seed data**

```bash
make db/migrations/up/seed
```

5. **Start the application**

```bash
make run
```

The API will be available at `http://localhost:8080`

## 📚 Available Commands

### General

```bash
# List all available commands
make help
```

### Development

```bash
# Run the application
make run

# Build the application
make build/api
```

### Database

```bash
# Connect to the database using psql
make db/psql

# Create a new database migration
make db/migrations/new name=create_users_table

# Apply all pending migrations
make db/migrations/up

# Apply migrations and seed test data
make db/migrations/up/seed

# Rollback all migrations
make db/migrations/down

# Refresh database (down + up)
make db/migrations/refresh

# Refresh database with seed data
make db/migrations/refresh/seed
```

## 🏗️ Project Structure

```
gofi/
├── cmd/
│   ├── api/          # API server entry point
│   └── migrate/      # Migration CLI tool
├── internal/         # Private application code
│   ├── handlers/     # HTTP request handlers
│   ├── models/       # Data models
│   ├── repository/   # Data access layer
│   └── services/     # Business logic
├── migrations/       # Database migration files
├── public/           # Static files
├── script/           # Utility scripts
├── templates/        # Email/HTML templates
├── .envrc.example    # Environment variables template
├── Dockerfile        # Container configuration
├── Makefile          # Build and development commands
└── go.mod            # Go module dependencies
```

## 🐳 Docker Support

Build and run with Docker:

```bash
# Build the Docker image
docker build -t gofi:latest .

# Run the container
docker run -d -p 8080:8080 gofi:latest
```

## 🚢 Deployment

### Automated Deployment

Set up GitHub Actions by configuring the workflows in `.github/workflows`. The CI/CD pipeline can push images to:

- Google Container Registry (GCR)
- AWS Elastic Container Registry (ECR)
- Docker Hub
- Any OCI-compatible registry

### Manual Deployment

Build the application for your target environment:

```bash
make build/api
```

This generates two binaries:

- **`/bin/api`** - Compiled for your local machine's architecture
- **`/bin/linux_amd64/api`** - Cross-compiled for Linux AMD64 (production servers)

Deploy the appropriate binary to your server and run it with the required environment variables.

### Environment Variables for Production

Ensure these are properly configured in production:

- `ENV=production`
- `DEBUG=false`
- `JWT_SECRET` - Use a strong, randomly generated secret
- `DB_DSN` - Production database connection string
- `CLIENT_URL` - Your frontend application URL
- `SERVER_URL` - Your API server URL

## 🔒 Security

- JWT-based authentication
- Helmet middleware for security headers
- CORS configuration
- Rate limiting support
- Environment-based configuration
- SQL injection protection via parameterized queries

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

## 👤 Author

GitHub: [@masb0ymas](https://github.com/masb0ymas)
<br/>
Email: [me@masb0ymas.com](mailto:me@masb0ymas.com)
<br/>

Credit: [@edwardanthony](https://github.com/edwardanthony)

## ⭐ Show Your Support

Give a ⭐️ if this project helped you!
