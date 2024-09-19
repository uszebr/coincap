TEMPLE_CMD := templ generate
TEST_CMD := go test ./...
TIDY_CMD := go mod tidy

# Define variables (adjust as needed)
GOARCH := amd64
GOOS := linux  # You can add other OS options here (e.g., windows, darwin)
BINARY_NAME := coincap

# Generate target
generate:
	$(TEMPLE_CMD)

build: generate
	go build -o $(BINARY_NAME) ./cmd/coincap.go

clean:
	rm -f $(BINARY_NAME)

run: build  # Specify build as a dependency here
	./$(BINARY_NAME)

tidy: 
	$(TIDY_CMD)

test: 
	$(TEST_CMD)
