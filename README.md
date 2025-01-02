## Gofi (Golang Fiber)

### Prerequisites
- Go 1.19+
- Docker
- Docker Compose

### Setup
1. Clone the repository: `git clone https://github.com/masb0ymas/gofi.git`
2. Navigate to the project directory: `cd gofi`
3. Run the setup script: `bash setup.sh`
4. Create a database according to `.env`
5. Install golang-migrate: `go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest`
6. Run migrations: `make migration-up`

### Usage
1. Navigate to the project directory: `cd gofi`
2. Run the server: `make dev`
3. For production, run the server with command: `make start`
4. Access the API: `http://localhost:8000`