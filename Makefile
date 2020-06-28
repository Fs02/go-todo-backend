export RELEASE_VERSION	?= $(shell git show -q --format=%h)
export DOCKER_REGISTRY	?= docker.pkg.github.com/fs02/go-todo-backend
export DEPLOY			?= api

all: build start
bundle:
	bundle install --path bundle
db-create: bundle
	bundle exec rake db:create
db-migrate: bundle
	bundle exec rake db:migrate
db-rollback: bundle
	bundle exec rake db:rollback
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
