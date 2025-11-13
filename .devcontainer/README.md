# Development Container Setup

This devcontainer provides a complete development environment for the Go Skeleton project with all required dependencies.

## What's Included

### Development Environment
- **Go 1.24** - Latest Go version with all tools
- **Docker-in-Docker** - Run Docker commands inside the container
- **Zsh with Oh My Zsh** - Enhanced shell experience

### Infrastructure Services
- **MySQL 8.0** - Primary database (port 3306)
  - Database: `go_skeleton`
  - User: `root` / Password: `root`
  
- **PostgreSQL 15** - Alternative database (port 5432)
  - Database: `go_skeleton`
  - User: `postgres` / Password: `postgres`
  
- **Redis 7** - Caching and session storage (port 6379)
  - No password configured by default
  
- **RabbitMQ 3** - Message queue (ports 5672, 15672)
  - User: `guest` / Password: `guest`
  - Management UI: http://localhost:15672
  
- **MongoDB 6** - Document database (port 27017)
  - Database: `go_skeleton`

### VS Code Extensions
- Go language support
- Docker extension
- Code spell checker
- GitLens
- MongoDB support
- SQL Tools for MySQL and PostgreSQL

## Getting Started

1. **Open in VS Code**
   - Install the "Dev Containers" extension in VS Code
   - Open the project folder
   - Click "Reopen in Container" when prompted, or use Command Palette: `Dev Containers: Reopen in Container`

2. **Wait for Setup**
   - The container will build and start all services
   - Go dependencies will be downloaded automatically
   - This may take a few minutes on first run

3. **Verify Services**
   ```bash
   # Check MySQL
   mysql -h db -u root -proot -e "SHOW DATABASES;"
   
   # Check PostgreSQL
   psql -h postgres -U postgres -c "\l"
   
   # Check Redis
   redis-cli -h redis ping
   
   # Check RabbitMQ
   curl -u guest:guest http://rabbitmq:15672/api/overview
   
   # Check MongoDB
   mongosh --host mongodb --eval "db.version()"
   ```

4. **Run the Application**
   ```bash
   # Copy environment file
   cp .devcontainer/.env.devcontainer .env
   
   # Run database migrations
   make migrate-up  # or your migration command
   
   # Start the API server
   go run cmd/api/main.go
   ```

## Environment Configuration

The `.devcontainer/.env.devcontainer` file contains pre-configured environment variables for the devcontainer services. Copy this file to `.env` in your project root:

```bash
cp .devcontainer/.env.devcontainer .env
```

Key differences from local development:
- Database hosts use service names (e.g., `db`, `postgres`, `redis`)
- Ports are internal Docker network ports
- All services are on the same Docker network

## Port Forwarding

The following ports are automatically forwarded to your local machine:
- **7011** - API Server
- **3306** - MySQL
- **5432** - PostgreSQL
- **6379** - Redis
- **5672** - RabbitMQ AMQP
- **15672** - RabbitMQ Management UI
- **27017** - MongoDB

## Troubleshooting

### Services not responding
```bash
# Check if services are running
docker-compose -f .devcontainer/docker-compose.yml ps

# View service logs
docker-compose -f .devcontainer/docker-compose.yml logs [service-name]
```

### Reset everything
```bash
# Remove volumes and rebuild
docker-compose -f .devcontainer/docker-compose.yml down -v
# Then rebuild the container in VS Code
```

### Go module issues
```bash
go mod tidy
go mod download
```

## Tips

- Use the integrated terminal in VS Code for all commands
- RabbitMQ Management UI provides a visual interface for queue management
- Data persists in Docker volumes between container rebuilds
- Use `make` commands if available in the project's Makefile
