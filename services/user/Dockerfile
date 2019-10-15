FROM golang:1.12.6

WORKDIR /app

RUN mkdir -p /app/src
COPY . /app

# RUN go get "github.com/go-redis/redis"
RUN go get "go.mongodb.org/mongo-driver/mongo"
RUN go get "google.golang.org/grpc"

CMD go run main.go

EXPOSE 9000