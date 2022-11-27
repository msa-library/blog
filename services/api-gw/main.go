package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"net/http"
	"os"
	"time"

	categoryService "./services/category/protobuf"
	commentService "./services/comment/protobuf"
	postService "./services/post/protobuf"
	userService "./services/user/protobuf"

	app "./app"
)

var (
	// gRPC services
	userServerAdress     = fmt.Sprintf("%s:%s", os.Getenv("USER_HOST"), os.Getenv("USER_PORT"))
	postServerAdress     = fmt.Sprintf("%s:%s", os.Getenv("POST_HOST"), os.Getenv("POST_PORT"))
	commentServerAdress  = fmt.Sprintf("%s:%s", os.Getenv("COMMENT_HOST"), os.Getenv("COMMENT_PORT"))
	categoryServerAdress = fmt.Sprintf("%s:%s", os.Getenv("CATEGORY_HOST"), os.Getenv("CATEGORY_PORT"))
)

func main() {
	proxyAddr := fmt.Sprintf(":%s", os.Getenv("PORT"))
	HTTPProxy(proxyAddr)
}

func HTTPProxy(proxyAddr string) {

	grpcGwMux := runtime.NewServeMux()

	//----------------------------------------------------------------
	// настройка подключений со стороны gRPC
	//----------------------------------------------------------------
	//Подключение к сервису User
	grpcUserConn, err := grpc.Dial(
		userServerAdress,
		//grpc.WithPerRPCCredentials(&reqData{}),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalln("Filed to connect to User service", err)
	}
	defer grpcUserConn.Close()

	err = userService.RegisterUserServiceHandler(
		context.Background(),
		grpcGwMux,
		grpcUserConn,
	)
	if err != nil {
		log.Fatalln("Filed to start HTTP server", err)
	}

	//----------------------------------------------------------------
	//Подключение к сервису Post
	grpcPostConn, err := grpc.Dial(
		postServerAdress,
		//grpc.WithPerRPCCredentials(&reqData{}),
		grpc.WithUnaryInterceptor(AccessLogInterceptor),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalln("Filed to connect to Post service", err)
	}
	defer grpcPostConn.Close()

	err = postService.RegisterPostServiceHandler(
		context.Background(),
		grpcGwMux,
		grpcPostConn,
	)
	if err != nil {
		log.Fatalln("Filed to start HTTP server", err)
	}

	//----------------------------------------------------------------
	//Подключение к сервису Comment
	grpcCommentConn, err := grpc.Dial(
		commentServerAdress,
		//grpc.WithPerRPCCredentials(&reqData{}),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalln("Filed to connect to Comment service", err)
	}
	defer grpcCommentConn.Close()

	err = commentService.RegisterCommentServiceHandler(
		context.Background(),
		grpcGwMux,
		grpcCommentConn,
	)
	if err != nil {
		log.Fatalln("Filed to start HTTP server", err)
	}

	//----------------------------------------------------------------
	//Подключение к сервису Category
	grpcCategoryConn, err := grpc.Dial(
		categoryServerAdress,
		//grpc.WithPerRPCCredentials(&reqData{}),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalln("Filed to connect to Category service", err)
	}
	defer grpcCategoryConn.Close()

	err = categoryService.RegisterCategoryServiceHandler(
		context.Background(),
		grpcGwMux,
		grpcCategoryConn,
	)
	if err != nil {
		log.Fatalln("Filed to start HTTP server", err)
	}

	//----------------------------------------------------------------
	//	Настройка маршрутов с стороны REST
	//----------------------------------------------------------------
	mux := http.NewServeMux()

	mux.Handle("/api/v1/", grpcGwMux)
	mux.HandleFunc("/", helloworld)

	fmt.Println("starting HTTP server at " + proxyAddr)
	log.Fatal(http.ListenAndServe(proxyAddr, mux))
}

func helloworld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "URL:", r.URL.String())
}

func AccessLogInterceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	md, _ := metadata.FromOutgoingContext(ctx)
	start := time.Now()

	var traceId, userId, userRole string
	if len(md["authorization"]) > 0 {
		tokenString := md["authorization"][0]
		if tokenString != "" {
			err, token := userService.CheckGetJWTToken(tokenString)
			if err != nil {
				return err
			}
			userId = fmt.Sprintf("%s", token["UserID"])
			userRole = fmt.Sprintf("%s", token["UserRole"])
		} else {
			return errors.New("error authorization")
		}
	} else {
		return errors.New("error authorization")
	}
	//Присваиваю ID запроса
	traceId = fmt.Sprintf("%d", time.Now().UTC().UnixNano())

	callContext := context.Background()
	mdOut := metadata.Pairs(
		"trace-id", traceId,
		"user-id", userId,
		"user-role", userRole,
	)
	callContext = metadata.NewOutgoingContext(callContext, mdOut)

	err := invoker(callContext, method, req, reply, cc, opts...)

	msg := fmt.Sprintf("Call:%v, traceId: %v, userId: %v, userRole: %v, time: %v", method, traceId, userId, userRole, time.Since(start))
	app.AccesLog(msg)

	return err
}
