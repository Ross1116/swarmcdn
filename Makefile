# Makefile for SwarmCDN

# Output directory
BIN_DIR := bin

# Binary names
MAIN_SERVER := $(BIN_DIR)/central-server
PEER_CLIENT := $(BIN_DIR)/peer-client
PEER_SERVER := $(BIN_DIR)/peer-server

# Build all binaries
.PHONY: all
all: $(MAIN_SERVER) $(PEER_CLIENT) $(PEER_SERVER)

# Build main server
$(MAIN_SERVER): main.go
	@echo "Building main server..."
	@mkdir -p $(BIN_DIR)
	go build -o $(MAIN_SERVER) .

# Build peer client
$(PEER_CLIENT):
	@echo "Building peer client..."
	@mkdir -p $(BIN_DIR)
	go build -o $(PEER_CLIENT) ./peer/client

# Build peer server
$(PEER_SERVER):
	@echo "Building peer server..."
	@mkdir -p $(BIN_DIR)
	go build -o $(PEER_SERVER) ./peer/server

# Individual build targets
.PHONY: build-main build-client build-server
build-main: $(MAIN_SERVER)

build-client: $(PEER_CLIENT)

build-server: $(PEER_SERVER)

# Run targets
.PHONY: run-main run-client run-server

run-main: build-main
	@echo "Running central server..."
	$(MAIN_SERVER)

run-client: build-client
	@echo "Running peer client..."
	$(PEER_CLIENT)

run-server: build-server
	@echo "Running peer server..."
	$(PEER_SERVER)

# Clean binaries
.PHONY: clean
clean:
	@echo "Cleaning binaries..."
	rm -rf $(BIN_DIR)

# Rebuild all
.PHONY: rebuild
rebuild: clean all

# Always build all
.PHONY: build-all build

build: build-all

build-all:
	@echo "Building all binaries (forced rebuild)..."
	@mkdir -p $(BIN_DIR)
	go build -o $(MAIN_SERVER) .
	go build -o $(PEER_CLIENT) ./peer/client
	go build -o $(PEER_SERVER) ./peer/server
