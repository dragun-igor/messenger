# FROM golang:latest
# RUN mkdir /app
# RUN mkdir /app/cmd/
# RUN mkdir /app/cmd/client
# RUN mkdir /app/internal/
# RUN mkdir /app/internal/client
# RUN mkdir /app/proto
# RUN mkdir /app/internal/pkg
# COPY ./cmd/client /app/cmd/client
# COPY ./internal/client /app/internal/client
# COPY ./proto /app/proto
# COPY go.mod /app
# COPY /internal/pkg /app/internal/pkg
# WORKDIR /app
# RUN go mod tidy
# RUN go build -o main cmd/client/main.go
# CMD ["/app/main"]

# FROM golang:latest
# RUN mkdir /app
# RUN mkdir /app/cmd/
# RUN mkdir /app/internal/
# RUN mkdir /app/cmd/server
# RUN mkdir /app/internal/server
# RUN mkdir /app/proto
# RUN mkdir /app/pkg
# RUN mkdir /app/migrations
# RUN mkdir /app/config
# RUN mkdir /app/internal/pkg
# COPY /cmd/server /app/cmd/server
# COPY /internal/server /app/internal/server
# COPY /proto /app/proto
# COPY /migrations /app/migrations
# COPY /pkg /app/pkg
# COPY /config /app/config
# COPY /internal/pkg /app/internal/pkg
# COPY go.mod /app
# WORKDIR /app
# RUN go mod tidy
# RUN go build -o main cmd/server/main.go
# CMD ["/app/main"]

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