package client

import (
	"fmt"
	"google.golang.org/grpc"
	"time"
)

func newGrpcConn(host string) *grpc.ClientConn {
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	for err != nil {
		fmt.Println("grpc dial", "err", err)
		time.Sleep(time.Millisecond * 100)
		conn, err = grpc.Dial(host, grpc.WithInsecure())
	}
	return conn
}

