# social-go

A REST API for a social network built with Go. It provides user registration, JWT authentication, posts, comments, followers, and a personalized feed.

## Stack

- **Go** with [Chi](https://github.com/go-chi/chi) router
- **PostgreSQL** for persistence
- **Redis** for caching and rate limiting
- **JWT** for auth, **Mailtrap** for activation emails
- **Swagger** for API docs at `/v1/swagger/`

## Getting started

1. Start dependencies:

```bash
docker compose up -d
```

2. Set environment variables (see `.envrc` for reference) and run migrations:

```bash
make migrate-up
```

3. Run the API:

```bash
go run ./cmd/api
```

The server listens on `:3030` by default.

## Scripts

| Command | Description |
|---------|-------------|
| `make migrate-up` | Run database migrations |
| `make migrate-down` | Roll back last migration |
| `make seed` | Seed the database |
| `make gen-docs` | Regenerate Swagger docs |
| `make test` | Run tests |
