package main

import (
	"context"
	"github.com/go-kratos/kratos/v2"
	"log"

	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	jwtv4 "github.com/golang-jwt/jwt/v4"
	"github.com/noovertime7/dailytest/helloworld/api/helloworld/v1"
)

type server struct {
	v1.UnimplementedGreeterServer

	hc v1.GreeterClient
}

func (s *server) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	return &v1.HelloReply{Message: "hello from service"}, nil
}

func main() {
	testKey := "testKey"
	httpSrv := http.NewServer(
		http.Address(":8000"),
		http.Middleware(
			jwt.Server(func(token *jwtv4.Token) (interface{}, error) {
				return []byte(testKey), nil
			}),
		),
	)
	grpcSrv := grpc.NewServer(
		grpc.Address(":9000"),
		grpc.Middleware(
			jwt.Server(func(token *jwtv4.Token) (interface{}, error) {
				return []byte(testKey), nil
			}),
		),
	)
	serviceTestKey := "serviceTestKey"
	con, _ := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("dns:///127.0.0.1:9001"),
		grpc.WithMiddleware(
			jwt.Client(func(token *jwtv4.Token) (interface{}, error) {
				return []byte(serviceTestKey), nil
			}),
		),
	)
	s := &server{
		hc: v1.NewGreeterClient(con),
	}
	v1.RegisterGreeterServer(grpcSrv, s)
	v1.RegisterGreeterHTTPServer(httpSrv, s)
	app := kratos.New(
		kratos.Name("helloworld"),
		kratos.Server(
			httpSrv,
			grpcSrv,
		),
	)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
