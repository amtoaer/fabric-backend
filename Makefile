.PHONY: all dev clean build env-up env-down run

all: clean build env-up run

dev: build run

build:
	@echo "开始下载依赖..."
	@go mod download
	@echo "开始构建程序..."
	@go build
	@echo "构建成功。"

env-up:
	@echo "启动网络中..."
	@cd fixtures && docker-compose up --force-recreate -dev
	@echo "网络已启动。"

env-down:
	@echo "停止网络中..."
	@cd fixtures && docker-compose down
	@echo "网络已停止。"

run:
	@echo "启动程序中..."
	@./fabric-backend

clean: env-down
	@echo "清理环境中..."
	@rm -rf /tmp/fabric-backend-*
	@-docker rm -f -v `docker ps -a --no-trunc | grep "fabric-backend" | cut -d ' ' -f 1` 2>/dev/null
	@-docker rmi `docker images --no-trunc | grep "fabric-backend" | cut -d ' ' -f 1` 2>/dev/null
	@echo "清理环境完成"