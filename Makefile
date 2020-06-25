all: build start
dep:
	go get github.com/Fs02/kamimai
migrate:
	export $$(cat .env | grep -v ^\# | xargs) && \
	kamimai --driver=mysql --dsn="mysql://$$MYSQL_USERNAME:$$MYSQL_PASSWORD@($$MYSQL_HOST:$$MYSQL_PORT)/$$MYSQL_DATABASE" --directory=./migrations sync
rollback:
	export $$(cat .env | grep -v ^\# | xargs) && \
	kamimai --driver=mysql --dsn="mysql://$$MYSQL_USERNAME:$$MYSQL_PASSWORD@($$MYSQL_HOST:$$MYSQL_PORT)/$$MYSQL_DATABASE" --directory=./migrations down
gen:
	go generate ./...
build: gen
	go build -o bin/api .
start:
	export $$(cat .env | grep -v ^\# | xargs) && ./bin/api
