#!/bin/bash

#----------------------------------------------------------------------------------
#USER
#----------------------------------------------------------------------------------
echo "Генерация protobuf - user"
#(user server)
protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  services/user/protobuf/user.proto --go_out=plugins=grpc:.

#----------------------------------------------------------------------------------
# Post
#----------------------------------------------------------------------------------
echo "Генерация protobuf - Post"
#(post server)
protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  services/post/protobuf/post.proto --go_out=plugins=grpc:.

#(client для user сервиса)
 protoc -I/usr/local/include -I. \
   -I$GOPATH/src \
   -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
   services/user/protobuf/user.proto --go_out=plugins=grpc:./services/post/protobuf/

#(client для category сервиса)
 protoc -I/usr/local/include -I. \
   -I$GOPATH/src \
   -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
   services/category/protobuf/category.proto --go_out=plugins=grpc:./services/post/protobuf/

#(client для comment сервиса)
 protoc -I/usr/local/include -I. \
   -I$GOPATH/src \
   -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
   services/comment/protobuf/comment.proto --go_out=plugins=grpc:./services/post/protobuf/

#----------------------------------------------------------------------------------
# Comment
#----------------------------------------------------------------------------------
echo "Генерация protobuf - Comment"
protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  services/comment/protobuf/comment.proto --go_out=plugins=grpc:.

#Category
echo "Генерация protobuf - Category"
protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  services/category/protobuf/category.proto --go_out=plugins=grpc:.

#API
echo "Генерация protobuf - api-gw"
protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --grpc-gateway_out=logtostderr=true:services/api-gw \
  --go_out=plugins=grpc:services/api-gw \
  services/user/protobuf/user.proto 

  protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --grpc-gateway_out=logtostderr=true:services/api-gw \
  --go_out=plugins=grpc:services/api-gw \
  services/post/protobuf/post.proto 

  protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --grpc-gateway_out=logtostderr=true:services/api-gw \
  --go_out=plugins=grpc:services/api-gw \
  services/comment/protobuf/comment.proto 

  protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --grpc-gateway_out=logtostderr=true:services/api-gw \
  --go_out=plugins=grpc:services/api-gw \
  services/category/protobuf/category.proto 

#SWAGGER
echo "Генерация protobuf - swagger-ui"
protoc -I/usr/local/include -I. \
    -I$GOPATH/src \
    -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
    --swagger_out=logtostderr=true:web/frontend/src/swagger-ui/ \
    services/user/protobuf/user.proto

protoc -I/usr/local/include -I. \
    -I$GOPATH/src \
    -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
    --swagger_out=logtostderr=true:web/frontend/src/swagger-ui/ \
    services/post/protobuf/post.proto

protoc -I/usr/local/include -I. \
    -I$GOPATH/src \
    -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
    --swagger_out=logtostderr=true:web/frontend/src/swagger-ui/ \
    services/comment/protobuf/comment.proto

protoc -I/usr/local/include -I. \
    -I$GOPATH/src \
    -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
    --swagger_out=logtostderr=true:web/frontend/src/swagger-ui/ \
    services/category/protobuf/category.proto


#pretty-swag -i web/frontend/src/swagger-ui/services/user/protobuf/user.swagger.json -o web/frontend/public/swagger/user.html 
#pretty-swag -i web/frontend/src/swagger-ui/services/post/protobuf/post.swagger.json -o web/frontend/public/swagger/post.html
#pretty-swag -i web/frontend/src/swagger-ui/services/comment/protobuf/comment.swagger.json -o web/frontend/public/swagger/comment.html
#pretty-swag -i web/frontend/src/swagger-ui/services/category/protobuf/category.swagger.json -o web/frontend/public/swagger/category.html