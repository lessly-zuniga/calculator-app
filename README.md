# Calculator App

## Project Overview

Calculator App is a full-stack web calculator with a React and TypeScript frontend and a Go REST API backend. It supports basic and advanced arithmetic while keeping calculation rules separate from HTTP and presentation concerns.

## Features

- Addition, subtraction, multiplication, and division
- Percentage, square root, and exponentiation
- Calculator-style contextual percentage behavior
- Responsive calculator interface for mobile, tablet, and desktop viewports
- Request and operand validation
- Structured API error handling
- Frontend unit tests and backend domain and HTTP tests
- Dockerized full-stack execution through Docker Compose

## Tech Stack

### Frontend

- React
- TypeScript
- Vite
- Vitest
- ESLint
- Prettier

### Backend

- Go
- `net/http`
- Go standard-library testing tools

### Infrastructure

- Docker
- Docker Compose
- Nginx
- GitHub Actions

## Architecture

React components render the calculator and keep display and keypad concerns separate. The `useCalculator` hook owns interaction state and coordinates operations through the calculator API client. The client uses relative `/api` URLs and performs no arithmetic.

In production, Nginx serves the built React application and proxies `/api/` requests to the Go service. Go HTTP handlers validate and translate requests and responses, while the calculator package contains HTTP-independent domain logic.

```text
Browser
  -> Nginx / React
  -> /api/*
  -> Go REST API
  -> calculator domain
```

## Project Structure

```text
calculator-app/
├── .github/
│   └── workflows/
│       └── ci.yml
├── backend/
│   ├── cmd/server/
│   │   ├── main.go
│   │   └── main_test.go
│   ├── internal/
│   │   ├── calculator/
│   │   └── handler/
│   ├── Dockerfile
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   ├── calculator/
│   │   │   └── layout/
│   │   ├── hooks/
│   │   ├── services/
│   │   └── types/
│   ├── Dockerfile
│   ├── nginx.conf
│   ├── package.json
│   ├── vite.config.ts
│   └── vitest.config.ts
├── docker-compose.yml
└── README.md
```

## Quick Start with Docker

Build and start both services:

```bash
docker compose up --build
```

Open [http://localhost:3000](http://localhost:3000).

Stop and remove the containers and Compose network:

```bash
docker compose down
```

The backend is available only inside the Compose network. Browser API requests go through the frontend Nginx container.

## Local Development

### Backend

Start the Go API:

```bash
cd backend
go run ./cmd/server
```

The API runs at `http://localhost:8080`. Verify it with:

```bash
curl http://localhost:8080/api/v1/health
```

### Frontend

Install dependencies and start Vite:

```bash
cd frontend
npm ci
npm run dev
```

The development application is available at `http://localhost:5173` by default. Vite proxies relative `/api` requests to `http://localhost:8080`, so the backend must also be running.

## API Documentation

### Endpoints

| Method | Endpoint | Operation |
| --- | --- | --- |
| `GET` | `/api/v1/health` | Health check |
| `POST` | `/api/v1/add` | Addition |
| `POST` | `/api/v1/subtract` | Subtraction |
| `POST` | `/api/v1/multiply` | Multiplication |
| `POST` | `/api/v1/divide` | Division |
| `POST` | `/api/v1/percentage` | Percentage conversion |
| `POST` | `/api/v1/square-root` | Square root |
| `POST` | `/api/v1/power` | Exponentiation |

The health endpoint returns:

```json
{
  "status": "ok"
}
```

Binary operations accept exactly two numeric operands. For example:

```http
POST /api/v1/add
Content-Type: application/json
```

```json
{
  "operands": [2, 3]
}
```

```json
{
  "result": 5
}
```

Percentage and square root accept exactly one numeric operand. For example:

```http
POST /api/v1/square-root
Content-Type: application/json
```

```json
{
  "operands": [81]
}
```

```json
{
  "result": 9
}
```

### Error Responses

Malformed requests, non-numeric operands, and incorrect binary operand counts return HTTP 400:

```json
{
  "error": {
    "code": "INVALID_REQUEST",
    "message": "Exactly two numeric operands are required"
  }
}
```

Unary endpoints use the message `Exactly one numeric operand is required`.

Division by zero returns HTTP 422:

```json
{
  "error": {
    "code": "DIVISION_BY_ZERO",
    "message": "Cannot divide by zero"
  }
}
```

A negative square root returns HTTP 422:

```json
{
  "error": {
    "code": "NEGATIVE_SQUARE_ROOT",
    "message": "Cannot calculate the square root of a negative number"
  }
}
```

A calculation that produces `NaN` or infinity returns HTTP 422:

```json
{
  "error": {
    "code": "INVALID_RESULT",
    "message": "Calculation produced an invalid numeric result"
  }
}
```

## Percentage Behavior

The frontend implements calculator-style contextual percentages while keeping arithmetic backend-driven:

| Input | Result |
| --- | ---: |
| `10%` | `0.1` |
| `200 + 10%` | `220` |
| `200 - 10%` | `180` |
| `200 × 10%` | `20` |
| `200 ÷ 10%` | `2000` |

For addition and subtraction, the percentage is relative to the first operand. For multiplication and division, it is converted to its decimal value. The hook coordinates calls to the percentage and arithmetic API endpoints; it does not calculate final results locally.

## Testing

### Frontend

```bash
cd frontend
npm test
npm run test:coverage
```

The current hook test suite reports 92.3% statement coverage, 77.61% branch coverage, 100% function coverage, and 92.3% line coverage for `useCalculator.ts`.

### Backend

```bash
cd backend
go test ./...
go test -cover ./...
```

The current package coverage is:

- `cmd/server`: 72.1%
- `internal/calculator`: 100.0%
- `internal/handler`: 84.8%

## Code Quality

Frontend checks:

```bash
cd frontend
npm run lint
npm run format:check
npm run build
```

Backend checks:

```bash
cd backend
gofmt -w .
go vet ./...
```

GitHub Actions runs frontend and backend validation in separate CI jobs on pushes and pull requests.

## Design Decisions

- Frontend presentation and interaction state are separated from networking.
- Final arithmetic is performed by the Go backend.
- Relative `/api` URLs work with both the Vite development proxy and production Nginx.
- Nginx serves the SPA and provides a single browser-facing origin for API traffic.
- React local state is sufficient; no global state library is used.
- The backend uses the Go standard library without unnecessary external dependencies.
- The responsive visual design is Apple-inspired in clarity and interaction quality, but is an original web interface rather than an Apple clone.

## AI Usage

AI tooling was used to accelerate implementation, testing, and review. Architecture, requirements, behavior, validation, and final decisions remained developer-directed and were manually verified.

See [AI usage notes](docs/ai-usage.md) for additional details.

## Future Improvements

- Add browser-level end-to-end coverage for critical calculator flows.
- Add automated container smoke tests to CI.
- Track frontend bundle size and runtime performance as the application evolves.
