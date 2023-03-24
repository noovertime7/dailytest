package main

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2"
	server2 "github.com/noovertime7/dailytest/kratos_jwt_demo/server"
	"log"

	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	jwtv4 "github.com/golang-jwt/jwt/v4"
	"github.com/noovertime7/dailytest/helloworld/api/helloworld/v1"
)

type server struct {
	v1.UnimplementedGreeterServer

	hc v1.GreeterClient
}

func (s *server) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	claims, ok := jwt.FromContext(ctx)
	fmt.Println(claims, ok)
	return &v1.HelloReply{Message: "hello from service"}, nil
}

type BaseClaims struct {
	ID          int
	Username    string
	NickName    string
	AuthorityId uint
	TenantID    string
}

// CustomClaims 自定义token中携带的信息
type CustomClaims struct {
	BaseClaims
	jwtv4.StandardClaims
}

func main() {
	testKey := "testKey"
	//httpSrv := http.NewServer(
	//	http.Address(":8000"),
	//	http.Middleware(
	//		jwt.Server(func(token *jwtv4.Token) (interface{}, error) {
	//			return []byte(testKey), nil
	//		}),
	//	),
	//)
	grpcSrv := grpc.NewServer(
		grpc.Address(":9000"),
		grpc.Middleware(
			jwt.Server(func(token *jwtv4.Token) (interface{}, error) {
				return []byte(testKey), nil
			}),
		),
	)
	serviceTestKey := "serviceTestKey"

	jwt.WithClaims(func() jwtv4.Claims { return &CustomClaims{} })

	con, _ := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("dns:///127.0.0.1:9000"),
		grpc.WithMiddleware(
			jwt.Client(func(token *jwtv4.Token) (interface{}, error) {
				return []byte(serviceTestKey), nil
			}),
		),
	)
	s := &server{
		hc: v1.NewGreeterClient(con),
	}

	//注册中心
	register, _, err := server2.RegistrarServer()
	if err != nil {
		log.Fatal(err)
	}

	v1.RegisterGreeterServer(grpcSrv, s)
	//v1.RegisterGreeterHTTPServer(httpSrv, s)
	app := kratos.New(
		kratos.Name("helloworld"),
		kratos.Server(
			//httpSrv,
			grpcSrv,
		),
		kratos.Registrar(register),
	)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
