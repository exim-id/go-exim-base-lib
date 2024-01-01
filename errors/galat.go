package errors

import (
	"context"
	"log"
	"runtime/debug"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

func StopConnClientGrpc(conn *grpc.ClientConn) {
	NormalError()
	if conn != nil {
		if err := conn.Close(); err != nil {
			panic(err)
		}
	}
}

func StopGrpcServer(svc *grpc.Server) {
	NormalError()
	svc.Stop()
}

func MongoCursorClose(cursor *mongo.Cursor) {
	NormalError()
	if cursor != nil {
		if err := cursor.Close(context.Background()); err != nil {
			panic(err)
		}
	}
}

func MongoClose(cli *mongo.Client) {
	NormalError()
	if cli != nil {
		if err := cli.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}
}

func ServerStop(e *echo.Echo) {
	NormalError()
	if e != nil {
		if err := e.Close(); err != nil {
			panic(err)
		}
	}
}

func NormalError() {
	if r := recover(); r != nil {
		log.Println("Catch", r)
		log.Println("Error", string(debug.Stack()))
	}
}
