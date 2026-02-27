# Dormitory Management System — Microservice Starter Template

A starter template for a Dormitory Management System built with a **microservice architecture** using **Go** and the standard library (`net/http`).

## Architecture

```
                        ┌─────────────────┐
  Client ──────────────▶│   API Gateway   │  :3000
                        └────────┬────────┘
                                 │
              ┌──────────────────┼──────────────────┐
              ▼                  ▼                  ▼
   ┌──────────────────┐ ┌──────────────┐ ┌──────────────────┐
   │ student-service  │ │ room-service │ │ booking-service  │
   │      :3001       │ │    :3002     │ │      :3003       │
   └──────────────────┘ └──────────────┘ └──────────────────┘
```

| Service | Port | Responsibility |
|---|---|---|
| `api-gateway` | 3000 | Reverse-proxies all incoming requests to the appropriate service |
| `student-service` | 3001 | Manage student / resident records |
| `room-service` | 3002 | Manage dormitory rooms |
| `booking-service` | 3003 | Manage room bookings (check-in / check-out) |

## Quick Start

### With Docker Compose (recommended)

```bash
docker-compose up --build
```

### Without Docker (run each service manually)

```bash
# In separate terminals (requires Go 1.22+):
cd student-service && go run .
cd room-service    && go run .
cd booking-service && go run .
cd api-gateway     && go run .
```

### Build binaries

```bash
cd student-service && go build -o student-service .
cd room-service    && go build -o room-service    .
cd booking-service && go build -o booking-service .
cd api-gateway     && go build -o api-gateway     .
```

## API Endpoints

All requests go through the **API Gateway** on port `3000`.

### Students — `/students`

| Method | Path | Description |
|---|---|---|
| GET | `/students` | List all students |
| GET | `/students/{id}` | Get a student by ID |
| POST | `/students` | Create a new student |
| PUT | `/students/{id}` | Update a student |
| DELETE | `/students/{id}` | Remove a student |

**POST body example:**
```json
{ "name": "Alice Smith", "email": "alice@example.com", "phone": "555-0101" }
```

### Rooms — `/rooms`

| Method | Path | Description |
|---|---|---|
| GET | `/rooms` | List all rooms (supports `?available=true`) |
| GET | `/rooms/{id}` | Get a room by ID |
| POST | `/rooms` | Create a new room |
| PUT | `/rooms/{id}` | Update a room |
| DELETE | `/rooms/{id}` | Remove a room |

**POST body example:**
```json
{ "number": "301", "type": "single", "capacity": 1 }
```

### Bookings — `/bookings`

| Method | Path | Description |
|---|---|---|
| GET | `/bookings` | List all bookings |
| GET | `/bookings/{id}` | Get a booking by ID |
| POST | `/bookings` | Create a new booking |
| PUT | `/bookings/{id}` | Update a booking |
| DELETE | `/bookings/{id}` | Cancel a booking (soft cancel) |

**POST body example:**
```json
{ "studentId": 1, "roomId": 2, "checkIn": "2026-03-01", "checkOut": "2026-06-30" }
```

### Health Checks

Each service exposes `GET /health`:

```
GET http://localhost:3000/health  →  api-gateway
GET http://localhost:3001/health  →  student-service (direct)
GET http://localhost:3002/health  →  room-service (direct)
GET http://localhost:3003/health  →  booking-service (direct)
```

## Project Structure

```
├── docker-compose.yml
├── api-gateway/
│   ├── Dockerfile
│   ├── go.mod
│   └── main.go          # httputil.ReverseProxy routing
├── student-service/
│   ├── Dockerfile
│   ├── go.mod
│   └── main.go          # CRUD handlers + in-memory store
├── room-service/
│   ├── Dockerfile
│   ├── go.mod
│   └── main.go          # CRUD handlers + in-memory store
└── booking-service/
    ├── Dockerfile
    ├── go.mod
    └── main.go          # CRUD handlers + in-memory store
```

## Extending the Template

* **Replace the in-memory stores** in each `main.go` with a real database (e.g. PostgreSQL with `database/sql`, or MongoDB with the official Go driver).
* **Add authentication** to the API Gateway (e.g. JWT middleware).
* **Add inter-service communication** using HTTP calls or a message broker (e.g. RabbitMQ, NATS).
* **Add a new service** by following the same pattern and registering a new proxy route in `api-gateway/main.go`.

