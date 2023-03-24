package main

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	jwtv4 "github.com/golang-jwt/jwt/v4"
	v1 "github.com/noovertime7/dailytest/helloworld/api/helloworld/v1"
	"log"
)

func main() {
	serviceTestKey := "testKey"
	con, _ := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("dns:///127.0.0.1:9000"),
		grpc.WithMiddleware(
			jwt.Client(func(token *jwtv4.Token) (interface{}, error) {
				return []byte(serviceTestKey), nil
			}),
		),
	)
	client := v1.NewGreeterClient(con)

	helloReply, err := client.SayHello(context.TODO(), &v1.HelloRequest{Name: "chenteng"})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(helloReply)

}
