FROM golang:1.19-alpine as builder
WORKDIR /app
ENV GO111MODULE=on

ARG version
ENV MESSENGER_APP_VERSION=$version

COPY . .
RUN go mod download
RUN go build -o messenger_client cmd/client/main.go

FROM alpine
RUN apk update && \
    adduser -D -H -h /app messenger_client && \
    mkdir -p /app/migrations && \
    chown -R messenger_client:messenger_client /app
WORKDIR /app
USER messenger_client

COPY --chown=messenger_client --from=builder /app/messenger_client /app
COPY --chown=messenger_client --from=builder /app/migrations /app/migrations

EXPOSE 5000
CMD ["/app/messenger_client"]