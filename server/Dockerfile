FROM golang:latest
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o main server/server.go
CMD ["/app/main"]
