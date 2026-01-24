VERSION ?= dev

.PHONY: build build-go test test-go test-go-integration test-unit release clean

# Build Go version (default)
build: build-go

build-go:
	CGO_ENABLED=0 go build -ldflags "-X 'main.version=$(VERSION)'" -o dist/cider ./cmd/cider

# Run Go tests
test-go:
	CGO_ENABLED=0 go test ./internal/...

# Run Go tests including integration (requires Apple Notes)
test-go-integration:
	CGO_ENABLED=0 go test -tags=integration ./internal/...

test-unit:
	SKIP_INTEGRATION=1 test/approve

# Clean build artifacts
clean:
	rm -f dist/cider-go
	rm -f cider-go
