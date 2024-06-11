tidy:
	go fmt ./...
	go mod tidy -v

test:
	go test -v ./...

test/cover:
	go test -v -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

dep:
	go mod download

run:
	go run ./cmd/local/orchestrator/

orders:
	go run ./cmd/local/order/

customers:
	go run ./cmd/local/customer/

kitchen:
	go run ./cmd/local/kitchen/

accounting:
	go run ./cmd/local/accounting/

services-up:
	docker compose up orchestrator orders customers accounting -d

services-down:
	docker compose down orchestrator orders customers accounting

services-logs:
	docker compose logs orchestrator orders customers accounting -f
