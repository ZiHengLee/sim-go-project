package rpc

import (
	"context"
	"fmt"
	"github.com/capell/capell_scan/lib/discovery"
	"github.com/capell/capell_scan/proto/swap"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"log"
	"time"
)

var (
	Register   *discovery.Resolver
	ctx        context.Context
	CancelFunc context.CancelFunc

	SwapClient swap.MsgClient
)

func Init() {
	Register = discovery.NewResolver([]string{"0.0.0.0:2379"}, logrus.New())
	resolver.Register(Register)
	ctx, CancelFunc = context.WithTimeout(context.Background(), 3*time.Second)

	defer Register.Close()
	initClient("swap", &SwapClient)
}

func initClient(serviceName string, client interface{}) {
	conn, err := connectServer(serviceName)

	if err != nil {
		panic(err)
	}
	switch c := client.(type) {
	case *swap.MsgClient:
		*c = swap.NewMsgClient(conn)
	default:
		panic("unsupported client type")
	}
}

func connectServer(serviceName string) (conn *grpc.ClientConn, err error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	addr := fmt.Sprintf("%s:///%s", Register.Scheme(), serviceName)

	// Load balance
	if true {
		log.Printf("load balance enabled for %s\n", serviceName)
		opts = append(opts, grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, "round_robin")))
	}

	conn, err = grpc.DialContext(ctx, addr, opts...)
	return
}
