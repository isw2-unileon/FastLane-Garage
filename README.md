# FastLane Garage - Full Stack Auto Parts Management

A complete full-stack application for managing automotive parts and customer orders, built with **Go** backend and **React** frontend in a monorepo structure.

## 🎯 Features

### Backend (Go + Gin + GORM)

✅ **Parts & Order Management**

- CRUD operations for car parts with filtering (by car zone) and full-text search.
- Complete order lifecycle (pending -> processing -> completed).
- Order items with quantity, unit pricing, and automatic total price calculation.
- Customer email validation and order status tracking.

✅ **AI Chatbot Integration**

- Automated technical assistant powered by n8n webhooks.
- Session management (creation and history retrieval).
- Persistent chat message storage linking users to the AI's responses.

✅ **System Analytics**

- Comprehensive stats tracking (revenue, total orders, parts count, etc.).

✅ **Architecture & Quality**

- Clean architecture (Handlers -> Service -> Repository).
- Dependency injection and interface-based design for testability.
- Comprehensive unit tests with mocks (100% test coverage for the service layer).
- Go linters (golangci-lint), structured logging with `slog`, and type-safe DTOs.
- SQLite database configured via GORM.

### Frontend (React + Three + TypeScript + Vite)

✅ **Core Architecture (React & TypeScript)**

- **React**: Main UI library handling modular component creation, global state management (e.g., cart items, active views), and efficient interface rendering.
- **TypeScript**: Strict static typing across the client. Defines precise interfaces (like `HeaderProps` or `ChatViewProps`) to guarantee that frontend data payloads sync perfectly with backend expectations.

✅ **Interactive 3D Digital Twin (Three.js)**

- Renders an immersive 3D digital twin of the vehicle (`coche.glb`) and environment.
- Features fully interactive camera controls (orbit, rotate, zoom).
- **Dynamic Interaction**: Users can click directly on car meshes to dynamically extract specific car zone tags for part filtering.

✅ **Modern UI & Aesthetics (Tailwind CSS & Lucide)**

- **Tailwind CSS**: Drives the responsive, dark industrial aesthetic. Controls scroll flows, Grid/Flexbox alignments, and neon glow/animation effects for the environment.
- **Lucide React**: Supplies clean vector iconography throughout the interface (e.g., `ShoppingCart`, `BarChart3`, `Bot`, `Send`).

✅ **Build & Dev Environment (Vite)**

- Ultra-fast bundler and development server.
- Automatically manages the reverse proxy, seamlessly redirecting `/api` requests from the client port (`5173`) to the Go backend server (`8080`).

### AI & Automation Workflow (Ovidium's Engine)

✅ **Intelligent Chatbot Integration**

- **n8n (Workflow Orchestration):** Acts as the central nervous system. It receives webhook triggers from the Go backend containing the user's message and vehicle context, orchestrating the entire logical flow.
- **SerpApi (Real-Time Search):** Before generating an answer, n8n queries SerpApi to fetch real-time data from the web. This allows the bot to dynamically check up-to-date part compatibilities, technical manuals, and vehicle specifications.
- **DeepSeek API (LLM Engine):** The core intelligence of the assistant. It ingests the user's prompt alongside the real-time context gathered by SerpApi to generate accurate, technical, and context-aware responses.
- **Persistent Memory:** The Go backend stores all sessions and messages in the SQLite database, creating a seamless history log.

---

## 🏗️ Project Structure

