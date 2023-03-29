.PHONY: run

run: run
	$(eval include ./cmd/.env.local)
	$(eval export)
	@go run cmd/main.go

