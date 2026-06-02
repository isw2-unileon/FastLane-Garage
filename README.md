# FastLane Garage - Full Stack Auto Parts Management

A complete full-stack application for managing automotive parts and customer orders, built with **Go** backend and **React** frontend in a monorepo structure.

## 🎯 Features

### Backend (Gon + Gin + GORM)

✅ **Parts Management**

- CRUD operations for car parts
- Filtering by car zone
- Full-text search by part name
- Price management
- Image URLs support

✅ **Order Management**

- Complete order lifecycle (pending -> processing -> completed)
- Order items with quantity and unit pricing
- Customer email validation
- Order status tracking
- Automatic total price calculation

✅ **Architecture**

- Clean architecture (Handlers -> Service -> Repository)
- Dependency injection
- Interface-based design for testability
- Comprehensive unit tests with mocks
- Error handling with context wrapping

✅ **Quality**

- 100% test coverage for service layer
- Go linters (golangci-lint)
- Structured logging with slog
- Type-safe DTOs with validation
- GORM with SQLite database

### Frontend (React + Three + TypeScript + Vite)

---

## 🏗️ Project Structure

```text
├── backend/                    # Go API Server
│   ├── cmd/server/            # Entry point
│   │   └── main.go            # Server initialization
│   └── internal/
│       ├── config/            # Configuration management
│       ├── database/          # Database migrations & seeding
│       ├── models/            # Domain models (Part, Order, OrderItem)
│       ├── dto/               # Data transfer objects with validation
│       ├── repository/        # Data access layer (PartsRepository, OrdersRepository)
│       ├── service/           # Business logic (PartsService, OrdersService)
│       ├── handlers/          # HTTP handlers (Parts, Orders, Helpers)
│
├── frontend/                   # React + TypeScript + Vite
│   └── src/                   # React components (to be implemented)
│
├── e2e/                       # Playwright E2E tests
│   └── tests/
│
├── docs/                      # Documentation
│   ├── getting-started.md     # Setup guide
│   ├── golang.md              # Go best practices
│   ├── monorepo.md            # Monorepo explanation
│   └── adr/                   # Architecture Decision Records
│
└── .github/workflows/         # CI/CD pipelines
```

## 📊 API Endpoint

### Parts Endpoints

| Method   | Path             | Description     | Query Parameters         |
| -------- | ---------------- | --------------- | ------------------------ |
| `GET`    | `/api/parts`     | List all parts  | `?zone=motor&name=Motor` |
| `GET`    | `/api/parts/:id` | Get part by ID  | -                        |
| `POST`   | `/api/parts`     | Create new part | -                        |
| `PUT`    | `/api/parts/:id` | Update part     | -                        |
| `DELETE` | `/api/parts/:id` | Delete part     | -                        |

**Examples:**

```bash
# Get all parts
curl http://localhost:8080/api/parts

# Filter by zone
curl http://localhost:8080/api/parts?zone=motor

# Search by name
curl http://localhost:8080/api/parts?name=Motor

# Get specific part
curl http://localhost:8080/api/parts/1
```

### Orders Endpoints

| Method   | Path                     | Description         | Query Parameters                         |
| -------- | ------------------------ | ------------------- | ---------------------------------------- |
| `GET`    | `/api/orders`            | List all orders     | `?status=pending&email=user@example.com` |
| `GET`    | `/api/orders/:id`        | Get order by ID     | -                                        |
| `POST`   | `/api/orders`            | Create new order    | -                                        |
| `PUT`    | `/api/orders/:id/status` | Update order status | -                                        |
| `DELETE` | `/api/orders/:id`        | Delete order        | -                                        |

**Examples:**

```bash
# Get all orders
curl http://localhost:8080/api/orders

# Filter by status
curl http://localhost:8080/api/orders?status=pending

# Search by email
curl http://localhost:8080/api/orders?email=customer@example.com

# Create order
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_name": "Juan García",
    "customer_email": "juan@example.com",
    "items": [
      {"part_id": 1, "quantity": 2},
      {"part_id": 3, "quantity": 1}
    ]
  }'
```

---

## 🚀 Getting Started

### Prerequisites

