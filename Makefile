BINARY_NAME=api

unit-test:
	@echo "Installing tparse (if needed)..."
	@go install github.com/mfridman/tparse@v0.11.1
	@echo "Running unit tests..."
	@go test ./... -cover -json | tparse

build:
	@echo "Building binary..."
	@go build -o bin/$(BINARY_NAME)
	@echo "Done!"

run: build
	@echo "Running binary..."
	@./bin/$(BINARY_NAME)

# Temporary deployment scripts that will be replaced by a CI/CD pipeline
build_docker:
	@echo "Building docker image..."
	@docker build -t webstradev/rsdb-backend:latest .
	@echo "Pushing docker image..."
	@docker push webstradev/rsdb-backend:latest

deploy-prod:
	@echo "Deleting old k8s production deployment"
	-@kubectl delete deployment rsdb-backend
	@echo "Deploying new k8s production deployment"
	@kubectl apply -f kube/deployment.yaml
	@echo "Done!"

deploy:
	@echo "Deleting old k8s deployment"
	-@kubectl delete deployment rsdb-dev-backend
	@echo "Deploying new k8s deployment"
	@kubectl apply -f kube/dev-deployment.yaml
	@echo "Done!"


dev: build_docker deploy
	@echo "Done!"

prod: build_docker deploy-prod
	@echo "Done!"