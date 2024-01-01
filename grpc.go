package goeximbaselib

import (
	"log"
	"net"

	"github.com/exim-id/go-exim-base-lib/env"
	"github.com/exim-id/go-exim-base-lib/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ClientGrpcService[T interface{}](port string, cli func(*grpc.ClientConn) T) T {
	conn, err := grpc.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer errors.StopConnClientGrpc(conn)
	if err != nil {
		panic(err)
	}
	return cli(conn)
}

func RegisterGrpcServer(registering func(*grpc.Server)) {
	svc := grpc.NewServer()
	defer errors.StopGrpcServer(svc)
	registering(svc)
	l, err := net.Listen("tcp", env.GetGrpcPort())
	if err != nil {
		panic(err)
	}
	log.Fatal(svc.Serve(l))
}
