lint:
	@echo ">> Linting with golangci-lint"
	@golangci-lint run

test:
	@echo ">> Running tests"
	@go clean -testcache
	@TZ="" go test -count=1 ./... -tags=faketime

