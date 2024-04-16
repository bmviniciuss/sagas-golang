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
	go run ./cmd/local/

run-orders:
	go run ./cmd/order/
