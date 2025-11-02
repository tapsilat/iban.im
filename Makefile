.PHONY: build build-frontend build-backend clean install-frontend test

# Build everything (frontend + backend)
build: build-frontend build-backend

# Install frontend dependencies
install-frontend:
	cd web && npm install

# Build frontend and copy to static/
build-frontend: install-frontend
	cd web && npm run build
	mkdir -p static/dist
	cp -r web/dist/* static/dist/

# Build Go binary
build-backend:
	go build -o iban.im

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf web/dist
	rm -rf static/dist
	rm -f iban.im

# Run the application
run: build
	./iban.im
