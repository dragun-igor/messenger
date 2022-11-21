FROM golang:latest
RUN mkdir /app
RUN mkdir /app/cmd/
RUN mkdir /app/internal/
RUN mkdir /app/cmd/server
RUN mkdir /app/internal/server
RUN mkdir /app/proto
RUN mkdir /app/pkg
RUN mkdir /app/migrations
RUN mkdir /app/config
RUN mkdir /app/internal/pkg
COPY /cmd/server /app/cmd/server
COPY /internal/server /app/internal/server
COPY /proto /app/proto
RUN rm /app/proto/proto.bat
COPY /migrations /app/migrations
COPY /pkg /app/pkg
COPY /config /app/config
COPY /internal/pkg /app/internal/pkg
COPY go.mod /app
WORKDIR /app
RUN go mod tidy
RUN go build -o main cmd/server/main.go
CMD ["/app/main"]