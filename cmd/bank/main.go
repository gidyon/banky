package main

import (
	"context"

	"github.com/gidyon/micros/utils/healthcheck"

	"github.com/gidyon/micros"
	"github.com/gidyon/micros/pkg/config"
	app_grpc_middleware "github.com/gidyon/micros/pkg/grpc/middleware"

	bank_app "github.com/gidyon/banky/internal/services/bank"
	"github.com/gidyon/banky/pkg/api/bank"
	"github.com/sirupsen/logrus"
)

// Configuration - [database configuration, tls security, external services]
// Customization - Customize http server, grpc server
// Bootstrapping

func main() {
	// 1. Initialize config
	cfg, err := config.New(config.FromFile)
	handleErr(err)

	ctx := context.Background()

	// 2. Create micro service
	bankSrv, err := micros.NewService(ctx, cfg, nil)
	handleErr(err)

	// 3. Add custom options to service
	// Add Recovery middleware
	recoveryUIs, recoverySIs := app_grpc_middleware.AddRecovery()
	bankSrv.AddGRPCUnaryServerInterceptors(recoveryUIs...)
	bankSrv.AddGRPCStreamServerInterceptors(recoverySIs...)

	// Set base endpoint for runtime mux
	bankSrv.SetServeMuxEndpoint("/")

	// Update readiness health check
	bankSrv.AddEndpoint("/health/ready", healthcheck.RegisterProbe(&healthcheck.ProbeOptions{
		Service:      bankSrv,
		Type:         healthcheck.ProbeReadiness,
		AutoMigrator: func() error { return nil },
	}))

	// Update liveness health check
	bankSrv.AddEndpoint("/health/live", healthcheck.RegisterProbe(&healthcheck.ProbeOptions{
		Service:      bankSrv,
		Type:         healthcheck.ProbeLiveNess,
		AutoMigrator: func() error { return nil },
	}))

	// 4. Starts the service
	bankSrv.Start(ctx, func() error {
		// API implemetation
		bankAPI, err := bank_app.NewBankAPI(&bank_app.Options{
			SQLDB:  bankSrv.GormDB(),
			Logger: bankSrv.Logger(),
		})
		handleErr(err)

		// Register Bank API implementation to gRPC server
		bank.RegisterBankAPIServer(bankSrv.GRPCServer(), bankAPI)

		// Register grpc-gateway muxer with gRPC server over network connection
		handleErr(bank.RegisterBankAPIHandler(ctx, bankSrv.RuntimeMux(), bankSrv.ClientConn()))

		return nil
	})
}

func handleErr(err error) {
	if err != nil {
		logrus.Fatalln(err)
	}
}
