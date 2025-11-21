.PHONY: build serve agent clean all

# 编译所有
all: serve agent

# 编译 serve
serve:
	@echo "Building serve..."
	@go build -o bin/serve ./cmd/serve

# 编译 agent
agent:
	@echo "Building agent..."
	@go build -o bin/agent ./cmd/agent

# 清理编译产物
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/

# 编译并运行 serve
run-serve: serve
	@echo "Running serve..."
	@./bin/serve

# 编译并运行 agent
run-agent: agent
	@echo "Running agent..."
	@./bin/agent

