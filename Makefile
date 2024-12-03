.PHONY: all build docker clean gprc $(MODULES)


WIRE := $(shell find . -name 'wire.go')

# 定义每个模块的构建规则
$(WIRE):
	wire

all: $(MODULES)

docker:
	@rm webook || true
	@go mod tidy
	@GOOS=linux GOARCH=arm go build -o webook .
	@docker rmi -f dahuang/webook:v0.0.1
	@docker build -t dahuang/webook:v0.0.1 .

.PHONY: fmt
fmt:
	@sh ./.script/fmt.sh
