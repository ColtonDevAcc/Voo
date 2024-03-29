## test: runs all tests
test:
	@go test -v ./...

## cover: opens coverage in browser
cover:
	@go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

## coverage: displays test coverage
coverage:
	@go test -cover ./...

## build_cli: builds the command line tool voo and copies it to myapp
build_cli:
	@go build -o ../myapp/voo ./cmd/cli

#run docker in the background
start_compose:
	@docker-compose up -d

end_compose:
	@docker-compose down