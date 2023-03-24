package main

import (
	"context"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	jwtv4 "github.com/golang-jwt/jwt/v4"
	v1 "github.com/noovertime7/dailytest/helloworld/api/helloworld/v1"
	clientv3 "go.etcd.io/etcd/client/v3"
	srcgrpc "google.golang.org/grpc"
	"log"
	"time"
)

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
	serviceTestKey := "testKey"

	cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	})

	if err != nil {
		panic(err)
	}
	r := etcd.New(cli)

	claims := &CustomClaims{
		BaseClaims: BaseClaims{
			ID:          1,
			Username:    "user1",
			NickName:    "nickname",
			AuthorityId: 111,
			TenantID:    "aaaa",
		},
		StandardClaims: jwtv4.StandardClaims{},
	}

	jwt.WithClaims(func() jwtv4.Claims { return claims })

	connGRPC, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///helloworld"),
		grpc.WithDiscovery(r),
		grpc.WithMiddleware(
			jwt.Client(func(token *jwtv4.Token) (interface{}, error) {
				return []byte(serviceTestKey), nil
			}, jwt.WithClaims(func() jwtv4.Claims {
				return claims
			})),
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer connGRPC.Close()

	//connHTTP, err := http.NewClient(
	//	context.Background(),
	//	http.WithEndpoint("discovery:///helloworld"),
	//	http.WithDiscovery(r),
	//	http.WithBlock(),
	//	http.WithMiddleware(
	//		jwt.Client(func(token *jwtv4.Token) (interface{}, error) {
	//			return []byte(serviceTestKey), nil
	//		}),
	//	),
	//)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer connHTTP.Close()

	for {
		//callHTTP(connHTTP)
		callGRPC(connGRPC)
		time.Sleep(time.Second)
	}
}

func callGRPC(conn *srcgrpc.ClientConn) {
	client := v1.NewGreeterClient(conn)
	reply, err := client.SayHello(context.Background(), &v1.HelloRequest{Name: "kratos"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[grpc] SayHello %+v\n", reply)
}

func callHTTP(conn *http.Client) {
	client := v1.NewGreeterHTTPClient(conn)
	reply, err := client.SayHello(context.Background(), &v1.HelloRequest{Name: "kratos"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[http] SayHello %+v\n", reply)
}

func CreateAccessJwtToken(secretKey []byte) string {
	claims := jwtv4.NewWithClaims(jwtv4.SigningMethodHS256,
		jwtv4.MapClaims{
			"id": "aaa",
		})

	signedToken, err := claims.SignedString(secretKey)
	if err != nil {
		return ""
	}

	return signedToken
}
