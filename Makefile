build:
	go build ./cmd

run:
	go run ./cmd/main.go

test:
	go test -v ./...

migrate:
  go run ./migrate/migrate.go
