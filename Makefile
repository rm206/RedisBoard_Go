# Makefile

# Variables
GO := go
SRC := $(wildcard *.go)
BIN := main

# Default target
.PHONY: all
all: $(BIN)

# Build target
$(BIN): $(SRC)
	$(GO) build -o $(BIN)

# Run target
.PHONY: run
run: $(BIN)
	./$(BIN)

# Clean target
.PHONY: clean
clean:
	rm -f $(BIN)