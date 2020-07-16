package main

import (
	"context"
	"net"
	"time"

	echo "github.com/gidyon/banky/examples/echo/pkg"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type echoService string

// context originates from client
func (echoService) Echo(
	ctx context.Context, echoMsg *echo.Payload,
) (*echo.Payload, error) {
	logrus.Infoln("starting some expensive work!")

	select {
	case <-time.After(5 * time.Second): // This rep some expensive operation
		logrus.Infoln("operation completed!")
		return echoMsg, nil
	case <-ctx.Done():
		logrus.Infoln("operation failed to complete: ", ctx.Err())
		return nil, ctx.Err()
	}
}

func handlerErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// grpc server
	srv := grpc.NewServer()

	echo.RegisterEchoServiceServer(srv, echoService("hello"))

	lis, err := net.Listen("tcp", ":9090")
	handlerErr(err)

	logrus.Infoln("server started on port :9090")
	srv.Serve(lis)
}
