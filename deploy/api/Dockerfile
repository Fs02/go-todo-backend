# Step 1:
FROM golang:1.13.5-alpine3.11 AS builder

RUN apk update && apk add --no-cache git make

WORKDIR $GOPATH/src/github.com/Fs02/go-todo-backend
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64\
    go build -mod=vendor -ldflags="-w -s" -o /go/bin/api ./cmd/api

# Step 2:
# you can also use scratch here, but I prefer to use alpine because it comes with basic command such as curl useful for debugging.
FROM alpine:3.11

RUN apk update && apk add --no-cache curl ca-certificates
RUN rm -rf /var/cache/apk/*

COPY --from=builder --chown=65534:0 /go/bin/api /go/bin/api

USER 65534
EXPOSE 3000

ENTRYPOINT ["/go/bin/api"]
