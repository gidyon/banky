package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	"google.golang.org/grpc/credentials"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	bank_app "github.com/gidyon/banky/internal/services/bank"
	"github.com/gidyon/banky/pkg/api/bank"
	"github.com/gidyon/micros"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // side effects
	"google.golang.org/grpc"
)

func main2() {
	ctx := context.Background()

	// Fetching configuration from environment variables
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	mysqlSchema := os.Getenv("MYSQL_SCHEMA")
	mysqlAddress := os.Getenv("MYSQL_ADDRESS")

	// Database connection
	dbClient, err := gorm.Open(
		"mysql", fmt.Sprintf(
			"%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			mysqlUser, mysqlPassword, mysqlAddress, mysqlSchema,
		),
	)
	handleErr(err)

	// Create logger
	logger := micros.NewLogger("bank_service")

	// API implemetation
	bankAPI, err := bank_app.NewBankAPI(&bank_app.Options{
		SQLDB:  dbClient,
		Logger: logger,
	})

	// gRPC server
	creds, err := credentials.NewServerTLSFromFile("cert.pem", "key.pem")
	handleErr(err)

	srv := grpc.NewServer(grpc.Creds(creds))

	opts := []grpc.DialOption{grpc.WithInsecure()}

	port := os.Getenv("PORT")

	// Registered the API to our server
	bank.RegisterBankAPIServer(srv, bankAPI)

	// Gateway proxy
	mux := newRuntimeMux()
	bank.RegisterBankAPIHandlerFromEndpoint(
		ctx, mux, fmt.Sprintf("localhost%s", port), opts,
	)

	// Http server
	err = http.ListenAndServe(port, mux)
	handleErr(err)

	// Network listener
	lis, err := net.Listen("tcp", port)
	handleErr(err)

	logger.Infoln("server started on port 9090")

	// Serves incoming request
	err = srv.Serve(lis)
	handleErr(err)
}

func newRuntimeMux() *runtime.ServeMux {
	return runtime.NewServeMux()
}
