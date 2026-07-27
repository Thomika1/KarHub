# KarHub - Restock Prioritization Engine

Go microservice for restock prioritization of auto parts. Calculates urgency scores to prioritize which products need restocking first.

## Tech Stack

- Go 1.26.5
- chi router
- go-kit endpoints
- GORM + pgx (PostgreSQL)
- slog (structured logging)
- go-playground/validator

## Prerequisites

- Go 1.26+
- PostgreSQL

## Setup

### 1. Clone and install dependencies

```bash
git clone https://github.com/Thomika1/KarHub.git
cd KarHub
go mod download
```

### 2. Environment variables

Create a `.env` file:

```env
PORT=8080
DB_ENGINE=postgres
DB_DSN="host=localhost port=5432 user=postgres password=postgres dbname=karhub sslmode=disable"
```

### 3. Database

Create the database:

```bash
psql -U postgres -c "CREATE DATABASE karhub-db;"
```

### 4. Run

```bash
go run cmd/api/main.go
```

Server starts at `http://localhost:8080`

## API Endpoints

Base path: `/api/v1`

### Products

#### Create Product

```bash
curl -X POST http://localhost:8080/api/v1/product \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Filtro de Óleo",
    "category": "filters",
    "currentStock": 15,
    "minimumStock": 20,
    "averageDailySales": 3,
    "leadTimeDays": 7,
    "unitCost": 1850,
    "criticalityLevel": 3
  }'
```

Response:
```json
{
  "id": "abc123...",
  "name": "Filtro de Óleo",
  "category": "filters",
  "currentStock": 15,
  "minimumStock": 20,
  "averageDailySales": 3,
  "leadTimeDays": 7,
  "unitCost": 1850,
  "criticalityLevel": 3,
  "createdAt": "2026-07-27T10:00:00Z",
  "updatedAt": "2026-07-27T10:00:00Z"
}
```

#### Get Product by ID

```bash
curl http://localhost:8080/api/v1/product/{id}
```

#### List Products (with filters, pagination, ordering)

```bash
curl -X POST http://localhost:8080/api/v1/product/list \
  -H "Content-Type: application/json" \
  -d '{
    "page": 0,
    "pageSize": 10,
    "orderBy": "name",
    "orderDirection": "asc",
    "filters": [
      {
        "field": "name",
        "value": "Filtro",
        "operator": "LIKE"
      }
    ]
  }'
```

#### Update Product (partial)

```bash
curl -X PUT http://localhost:8080/api/v1/product/{id} \
  -H "Content-Type: application/json" \
  -d '{
    "currentStock": 25,
    "minimumStock": 30
  }'
```

Only provided fields are updated.

#### Delete Product

```bash
curl -X DELETE http://localhost:8080/api/v1/product/{id}
```

### Restock Priorities

#### Get Priorities

Calculates urgency score: `(minimum_stock - projected_stock) * criticality_level`
where `projected_stock = current_stock - (average_daily_sales * lead_time_days)`

```bash
curl -X POST http://localhost:8080/api/v1/restock/priorities \
  -H "Content-Type: application/json" \
  -d '{
    "page": 0,
    "pageSize": 10
  }'
```

Response:
```json
{
  "priorities": [
    {
      "partId": "abc123...",
      "name": "Filtro de Óleo",
      "currentStock": 15,
      "projectedStock": -6,
      "minimumStock": 20,
      "urgencyScore": 78
    },
    {
      "partId": "def456...",
      "name": "Pastilha de Freio",
      "currentStock": 8,
      "projectedStock": -12,
      "minimumStock": 10,
      "urgencyScore": 110
    }
  ]
}
```

## Tests

To run tests:

```bash
go test ./...
```

Results are ordered by `urgencyScore DESC`, `criticalityLevel DESC`, `averageDailySales DESC`, `name ASC`.

## Categories

Valid values for `category` field:
- `engine`
- `brakes`
- `suspension`
- `transmission`
- `electrical`
- `cooling`
- `filters`
- `injection`
- `lighting`

## Validation Rules

| Field | Rule |
|---|---|
| `name` | required |
| `category` | required, must be valid enum |
| `currentStock` | >= 0 |
| `minimumStock` | >= 0 |
| `averageDailySales` | >= 0 |
| `leadTimeDays` | >= 0 |
| `unitCost` | >= 0 |
| `criticalityLevel` | 1-5 |

## Error Responses

| Status | Type | Description |
|---|---|---|
| 400 | ValidationError | Invalid field value |
| 409 | ConflictError | Duplicate product (same name + category) |
| 500 | Internal | Server error |
