.PHONY: all build docker clean gprc $(MODULES)

BUILD_DIR = build
SRC_FILES = $(shell find . -name 'main.go')

# 动态获取 build 目录下所有的文件名作为模块名称
MODULES := $(shell ls build)

# 定义每个模块的构建规则
$(MODULES):
	docker build -t dahuang/$@:latest --build-arg MODULE_NAME=$@ .


# 默认目标：构建所有模块
all: $(MODULES)



down:
	@kubectl delete -f ./k8s-webook-service.yaml || true
	@kubectl delete -f ./k8s-webook-redis.yaml || true
	@kubectl delete -f ./k8s-webook-mysql.yaml || true
	@kubectl delete -f ./k8s-webook-ingress.yaml || true
docker:
	@rm webook || true
	@go mod tidy
	@GOOS=linux GOARCH=arm go build -o webook .
	@docker rmi -f dahuang/webook:v0.0.1
	@docker build -t dahuang/webook:v0.0.1 .
k8s:
	@kubectl delete -f ./k8s-webook-service.yaml || true
	@kubectl apply -f ./k8s-webook-service.yaml
redis:
	@kubectl delete -f ./k8s-webook-redis.yaml || true
	@kubectl apply -f ./k8s-webook-redis.yaml
mysql:
	@kubectl delete -f ./k8s-webook-mysql.yaml || true
	@kubectl apply -f ./k8s-webook-mysql.yaml

grpc:
	@buf generate api/proto

# Build all Go applications with main.go
build:
	@mkdir -p $(BUILD_DIR)
	@for file in $(SRC_FILES); do \
		dir=$$(dirname $$file); \
		dir_name=$$(basename $$dir); \
		echo "Building $$file in directory $$dir"; \
		GOOS=linux GOARCH=arm go build -o $(BUILD_DIR)/$$dir_name $$dir; \
	done

# Clean build artifacts
clean:
	@for module in $(MODULES); do \
		echo "Removing Docker image for: $$module"; \
		docker rmi dahuang/$$module:latest; \
	done

update:
	export ETCDCTL_ENDPOINTS="http://127.0.0.1:12379"
	cat ./config/dev.yaml | etcdctl put /config/config.yaml

cleanDir:
	rm -rf $(BUILD_DIR)
