.PHONY: docker
down:
	@kubectl delete -f ./k8s-webook-service.yaml || true
	@kubectl delete -f ./k8s-webook-redis.yaml || true
	@kubectl delete -f ./k8s-webook-mysql.yaml || true
	@kubectl delete -f ./k8s-webook-ingress.yaml || true
mock:
	@mockgen -source=internal/service/code.go -package=svcmocks -destination=internal/service/mocks/code.mock.go
	@mockgen -source=internal/service/user.go -package=svcmocks -destination=internal/service/mocks/user.mock.go
	@mockgen -source=internal/repository/user.go -package=repomocks -destination=internal/repository/mocks/user.mock.go
	@mockgen -source=internal/repository/code.go -package=repomocks -destination=internal/repository/mocks/code.mock.go
	@mockgen -source=internal/repository/dao/user.go -package=daomocks -destination=internal/repository/dao/mocks/user.mock.go
	@mockgen -source=internal/repository/cache/user.go -package=cachemocks -destination=internal/repository/cache/mocks/user.mock.go
	@mockgen -source=pkg/ratelimit/types.go -package=limitmocks -destination=pkg/ratelimit/mocks/ratelimit.mock.go
	@mockgen -package=redismocks -destination=internal/repository/cache/redismocks/cmdable.mock.go github.com/redis/go-redis/v9 Cmdable
	@go mod tidy
docker:
	@rm webook || true
	@go mod tidy
	@GOOS=linux GOARCH=arm go build -tags=k8s -o webook .
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
remote:
	@rm webook || true
	@go mod tidy
	@GOOS=linux GOARCH=amd64 go build -tags=k8s -o webook .
.PHONY: gprc
grpc:
	@buf generate api/proto
