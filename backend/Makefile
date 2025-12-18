.PHONY: dev build run swagger swagger-install tidy clean

# 开发模式运行
dev:
	go run ./cmd/server

# 构建二进制
build:
	go build -o bin/server ./cmd/server

# 运行已构建的二进制
run: build
	./bin/server

# 生成 Swagger 文档
swagger:
	swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal

# 安装 swag 工具
swagger-install:
	go install github.com/swaggo/swag/cmd/swag@latest

# 整理依赖
tidy:
	go mod tidy

# 清理构建产物
clean:
	rm -rf bin/
	rm -rf docs/docs.go docs/swagger.json docs/swagger.yaml
