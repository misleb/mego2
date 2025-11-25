.PHONY: build run clean wasm server air cmd

# Build all components
build: wasm air cmd

cmd:
	go build -o bin/cmd ./cmd

# Build WASM frontend
wasm:
	GOOS=js GOARCH=wasm go build -o server/web/main.wasm ./app

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

migrate-up:
	@DATABASE_URL=$$(heroku config:get DATABASE_URL); \
	migrate -path server/migrations -database $$DATABASE_URL up

migrate-test-up:
	migrate -path server/migrations -database postgresql://mego2:mego2_dev@localhost:5432/mego2_test?sslmode=disable up

migrate-down:
	@DATABASE_URL=$$(heroku config:get DATABASE_URL); \
	migrate -path server/migrations -database $$DATABASE_URL down

migrate-force:
	@DATABASE_URL=$$(heroku config:get DATABASE_URL); \
	migrate -path server/migrations -database $$DATABASE_URL force $(VERSION)