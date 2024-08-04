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