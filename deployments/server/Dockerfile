FROM golang:1.19-alpine as builder
WORKDIR /app
ENV GO111MODULE=on

ARG version
ENV MESSENGER_APP_VERSION=$version

COPY . .
RUN go mod download
RUN go build -o messenger_server cmd/server/main.go

FROM alpine
RUN apk update && \
    adduser -D -H -h /app messenger_server && \
    mkdir -p /app/migrations && \
    chown -R messenger_server:messenger_server /app
WORKDIR /app
USER messenger_server

COPY --chown=messenger_server --from=builder /app/messenger_server /app
COPY --chown=messenger_server --from=builder /app/migrations /app/migrations

EXPOSE 5000
CMD ["/app/messenger_server"]