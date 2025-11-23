# GoFi (Golang with Fiber)

GoFi is a robust and scalable web application boilerplate built with Go and the Fiber v2 framework. It follows Clean Architecture principles to ensure maintainability and testability. This project provides a solid foundation for building RESTful APIs with built-in authentication, role-based access control, and database management.

## Features

-   **High Performance**: Built on [Fiber](https://gofiber.io/), one of the fastest Go web frameworks.
-   **Clean Architecture**: Separation of concerns for better maintainability.
-   **Authentication**: Secure user sign-up, sign-in, and session management using JWT.
-   **RBAC**: Role-Based Access Control to manage user permissions.
-   **Database**: PostgreSQL integration with `pgx` driver.
-   **Migrations**: Database schema management using `golang-migrate`.
-   **User Management**: Full CRUD operations for users with soft delete functionality.
-   **Configuration**: Environment-based configuration using `.envrc`.

## Tech Stack

-   **Language**: [Go](https://go.dev/) (1.20+)
-   **Framework**: [Fiber v2](https://github.com/gofiber/fiber)
-   **Database**: [PostgreSQL](https://www.postgresql.org/)
-   **Migrations**: [golang-migrate](https://github.com/golang-migrate/migrate)
-   **Containerization**: [Docker](https://www.docker.com/)

## Prerequisites

Before you begin, ensure you have the following installed:

-   [Go](https://go.dev/dl/) (version 1.20 or higher)
-   [Docker](https://docs.docker.com/get-docker/) & Docker Compose
-   [Make](https://www.gnu.org/software/make/)

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/masb0ymas/gofi.git
cd gofi
```

### 2. Configuration

Copy the example environment file and configure your variables:

```bash
cp .envrc.example .envrc
# Edit .envrc with your specific configuration (DB credentials, secrets, etc.)
source .envrc # Or allow direnv if you use it
```

### 3. Database Setup

Ensure your PostgreSQL database is running. You can use Docker to start a Postgres instance if you don't have one locally:

```bash
docker run --name gofi-postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=gofi -p 5432:5432 -d postgres
```

Run the database migrations to set up the schema:

```bash
make db/migrations/up
```

To seed the database with initial data (if available):

```bash
make db/migrations/up/seed
```

### 4. Running the Application

Start the API server:

```bash
make run
```

The server will start on the port specified in your `.envrc` (default is `8080`).

## Project Structure

```
gofi/
├── cmd/
│   ├── api/            # Main entry point for the API
│   └── migrate/        # Entry point for migration tool
├── internal/
│   ├── app/            # App bootstrapping
│   ├── config/         # Configuration loading
│   ├── handlers/       # HTTP handlers (Controllers)
│   ├── middlewares/    # HTTP middlewares
│   ├── models/         # Domain models
│   ├── repositories/   # Data access layer
│   ├── services/       # Business logic
│   └── ...
├── migrations/         # SQL migration files
├── Makefile            # Build and utility commands
└── README.md           # Project documentation
```

## API Endpoints

### Authentication
-   `POST /v1/auth/sign-up`: Register a new user
-   `POST /v1/auth/sign-in`: Login
-   `POST /v1/auth/verify-registration`: Verify email registration
-   `GET /v1/auth/verify-session`: Check current session status
-   `POST /v1/auth/sign-out`: Logout

### Users
-   `GET /v1/users`: List users
-   `GET /v1/users/:userID`: Get user details
-   `POST /v1/users`: Create user (Admin only)
-   `PUT /v1/users/:userID`: Update user (Admin only)
-   `DELETE /v1/users/:userID`: Delete user (Admin only)

### Roles
-   `GET /v1/roles`: List roles
-   `POST /v1/roles`: Create role (Admin only)
-   ...and more.

## Useful Commands

| Command | Description |
| :--- | :--- |
| `make run` | Run the application locally |
| `make build/api` | Build the API binary |
| `make db/migrations/new name=foo` | Create a new migration file |
| `make db/migrations/up` | Apply up migrations |
| `make db/migrations/down` | Rollback migrations |
| `make help` | Show all available make commands |

## Creator

N. Fajri (masb0ymas)
<br />
me@masb0ymas.com

## License

This project is licensed under the [MIT License](LICENSE.md).
