.PHONY: proto migrate-up migrate-down test run-auth run-forum swag

# Proto generation
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/*.proto

# Database migrations
migrate-up:
	migrate -path migrations/auth -database "postgresql://localhost:5432/forum_auth?sslmode=disable" up
	migrate -path migrations/forum -database "postgresql://localhost:5432/forum_main?sslmode=disable" up

migrate-down:
	migrate -path migrations/auth -database "postgresql://localhost:5432/forum_auth?sslmode=disable" down
	migrate -path migrations/forum -database "postgresql://localhost:5432/forum_main?sslmode=disable" down

# Testing
test:
	go test -v -cover ./...

# Run services
run-auth:
	go run cmd/auth-service/main.go

run-forum:
	go run cmd/forum-service/main.go

# Swagger documentation
swag:
	swag init -g cmd/auth-service/main.go -o docs/auth
	swag init -g cmd/forum-service/main.go -o docs/forum 