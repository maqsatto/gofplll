.PHONY: test test-integration smoke fmt vet

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./...

test-integration:
	go test ./... -run Integration

smoke:
	go run ./cmd/gofplll-smoke
