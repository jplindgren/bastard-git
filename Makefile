PROJECT_NAME := "Bastard GIT"
BINARY_NAME := "bgit"

.PHONY: all clean test test-verbose run build help init

all: clean build test
prepare: clean build init

clean:
	@echo "Cleaning"
	rm -f  $(BINARY_NAME) || true
	rm -rf srctest/.bgit || true

test: go test -v ./...

build:
	@echo "Building $(PROJECT_NAME)"
	go build -o $(BINARY_NAME) main.go
	chmod +x $(BINARY_NAME)

build-mac:
	@echo "Building $(PROJECT_NAME) for macOS m1"
	GOOS=darwin GOARCH=amd64 go build -o "$(BINARY_NAME)-m1" main.go
	chmod +x "$(BINARY_NAME)-m1"

run:
	@echo "Running $(PROJECT_NAME)"
	go run main.go

init:
	@echo "Initializing a new bgit repo..."
	"./$(BINARY_NAME)" init
