all: build start
dep:
	go get github.com/Fs02/kamimai
migrate:
	export $$(cat .env | grep -v ^\# | xargs) && \
	kamimai --driver=mysql --dsn="mysql://$$MYSQL_USERNAME:$$MYSQL_PASSWORD@($$MYSQL_HOST:$$MYSQL_PORT)/$$MYSQL_DATABASE" --directory=./db/migrations sync
rollback:
	export $$(cat .env | grep -v ^\# | xargs) && \
	kamimai --driver=mysql --dsn="mysql://$$MYSQL_USERNAME:$$MYSQL_PASSWORD@($$MYSQL_HOST:$$MYSQL_PORT)/$$MYSQL_DATABASE" --directory=./db/migrations down
gen:
	go generate ./...
build: gen
	go build -mod=vendor -o bin/api ./cmd/api
test: gen
	go test -mod=vendor -race ./...
start:
	export $$(cat .env | grep -v ^\# | xargs) && ./bin/api
