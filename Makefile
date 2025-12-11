.PHONY: build serve agent clean all dashboard-dist print-vars

VERSION := $(shell git describe --tags --always --dirty)
DEBUG := $(shell if [ -n "$$DEBUG" ]; then echo "true"; else echo "false"; fi)

ifeq ($(DEBUG), true)
	LDFLAGS := -X 'github.com/benlocal/lai-panel/pkg/version.Version=$(VERSION)'
else
	LDFLAGS := -s -w -X 'github.com/benlocal/lai-panel/pkg/version.Version=$(VERSION)'
endif

print-vars:
	@echo "VERSION: $(VERSION)"
	@echo "DEBUG: $(DEBUG)"
	@echo "LDFLAGS: $(LDFLAGS)"

all: print-vars serve agent

dashboard-dist:
	@echo "Ensuring dashboard/dist exists..."
	@mkdir -p dashboard/dist
	@if [ -z "$$(ls -A dashboard/dist 2>/dev/null)" ]; then \
		touch dashboard/dist/.keep; \
	fi

serve: dashboard-dist
	@echo "Building serve..."
	@go build -ldflags "$(LDFLAGS)" -o bin/serve ./cmd/serve

agent:
	@echo "Building agent..."
	@go build -ldflags "$(LDFLAGS)" -o bin/agent ./cmd/agent

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/

run-serve: serve
	@echo "Running serve..."
	@./bin/serve

run-agent: agent
	@echo "Running agent..."
	@./bin/agent

test:
	@echo "Running tests..."
	@go test -v ./...

