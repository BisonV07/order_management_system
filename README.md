# Order Management System (OMS)

A high-integrity Order Management System built with React + Vite frontend and Go backend, capable of handling 100,000 daily orders with 1,000 concurrent submissions.

## Architecture

- **Client**: React + Vite SPA
- **Server**: Go service with three-layer architecture (API → Service → DataStore)
- **Database**: PostgreSQL with ACID compliance
- **State Management**: Finite State Machine (FSM) for order lifecycle
- **Concurrency Control**: Pessimistic locking (SELECT FOR UPDATE) for inventory

## Features

- Zero overselling with strict inventory consistency
- Order state machine (ORDERED → SHIPPED → DELIVERED, CANCELLED)
- JWT-based authentication
- Audit logging for order status changes
- Rate limiting to prevent spam

## Project Structure

```
order_management_system/
├── client/            # React + Vite SPA
├── server/            # Go service
├── migrations/        # Database migrations
└── docker-compose.yml # Local development setup
```

## Getting Started

See [RUN.md](RUN.md) for detailed setup and run instructions.

### Quick Start

```bash
# 1. Start database
docker compose up -d postgres
# OR if that doesn't work: docker-compose up -d postgres

# 2. Server (Terminal 1)
cd server
go mod download
cp .env.sample .env
go run cmd/main.go --migrate
go run cmd/main.go --api --port=8080

# 3. Client (Terminal 2)
cd client
npm install
cp .env.sample .env
npm run dev
```

The client will be available at `http://localhost:3000` and the server at `http://localhost:8080`.

## API Endpoints

### Create Order
- **POST** `/api/v1/orders`
- **Body**: `{ "product_id": 101, "quantity": 2 }`
- **Response**: `{ "order_id": "...", "current_status": "ORDERED", "message": "Order placed successfully" }`

### Update Order Status
- **PATCH** `/api/v1/orders/{orderId}`
- **Body**: `{ "current_status": "SHIPPED" }`
- **Response**: `{ "order_id": "...", "previous_status": "ORDERED", "current_status": "SHIPPED", ... }`

## Database Schema

- **products**: Product catalog with SKU, name, price, metadata
- **inventory**: Stock quantities per product
- **orders**: Order records with status tracking
- **order_state_logs**: Audit trail of status changes

## Development

### Running Tests
```bash
cd server
make test
```

### Building
```bash
cd server
make build
```

## License

MIT

