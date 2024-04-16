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

run-orders:
	go run ./cmd/local/order/
