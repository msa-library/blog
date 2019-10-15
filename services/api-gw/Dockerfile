FROM golang:1.12.6

WORKDIR /app

RUN mkdir -p /app/src
COPY . /app

RUN go get "google.golang.org/grpc"
RUN go get "github.com/grpc-ecosystem/grpc-gateway/runtime"
RUN go get "github.com/dgrijalva/jwt-go"

CMD go run main.go

EXPOSE 8888