- [Go](https://go.dev/dl/) 1.24+
- [Node.js](https://nodejs.org/) 22+

### Installation

```bash
# Clone repository
git clone https://github.com/isw2-unileon/FastLane-Garage.git
cd FastLane-Garage

# Install all dependencies
make install
```

### Development

```bash
# Terminal 1: Start backend (port 8080)
make run-backend

# Terminal 2: Start frontend (port 5173)
make run-frontend
```

The Vite dev server automatically proxies `/api` requests to the backend at `http://localhost:8080`.

### Testing & Quality

```bash
# Run all tests
make test

# Run linters
make lint

# Run E2E tests (requires services running)
make e2e

# Build for production
make build-backend
make build-frontend
```

---

## 🧪 Testing

### Unit Tests

The project includes comprehensive unit tests for all service layers:

```bash
go test -v -race ./backend/internal/service/
```

**Test Coverage:**

- ✅ Parts Service: 6 tests (GetAll, GetById, Create, Update, Delete, Validation)
- ✅ Orders Service: 10 tests (GetAll, GetById, Create, Update Status, Delete, Filters, Validation)

Run with coverage:

```bash
go test -v -race -coverprofile=coverage.out ./...
```

---

## 📦 Database Schema

### Parts Table

```sql
CREATE TABLE parts (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL,
  car_zone TEXT NOT NULL,
  image_url TEXT,
  price REAL NOT NULL
);
```

### Orders Table

```sql
CREATE TABLE orders (
  id INTEGER PRIMARY KEY,
  customer_name TEXT NOT NULL,
  customer_email TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'pending',
  total_price REAL NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Order Items Table

```sql
CREATE TABLE order_items (
  id INTEGER PRIMARY KEY,
  order_id INTEGER NOT NULL,
  part_id INTEGER NOT NULL,
  quantity INTEGER NOT NULL,
  unit_price REAL NOT NULL,
  FOREIGN KEY (order_id) REFERENCES orders(id),
  FOREIGN KEY (part_id) REFERENCES parts(id)
);
```

---

## 🏛️ Architecture

### Clean Architecture Layers

```
HTTP Request
    ↓
[Handlers]        ← HTTP request/response handling
    ↓
[Service]         ← Business logic & validation
    ↓
[Repository]      ← Data access & persistence
    ↓
[Database]        ← SQLite (GORM)
```

### Design Patterns Used

1. **Repository Pattern** - Abstract data access
2. **Service Pattern** - Encapsulate business logic
3. **Dependency Injection** - Loose coupling via constructors
4. **Interface-based Design** - For testability and flexibility
5. **DTO Pattern** - Request/response validation and transformation
6. **Mock Pattern** - Isolated unit testing

---

## 📚 Documentation

- [Getting Started Guide](docs/getting-started.md) - How to set up and use the template
- [Go Best Practices](docs/golang.md) - Guidelines for backend development
- [Monorepo Structure](docs/monorepo.md) - Why we use a monorepo
- [Architecture Decision Records](docs/adr/) - Design decisions and rationale

---

## 🔧 Available Commands

```bash
make install         # Install all dependencies
make run-backend     # Backend with hot reload (Air)
make run-frontend    # Frontend dev server (Vite)
make build-backend   # Build Go binary
make build-frontend  # Build frontend for production
make test            # Run all tests
make lint            # Run linters (Go + Frontend)
make e2e             # Run Playwright E2E tests
```

---

## 🚦 CI/CD Pipeline

The project includes automated testing and quality checks:

| Workflow      | Trigger                | What it does                |
| ------------- | ---------------------- | --------------------------- |
| `Backend CI`  | Push/PR to `backend/`  | Lint + Test + Build Go      |
| `Frontend CI` | Push/PR to `frontend/` | ESLint + TypeScript + Build |
| `CodeQL`      | Weekly + Push/PR       | Security analysis           |
| `E2E Tests`   | Manual dispatch        | Playwright tests            |

---

## 💾 Database Seeding

The application automatically seeds sample data on first run:

**Sample Parts:**

- Motor V6 (motor zone, €2500)
- Turbo (motor zone, €1200)
- Neumático Michelin (neumaticos, €150)
- Disco de freno (frenos, €120)
- Pastillas de freno (frenos, €80)
- Puerta delantera (puertas, €300)
- Faro LED (iluminacion, €250)
- Batería 12V (electrico, €180)

---

## 📖 Code Examples

### Creating a Part

```bash
curl -X POST http://localhost:8080/api/parts \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Motor V8",
    "car_zone": "motor",
    "price": 3500.00,
    "image_url": "https://example.com/motor-v8.jpg"
  }'
```

### Creating an Order

```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_name": "Juan García",
    "customer_email": "juan@example.com",
    "items": [
      {"part_id": 1, "quantity": 1},
      {"part_id": 2, "quantity": 2}
    ]
  }'
```

### Updating Order Status

```bash
curl -X PUT http://localhost:8080/api/orders/1/status \
  -H "Content-Type: application/json" \
  -d '{
    "status": "completed"
  }'
```

---

## 🤝 Contributing

When adding new features:

1. Create a feature branch: `git checkout -b feature/your-feature`
2. Implement with tests
3. Ensure linters pass: `make lint`
4. Create a Pull Request with clear description
5. Use meaningful commit messages following conventional commits

---

## 📄 License

This project is part of the INSO2 course at Universidad de León.

---

## 👥 Team

**Backend Development:**

- Jorge (Go, Gin, GORM, Database Design, Tests)
- Ovidium

**Frontend Development:**

- Rodrigo

---

## 📞 Support

For issues or questions about the backend:

1. Check [docs/golang.md](docs/golang.md) for Go best practices
2. Review existing tests in `backend/internal/service/*_test.go`
3. Check API examples above
4. Open an issue on GitHub
