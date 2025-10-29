.PHONY: build run clean wasm server air

# Build all components
build: wasm air

# Build WASM frontend
wasm:
	GOOS=js GOARCH=wasm go build -o web/main.wasm ./app

# Build server
server:
	go build -o server ./server

air:
	go build -o tmp/main ./server

# Run server locally
run:
	go run ./server

# Run WASM build and serve
run-wasm: wasm
	go run ./server

# Clean build artifacts
clean:
	rm -f web/main.wasm server

# Add these to your Makefile
migrate-up:
	migrate -path server/migrations -database $(DATABASE_URL) up

migrate-down:
	migrate -path server/migrations -database $(DATABASE_URL) down

migrate-force:
	migrate -path server/migrations -database $(DATABASE_URL) force $(VERSION)
	