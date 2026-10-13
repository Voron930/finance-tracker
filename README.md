# Finance Tracker API

REST API for personal finance tracking, built with Go and PostgreSQL. Users record income and expenses by category, see their balance before and after each operation, and get monthly statistics.

> **Status:** in active development. See [Roadmap](#roadmap) for what is done and what is next.

## Features

- User registration and login with bcrypt password hashing and JWT authentication
- Income and expense records with amount, category, description and date
- Default categories: Food, Rent, Transport, Sport, Education, plus custom categories
- Current balance, and the balance before and after every operation
- Transaction list with filters by date range, type and category
- Monthly statistics: total income, total expenses, balance and spending per category
- Every user sees only their own data (ownership checks on every request)

## Tech Stack

- **Go** 1.22+, standard `net/http` router
- **PostgreSQL** 16 with `pgx/v5` / `pgxpool`
- **JWT** (`golang-jwt/jwt/v5`) and **bcrypt** (`golang.org/x/crypto/bcrypt`)
- **Docker** and **Docker Compose**
- Postman for API testing

## Architecture

Layered architecture with a clear split of responsibilities:

```
Client (HTTP/JSON)
   ↓
Router (net/http)
   ↓
Middleware (JWT auth → user_id in context)
   ↓
Handler     – parses requests, validates input, writes responses
   ↓
Service     – business logic: hashing, JWT, validation, balance
   ↓
Repository  – SQL queries via pgxpool
   ↓
PostgreSQL
```

## Project Structure

```
cmd/api/              – entry point: config, DB, wiring, HTTP server
internal/app/         – routes and middleware registration
internal/config/      – loads environment variables
internal/db/          – PostgreSQL connection pool
internal/handler/     – HTTP handlers
internal/middleware/  – JWT authentication middleware
internal/model/       – domain models (User, Category, Transaction)
internal/repository/  – database layer
internal/service/     – business logic
migrations/           – SQL schema
```

## Database

Three tables: `users`, `categories`, `transactions`.

- `transactions.type` is either `income` or `expense`
- `transactions.amount` is `DECIMAL(12,2)` and must be greater than 0
- Deleting a user cascades to their categories and transactions
- A category cannot be deleted while transactions use it

Full schema: [`migrations/init.sql`](migrations/init.sql)

## API Endpoints

### Public

| Method | Endpoint         | Description                 |
| ------ | ---------------- | --------------------------- |
| POST   | `/auth/register` | Register a new user         |
| POST   | `/auth/login`    | Log in and receive a JWT    |

### Protected (header `Authorization: Bearer <token>`)

| Method | Endpoint                     | Description                                   |
| ------ | ---------------------------- | --------------------------------------------- |
| GET    | `/api/v1/categories`         | List the user's categories                    |
| POST   | `/api/v1/categories`         | Create a custom category                      |
| GET    | `/api/v1/transactions`       | List transactions (with filters, see below)   |
| POST   | `/api/v1/transactions`       | Add income or expense                         |
| GET    | `/api/v1/transactions/{id}`  | Get a transaction by ID                       |
| PUT    | `/api/v1/transactions/{id}`  | Update a transaction                          |
| DELETE | `/api/v1/transactions/{id}`  | Delete a transaction                          |
| GET    | `/api/v1/balance`            | Current balance                               |
| GET    | `/api/v1/stats`              | Statistics for a period                       |

### Filters

```
GET /api/v1/transactions?from=2026-10-01&to=2026-10-31&type=expense&category_id=1
```

### Example: add an expense

Request:

```json
POST /api/v1/transactions
{
  "category_id": 1,
  "type": "expense",
  "amount": 250.00,
  "description": "Groceries"
}
```

Response `201 Created`:

```json
{
  "id": 42,
  "category": "Food",
  "type": "expense",
  "amount": 250.00,
  "description": "Groceries",
  "balance_before": 5000.00,
  "balance_after": 4750.00,
  "created_at": "2026-10-12T18:30:00Z"
}
```

### Example: monthly statistics

```
GET /api/v1/stats?from=2026-10-01&to=2026-10-31
```

```json
{
  "period": { "from": "2026-10-01", "to": "2026-10-31" },
  "total_income": 30000.00,
  "total_expense": 18450.00,
  "balance": 11550.00,
  "expenses_by_category": [
    { "category": "Rent", "amount": 12000.00 },
    { "category": "Food", "amount": 4200.00 },
    { "category": "Transport", "amount": 1250.00 },
    { "category": "Sport", "amount": 1000.00 }
  ]
}
```

## Getting Started

1. Clone the repository:

   ```bash
   git clone https://github.com/Voron930/finance-tracker.git
   cd finance-tracker
   ```

2. Create a `.env` file from the example and adjust the values:

   ```bash
   cp .env.example .env
   ```

   ```
   APP_PORT=8080
   DATABASE_URL=postgres://postgres:secret@localhost:5432/finance_tracker?sslmode=disable
   JWT_SECRET=change_me
   ```

3. Start PostgreSQL with Docker:

   ```bash
   docker compose up -d
   ```

4. Apply the schema from `migrations/init.sql`, then run the server:

   ```bash
   go run ./cmd/api
   ```

The API is available at `http://localhost:8080`.

## Roadmap

- [x] Project structure and Go modules
- [x] PostgreSQL in Docker, database schema
- [x] Config loading and DB connection pool
- [x] Models: User, Category, Transaction
- [x] User repository
- [ ] Auth service: registration and login with bcrypt
- [ ] JWT generation and auth middleware
- [ ] Transactions: repository, service, handlers
- [ ] Balance before and after each operation
- [ ] Filters by date, type and category
- [ ] Monthly statistics
- [ ] Unit tests
- [ ] Run the whole app in Docker Compose

## Author

Bogdan Voronik — [GitHub](https://github.com/Voron930)
