.PHONY: run

run: run
	$(eval include ./cmd/.env.local)
	$(eval export)
	@go mod tidy
	@go run cmd/main.go

