FROM golang:latest
RUN mkdir /app
RUN mkdir /app/cmd/
RUN mkdir /app/cmd/client
RUN mkdir /app/internal/
RUN mkdir /app/internal/client
RUN mkdir /app/proto
COPY ./cmd/client /app/cmd/client
COPY ./internal/client /app/internal/client
COPY ./proto /app/proto
RUN rm /app/proto/proto.bat
COPY go.mod /app
WORKDIR /app
RUN go mod tidy
RUN go build -o main cmd/client/main.go
CMD ["/app/main"]
