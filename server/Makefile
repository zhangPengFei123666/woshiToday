.PHONY: all build run test clean tidy fmt lint swagger

# 变量定义
APP_NAME=distributed-scheduler
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date +%Y-%m-%d_%H:%M:%S)
GO_VERSION=$(shell go version | cut -d ' ' -f 3)
LDFLAGS=-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GoVersion=${GO_VERSION}

# 默认目标
all: tidy fmt build

# 构建
build:
	@echo "Building..."
	@go build -ldflags "${LDFLAGS}" -o bin/admin ./cmd/admin

# 运行
run:
	@go run ./cmd/admin -config config.yaml

# 开发模式运行(热重载)
dev:
	@which air > /dev/null || go install github.com/cosmtrek/air@latest
	@air -c .air.toml

# 测试
test:
	@echo "Running tests..."
	@go test -v -cover ./...

# 测试覆盖率
cover:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# 清理
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html

# 整理依赖
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy

# 下载依赖
download:
	@echo "Downloading dependencies..."
	@go mod download

# 格式化代码
fmt:
	@echo "Formatting code..."
	@gofmt -s -w .

# 代码检查
lint:
	@which golangci-lint > /dev/null || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@golangci-lint run ./...

# 生成Swagger文档
swagger:
	@which swag > /dev/null || go install github.com/swaggo/swag/cmd/swag@latest
	@swag init -g cmd/admin/main.go -o docs

# 生成Mock
mock:
	@which mockgen > /dev/null || go install github.com/golang/mock/mockgen@latest
	@go generate ./...

# Docker构建
docker-build:
	@docker build -t ${APP_NAME}:${VERSION} .

# Docker运行
docker-run:
	@docker run -d --name ${APP_NAME} -p 8080:8080 ${APP_NAME}:${VERSION}

# 帮助
help:
	@echo "Available commands:"
	@echo "  make build    - Build the application"
	@echo "  make run      - Run the application"
	@echo "  make dev      - Run with hot reload"
	@echo "  make test     - Run tests"
	@echo "  make cover    - Generate test coverage report"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make tidy     - Tidy go modules"
	@echo "  make fmt      - Format code"
	@echo "  make lint     - Run linter"
	@echo "  make swagger  - Generate Swagger docs"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-run   - Run Docker container"

