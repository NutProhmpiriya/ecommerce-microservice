.PHONY: build test run clean

# Build the application
build:
	go build -o bin/order-service backend/order-service/main.go

# Run tests
test:
	cd backend/order-service && go test ./... -v

# Run the application
run:
	cd backend/order-service && go run main.go

# Clean build artifacts
clean:
	rm -rf bin/
