export RELEASE_VERSION	?= $(shell git show -q --format=%h)
export DOCKER_REGISTRY	?= docker.pkg.github.com/fs02/go-todo-backend
export DEPLOY			?= api

all: build start
db-migrate:
	export $$(cat .env | grep -v ^\# | xargs) && go run cmd/db/main.go migrate
db-rollback:
	export $$(cat .env | grep -v ^\# | xargs) && go run cmd/db/main.go rollback
gen:
	go generate ./...
build: gen
	go build -mod=vendor -o bin/api ./cmd/api
test: gen
	go test -mod=vendor -race ./...
start:
	export $$(cat .env | grep -v ^\# | xargs) && ./bin/api
docker:
	docker build -t $(DOCKER_REGISTRY)/$(DEPLOY):$(RELEASE_VERSION) -f ./deploy/$(DEPLOY)/Dockerfile .
push:
	docker push $(DOCKER_REGISTRY)/$(DEPLOY):$(RELEASE_VERSION)
