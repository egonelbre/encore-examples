.PHONY: dev run test setup teardown

setup:
	docker compose up -d

teardown:
	docker compose down

run:
	go run .

test:
	go test ./...
