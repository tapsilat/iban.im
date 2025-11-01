# Docker Development Setup

This project supports Docker-based development with hot reloading using [Air](https://github.com/air-verse/air).

## Quick Start

### Option 1: PostgreSQL (Recommended)

Start the development environment with PostgreSQL database:

```bash
docker compose -f docker-compose.dev.yml up
```

The application will be available at `http://localhost:8080`

Features:
- PostgreSQL 18 database
- Hot reloading with Air
- Persistent database volume
- GraphQL playground at `/graph`

### Option 2: SQLite (Lightweight)

Start the development environment with SQLite database:

```bash
docker compose -f docker-compose.dev.sqlite.yml up
```

Features:
- SQLite database (file-based)
- Hot reloading with Air
- Persistent database volume
- Lighter resource usage

## How It Works

### Hot Reloading

The development setup uses [Air](https://github.com/air-verse/air) for hot reloading:

1. Air watches for file changes in `.go`, `.tmpl`, and `.html` files
2. When changes are detected, it automatically rebuilds and restarts the application
3. Build errors are logged to `tmp/build-errors.log`

Configuration is in `.air.toml`. Customize it to:
- Change watched file extensions
- Modify excluded directories
- Adjust build commands

### Volume Mounts

The compose files mount your local code directory into the container at `/app`, so any changes you make locally are immediately reflected in the container.

## Database Access

### PostgreSQL

- **Host**: `localhost` (from your machine) or `db` (from within Docker network)
- **Port**: `5432`
- **Database**: `ibanim`
- **User**: `ibanim`
- **Password**: `ibanim`

Connect with psql:
```bash
psql -h localhost -p 5432 -U ibanim -d ibanim
```

Or using Docker:
```bash
docker compose -f docker-compose.dev.yml exec db psql -U ibanim -d ibanim
```

### SQLite

The SQLite database file is stored in a Docker volume at `/app/data/ibanim.db`.

Access the SQLite database:
```bash
docker compose -f docker-compose.dev.sqlite.yml exec dev sqlite3 /app/data/ibanim.db
```

## Common Commands

### View Logs

```bash
# All services
docker compose -f docker-compose.dev.yml logs -f

# Just the app
docker compose -f docker-compose.dev.yml logs -f dev

# Just the database
docker compose -f docker-compose.dev.yml logs -f db
```

### Rebuild Containers

If you change `go.mod` or `Dockerfile.dev`:

```bash
docker compose -f docker-compose.dev.yml up --build
```

### Stop Services

```bash
docker compose -f docker-compose.dev.yml down
```

### Remove Volumes (Fresh Start)

```bash
# PostgreSQL
docker compose -f docker-compose.dev.yml down -v

# SQLite
docker compose -f docker-compose.dev.sqlite.yml down -v
```

### Run Commands in Container

```bash
# Open a shell
docker compose -f docker-compose.dev.yml exec dev sh

# Run tests
docker compose -f docker-compose.dev.yml exec dev go test ./...

# Run linter
docker compose -f docker-compose.dev.yml exec dev go vet ./...
```

## Environment Variables

You can customize the environment by creating a `.env` file in the project root. See `.env.example` for available options.

For Docker development, the most important variables are already set in the `docker-compose.dev.yml` files:

- `DB_ADAPTER`: Database type (`postgres` or `sqlite`)
- `DB_HOST`: Database hostname (`db` for PostgreSQL in Docker)
- `DB_PORT`: Database port
- `DB_NAME`: Database name or file path
- `DB_USER`: Database user
- `DB_PASSWORD`: Database password
- `APP_PORT`: Application port (default: `8080`)
- `APP_ENV`: Environment (`development`)
- `APP_DEBUG`: Debug mode (`true`)
- `APP_KEY`: Secret key for JWT tokens
- `APP_REALM`: JWT realm

## Troubleshooting

### Port Already in Use

If port 8080 or 5432 is already in use, modify the port mapping in `docker-compose.dev.yml`:

```yaml
ports:
  - "8081:8080"  # Map to a different local port
```

### Database Connection Issues

1. Check if the database container is healthy:
```bash
docker compose -f docker-compose.dev.yml ps
```

2. Check database logs:
```bash
docker compose -f docker-compose.dev.yml logs db
```

3. The application waits for the database to be healthy before starting (PostgreSQL only).

### Hot Reload Not Working

1. Check if Air is watching the correct files (see `.air.toml`)
2. Check build errors in `tmp/build-errors.log`
3. Restart the services:
```bash
docker compose -f docker-compose.dev.yml restart dev
```

### Module Download Issues

If you see errors about downloading Go modules, ensure you have a working internet connection and try rebuilding:

```bash
docker compose -f docker-compose.dev.yml build --no-cache
```

## Production Deployment

This Docker setup is for **development only**. For production deployment, use the regular `Dockerfile` and `docker-compose.yml`:

```bash
docker compose up --build
```

The production Dockerfile creates an optimized, minimal container using multi-stage builds.
