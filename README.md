# Clients-Service

## Overview

This service is part of the `AI-Hackathon-2026` ecosystem.  
It exposes HTTP endpoints for:

- User registration and login
- Token refresh and logout (logout handler is currently a no-op)
- Profile read/update
- Graph list retrieval by username

## Architecture

The project follows a layered design:

- `transport` - HTTP handlers, routing, middleware, JSON helpers
- `service` - business logic (auth, profile, streaks)
- `repository` - data access abstraction over an in-memory fake database
- `model` - request/response and domain structures
- `pkg` - reusable utilities (`jwtutil`, `hashutil`)
- `app` - dependency wiring, server startup, logging, graceful shutdown

## Tech Stack

- Go 1.25.x
- Standard `net/http` router (`http.ServeMux`)
- JWT: `github.com/golang-jwt/jwt/v5`
- Password hashing: `bcrypt` (`golang.org/x/crypto`)
- UUID: `github.com/google/uuid`
- Optional `.env` loader: `github.com/joho/godotenv`

## Project Structure

```text
Clients-Service/
  cmd/ClientsService/main.go
  internal/
    app/
    config/
    model/
    repository/
    service/
    transport/
  pkg/
    hashutil/
    jwtutil/
  configs/.env
  Dockerfile
```

## Configuration

Environment variables:

| Variable | Required | Default | Description |
| --- | --- | --- | --- |
| `PORT` | No | `9191` | HTTP server port |
| `JWT_SECRET` | Yes | empty | Secret used to sign and validate JWT tokens |

Example `configs/.env`:

```env
PORT=9191
JWT_SECRET=your_super_secret_key
```

## Running Locally

```bash
go mod tidy
go run ./cmd/ClientsService
```

Service starts on `http://localhost:9191` by default.

## Running with Docker

```bash
docker build -t clients-service .
docker run --rm -p 9191:9191 --env JWT_SECRET=your_super_secret_key clients-service
```

## API Endpoints

### Public

| Method | Path | Description |
| --- | --- | --- |
| `POST` | `/auth/register` | Register user and return access + refresh tokens |
| `POST` | `/auth/login` | Login and return access + refresh tokens |
| `POST` | `/auth/refresh` | Refresh token pair |

### Protected (Bearer token required)

| Method | Path | Description |
| --- | --- | --- |
| `POST` | `/auth/logout` | Logout endpoint (currently returns `204`, token invalidation is not implemented yet) |
| `GET` | `/profile-edit/{username}` | Get profile by username |
| `PUT` | `/profile-edit/` | Update username/email |
| `GET` | `/graphs/{username}` | Get graph list for user |
| `GET` | `/get-streak` | Get current user's streak |

Authorization header format:

```http
Authorization: Bearer <access_token>
```

## Request Examples

Register:

```bash
curl -X POST http://localhost:9191/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"john","email":"john@example.com","password":"secret123"}'
```

Login:

```bash
curl -X POST http://localhost:9191/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","password":"secret123"}'
```

Get streak:

```bash
curl http://localhost:9191/get-streak \
  -H "Authorization: Bearer <access_token>"
```

## Behavior Notes

- Storage is in-memory (`FakeDB`), so all data is lost after restart.