```text
├── backend/                    # Go API Server
│   ├── cmd/server/             # Entry point (main.go)
│   └── internal/
│       ├── config/             # CORS and Environment configuration
│       ├── database/           # Database migrations (SQLite) & seeding
│       ├── models/             # Domain models (Part, Order, ChatSession, ChatMessage)
│       ├── dto/                # Data transfer objects with validation
│       ├── repository/         # Data access layer
│       ├── service/            # Business logic (Parts, Orders, Chat, Stats)
│       └── handlers/           # Gin HTTP handlers
│
├── frontend/                   # React + TypeScript + Vite
│   ├── src/
│   │   ├── assets/models/      # 3D .glb models (car, scenario, map)
│   │   ├── components/layout/  # Global layout (Header, etc.)
│   │   └── features/
│   │       ├── garage/         # 3D Scene components (Scene3D, CarModel)
│   │       └── parts/          # Dashboard, Analytics, ChatView, API services
│   └── index.html
│
├── e2e/                        # Playwright E2E tests
├── docs/                       # Documentation & Architecture Decision Records
└── .github/workflows/          # CI/CD pipelines (Backend, Frontend, CodeQL, E2E)
```

## 📊 API Endpoint

### Parts & Orders Endpoints

| Method   | Path                     | Description         | Query Parameters                         |
| -------- | ------------------------ | ------------------- | ---------------------------------------- |
| `GET`    | `/api/parts`             | List all parts      | `?zone=motor&name=Motor`                 |
| `GET`    | `/api/parts/:id`         | Get part by ID      | -                                        |
| `POST`   | `/api/parts`             | Create new part     | -                                        |
| `PUT`    | `/api/parts/:id`         | Update part         | -                                        |
| `DELETE` | `/api/parts/:id`         | Delete part         | -                                        |
| `GET`    | `/api/orders`            | List all orders     | `?status=pending&email=user@example.com` |
| `GET`    | `/api/orders/:id`        | Get order by ID     | -                                        |
| `POST`   | `/api/orders`            | Create new order    | -                                        |
| `PUT`    | `/api/orders/:id/status` | Update order status | -                                        |
| `DELETE` | `/api/orders/:id`        | Delete order        | -                                        |

### AI Chat Endpoints

| Method | Path                              | Description                  | Query Parameters |
| ------ | --------------------------------- | ---------------------------- | ---------------- |
| `POST` | `/api/chat/sessions`              | Create a new AI chat session | -                |
| `GET`  | `/api/chat/sessions/:id`          | Get chat session history     | -                |
| `POST` | `/api/chat/sessions/:id/messages` | Send message to the AI bot   | -                |

**Parts examples:**

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

**Orders examples:**

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

**Chat message example:**

```bash
// Sending a chat message
curl -X POST http://localhost:8080/api/chat/sessions/1/messages \
  -H "Content-Type: application/json" \
  -d '{
    "text": "¿Son compatibles las pastillas de freno con un Audi A4?"
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
# Run backend tests
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

### Chat Integration Table

```sql
CREATE TABLE chat_sessions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  status TEXT,
  vehicle_brand TEXT,
  vehicle_model TEXT,
  vehicle_year TEXT,
  parts TEXT,
  created_at DATETIME,
  updated_at DATETIME
);

CREATE TABLE chat_messages (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  session_id INTEGER,
  role TEXT,
  content TEXT,
  created_at DATETIME,
  FOREIGN KEY (session_id) REFERENCES chat_sessions(id)
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

```
[Frontend React App]
        ↓↑
[Go HTTP Handlers]    → Validates request & routes to Service
        ↓↑
[Go Service Layer]    → Calls Database Repo AND triggers AI Webhook
        ↓↑
  +-----------+
  |  n8n Node |       ← Orchestrates the AI logical flow
  +-----------+
    ↙       ↘
[SerpApi]  [DeepSeek] ← SerpApi fetches web context, DeepSeek generates the response
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
- Ovidium (AI Workflow Integration; n8n orchestration, SerpApi real-time search, DeepSeek LLM automation)

**Frontend Development:**

- Rodrigo (React, Three.js 3D Garage, ChatUI, Tailwind)

---

## 📞 Support

For issues or questions about the backend:

1. Check [docs/golang.md](docs/golang.md) for Go best practices
2. Review existing tests in `backend/internal/service/*_test.go`
3. Check API examples above
4. Open an issue on GitHub
