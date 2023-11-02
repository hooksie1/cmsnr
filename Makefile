PROJECT_NAME := "cmsnr"
PKG := "github.com/hooksie1/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)
VERSION := $(shell if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then git describe --exact-match --tags HEAD 2>/dev/null || echo "dev-$(shell git rev-parse --short HEAD)"; else echo "dev"; fi)


.PHONY: all build docker dep clean test coverage lint

all: build

lint: ## Lint the files
	@golint -set_exit_status ./...

test: ## Run unittests
	@go test ./...

coverage:
	@go test -cover ./...
	@go test ./... -coverprofile=cover.out && go tool cover -html=cover.out -o coverage.html

dep: ## Get the dependencies
	@go get -u golang.org/x/lint/golint

build: linux windows mac

linux: dep
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags "-w -X '$(PKG)/cmd.Version=$(VERSION)'" -o $(PROJECT_NAME)ctl

windows: dep
	CGO_ENABLED=0 GOOS=windows go build -a -ldflags "-w -X '$(PKG)/cmd.Version=$(VERSION)'" -o $(PROJECT_NAME)ctl.exe

mac: dep
	CGO_ENABLED=0 GOOS=darwin go build -a -ldflags "-s -w -X '$(PKG)/cmd.Version=$(VERSION)'" -o $(PROJECT_NAME)ctl-darwin

cmsnrctl: ## Builds the binary on the current platform
	go build -mod=vendor -a -ldflags "-w -X '$(PKG)/cmd.Version=$(VERSION)'" -o $(PROJECT_NAME)ctl

docker-local: ## Builds the container image and pushes to the local k8s registry
	docker build -t localhost:50000/cmsnr:latest .
	docker push localhost:50000/cmsnr:latest

docker-delete: ## Deletes the local docker image
	docker image rm localhost:50000/cmsnr:latest

update-local: docker-local ## Builds the container image and pushes to registry, rolls out the new container into the cluster
	kubectl rollout restart deployment/cmsnr-mutating-webhook
	kubectl rollout restart deployment/cmsnr-validating-webhook

deploy-local: k8s-up cmsnrctl docker-local ## Creates a local k8s cluster, builds a docker image of cmsnr, and pushes to local registry
	./cmsnrctl server deploy --registry k3d-cmsnr-registry:50000 --version latest| kubectl apply -f -
	kubectl wait pods -l app=cmsnr-mutating-webhook --for condition=Ready --timeout=30s
	kubectl wait pods -l app=cmsnr-validatin-webhook --for condition=Ready --timeout=30s

k8s-up: ## Creates a local kubernetes cluster with a registry
	k3d registry create cmsnr-registry --port 50000
	k3d cluster create cmsnr --registry-use k3d-cmsnr-registry:50000 --servers 3 -p "8080:80@loadbalancer"

k8s-down: ## Destroys the k8s cluster and registry
	k3d registry delete cmsnr-registry
	k3d cluster delete cmsnr

clean: ## Remove previous build
	git clean -fd
	git clean -fx
	git reset --hard

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
