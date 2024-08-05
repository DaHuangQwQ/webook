.PHONY: docker
docker:
	@rm webook || true
	@go mod tidy
	@GOOS=linux GOARCH=arm go build -tags=k8s -o webook .
	@docker rmi -f dahuang/webook:v0.0.1
	@docker build -t dahuang/webook:v0.0.1 .

k8s:
	@kubectl delete -f ./k8s-webook-service.yaml || true
	@kubectl apply -f ./k8s-webook-service.yaml
mock:
	@mockgen -source=internal/service/code.go -package=svcmocks -destination=internal/service/mocks/code.mock.go
	@mockgen -source=internal/service/user.go -package=svcmocks -destination=internal/service/mocks/user.mock.go
	@go mod tidy