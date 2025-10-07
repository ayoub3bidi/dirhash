#  ____  _      _               _     
# |  _ \(_)_ __| |__   __ _ ___| |__  
# | | | | | '__| '_ \ / _` / __| '_ \ 
# | |_| | | |  | | | | (_| \__ \ | | |
# |____/|_|_|  |_| |_|\__,_|___/_| |_|

# Built @Think-it 2022

build-cli:
	go build -o bin/dirhash

.PHONY: deps tidy fmt vet test lint build clean

deps:
	go mod download

tidy:
	go mod tidy

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./...

# Lint uses golangci-lint if available; otherwise no-op with message
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed; skipping lint"; \
	fi

build:
	go build ./...

clean:
	rm -rf bin/
