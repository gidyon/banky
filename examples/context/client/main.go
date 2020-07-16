package main

import (
	"context"
	"time"

	echo "github.com/gidyon/banky/examples/echo/pkg"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func handlErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	cc, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
	handlErr(err)

	client := echo.NewEchoServiceClient(cc)

	// timeout (this context will be passed to the server goroutine)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	res, err := client.Echo(ctx, &echo.Payload{Message: "Hello there"})
	if err != nil {
		logrus.Errorln(err)
		return
	}

	logrus.Infof("response was: %s\n", res.Message)

}
