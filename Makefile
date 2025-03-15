help:
	@echo "Available commands:"
	@echo "  go-install-air           Install Air for live reloading"
	@echo "  go-install-templ         Install Templ for template generation"
	@echo "  build                    Build the application"

.PHONY: go-install-air
go-install-air:
	go install github.com/air-verse/air@latest

.PHONY: go-install-templ
go-install-templ:
	go install github.com/a-h/templ/cmd/templ@latest

.PHONY: build
build:
	templ generate
	go build -o ./bin/http-server ./cmd/http-server/server.go
	go build -o ./bin/worker ./cmd/worker-server/main.go
	go build -o ./bin/scheduler ./cmd/scheduler/main.go
	go build -o ./bin/cli ./cmd/cli/main.go
