# Dormitory Management System — Microservice Starter Template

A starter template for a Dormitory Management System built with a **microservice architecture** using Node.js and Express.

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
| `api-gateway` | 3000 | Routes all incoming requests to the appropriate service |
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
# In separate terminals:
cd student-service && npm install && npm start
cd room-service    && npm install && npm start
cd booking-service && npm install && npm start
cd api-gateway     && npm install && npm start
```

## API Endpoints

All requests go through the **API Gateway** on port `3000`.

### Students — `/students`

| Method | Path | Description |
|---|---|---|
| GET | `/students` | List all students |
| GET | `/students/:id` | Get a student by ID |
| POST | `/students` | Create a new student |
| PUT | `/students/:id` | Update a student |
| DELETE | `/students/:id` | Remove a student |

**POST body example:**
```json
{ "name": "Alice Smith", "email": "alice@example.com", "phone": "555-0101" }
```

### Rooms — `/rooms`

| Method | Path | Description |
|---|---|---|
| GET | `/rooms` | List all rooms (supports `?available=true`) |
| GET | `/rooms/:id` | Get a room by ID |
| POST | `/rooms` | Create a new room |
| PUT | `/rooms/:id` | Update a room |
| DELETE | `/rooms/:id` | Remove a room |

**POST body example:**
```json
{ "number": "301", "type": "single", "capacity": 1 }
```

### Bookings — `/bookings`

| Method | Path | Description |
|---|---|---|
| GET | `/bookings` | List all bookings |
| GET | `/bookings/:id` | Get a booking by ID |
| POST | `/bookings` | Create a new booking |
| PUT | `/bookings/:id` | Update a booking |
| DELETE | `/bookings/:id` | Cancel a booking |

**POST body example:**
```json
{ "studentId": 1, "roomId": 2, "checkIn": "2026-03-01", "checkOut": "2026-06-30" }
```

### Health Checks

Each service exposes `GET /health`. Through the gateway:

```
GET http://localhost:3000/students/../health  →  student-service
GET http://localhost:3002/health              →  room-service (direct)
GET http://localhost:3003/health              →  booking-service (direct)
```

## Project Structure

```
├── docker-compose.yml
├── api-gateway/
│   ├── Dockerfile
│   ├── package.json
│   └── src/index.js          # Proxy routes to downstream services
├── student-service/
│   ├── Dockerfile
│   ├── package.json
│   └── src/
│       ├── index.js
│       └── routes/students.js
├── room-service/
│   ├── Dockerfile
│   ├── package.json
│   └── src/
│       ├── index.js
│       └── routes/rooms.js
└── booking-service/
    ├── Dockerfile
    ├── package.json
    └── src/
        ├── index.js
        └── routes/bookings.js
```

## Extending the Template

* **Replace the in-memory stores** in each `routes/*.js` file with a real database (e.g. MongoDB, PostgreSQL).
* **Add authentication** to the API Gateway (e.g. JWT middleware).
* **Add inter-service communication** using HTTP calls or a message broker (e.g. RabbitMQ).
* **Add a new service** by following the same pattern and registering a new proxy route in `api-gateway/src/index.js`.
