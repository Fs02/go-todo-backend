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
