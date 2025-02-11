help:
	@echo "Available commands:"
	@echo "  go-install-air         Install Air for live reloading"
	@echo "  install-tailwind       Install TailwindCSS"
	@echo "  tailwind-watch         Watch and compile TailwindCSS files"
	@echo "  tailwind-build         Build TailwindCSS files"
	@echo "  build-http-server      Build the HTTP server"

.PHONY: go-install-air
go-install-air:
	go install github.com/air-verse/air@latest

.PHONY: install-tailwind
install-tailwind:
	@if [ "$$(uname)" = "Darwin" ]; then \
		if [ "$$(uname -m)" = "arm64" ]; then \
			curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-macos-arm64; \
			chmod +x tailwindcss-macos-arm64; \
			mv tailwindcss-macos-arm64 tailwindcss; \
		else \
			curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-macos-x64; \
			chmod +x tailwindcss-macos-x64; \
			mv tailwindcss-macos-x64 tailwindcss; \
		fi \
	elif [ "$$(uname)" = "Linux" ]; then \
		curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64; \
		chmod +x tailwindcss-linux-x64; \
		mv tailwindcss-linux-x64 tailwindcss; \
	else \
		echo "Unsupported operating system"; \
		exit 1; \
	fi

.PHONY: tailwind-watch
tailwind-watch:
	./tailwindcss -i ./static/css/custom.css -o ./static/css/style.css --watch --config ./tailwind.config.js

.PHONY: tailwind-build
tailwind-build:
	./tailwindcss -i ./static/css/custom.css -o ./static/css/style.css --config ./tailwind.config.js

.PHONY: build-http-server
build-http-server:
	./tailwindcss -i ./static/css/custom.css -o ./static/css/style.css --config ./tailwind.config.js
	templ generate
	go build -o ./bin/http-server ./cmd/http-server/server.go

