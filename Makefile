PROJECT_NAME := "Bastard GIT"
BINARY_NAME := "bgit"

.PHONY: all clean test test-verbose run build help init

all: clean build test
prepare: clean build init

clean:
	@echo "Cleaning"
	rm -f  $(BINARY_NAME) || true
	rm -rf srctest/.bgit || true

test:
	@echo "Testing"
	go test -v ./...

build:
	@echo "Building $(PROJECT_NAME)"
	CGO_ENABLED=0 go build -o $(BINARY_NAME) main.go
	chmod +x $(BINARY_NAME)

release:
	@echo "Pushing a new tag"
	git fetch --tags --force
	latest_tag=$(git describe --tags `git rev-list --tags --max-count=1` || echo "v0.0.0")
	latest_version=$(echo $latest_tag | sed 's/^v//')
	new_version=$(echo $latest_version | awk -F. -v OFS=. '{$NF++;print}')
	@echo "New version: $new_version"
	git tag v$new_version
	git push origin v$new_version

run:
	@echo "Running $(PROJECT_NAME)"
	go run main.go

init:
	@echo "Initializing a new bgit repo..."
	"./$(BINARY_NAME)" init
