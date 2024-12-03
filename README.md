# E-Commerce Microservices

This repository contains a microservices-based e-commerce application built with Go.

## Services

### Order Service

The Order Service handles all order-related operations in the e-commerce system. It is built using Clean Architecture principles and provides RESTful APIs for order management.

#### Features
- Create new orders
- Retrieve order details
- List all orders
- Update order status
- Delete orders

#### Tech Stack
- Go
- MongoDB
- Gorilla Mux (HTTP Router)
- Clean Architecture

#### API Endpoints
- `POST /api/v1/orders` - Create a new order
- `GET /api/v1/orders` - Get all orders
- `GET /api/v1/orders/{id}` - Get order by ID
- `PUT /api/v1/orders/{id}` - Update order
- `DELETE /api/v1/orders/{id}` - Delete order
- `GET /health` - Health check endpoint

#### Project Structure
```
backend/
└── order-service/
    ├── cmd/
    │   └── main.go
    └── internal/
        ├── domain/
        │   └── order.go
        ├── delivery/
        │   └── http/
        │       ├── order_handler.go
        │       └── order_handler_test.go
        ├── repository/
        │   ├── mongo/
        │   │   └── order_repository.go
        │   └── mock/
        │       └── order_repository_mock.go
        └── usecase/
            ├── order_usecase.go
            └── order_usecase_test.go
```

#### Getting Started

1. Prerequisites
   - Go 1.16 or later
   - MongoDB
   - Make (optional)

2. Environment Variables
   ```
   MONGODB_URI=mongodb://localhost:27017
   PORT=8080
   ```

3. Running the Service
   ```bash
   # From the order-service directory
   go run main.go
   ```

4. Running Tests
   ```bash
   go test ./... -v
   ```

## Contributing
1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License
This project is licensed under the MIT License - see the LICENSE file for details.
