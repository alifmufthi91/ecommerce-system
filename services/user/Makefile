start:
	`go env GOPATH`/bin/swag init --parseDependency --parseInternal &&`go env GOPATH`/bin/air

swaggo:
	@`go env GOPATH`/bin/swag init --parseDependency --parseInternal

test:
	@go test -v -tags dynamic `go list ./... | grep -i 'internal\|pkg'` -cover -coverprofile=coverage.out
