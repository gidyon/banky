package bank

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gidyon/banky/pkg/api/bank"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc/grpclog"
)

type bankAPIServer struct {
	sqlDB  *gorm.DB
	logger grpclog.LoggerV2
}

// Options contains parameters for NewBankAPI
type Options struct {
	SQLDB  *gorm.DB
	Logger grpclog.LoggerV2
}

// NewBankAPI is a factory for creating bank API singleton
func NewBankAPI(opt *Options) (bank.BankAPIServer, error) {
	// Validation
	switch {
	case opt.SQLDB == nil:
		return nil, status.Error(codes.InvalidArgument, "nil database not allowed")
	case opt.Logger == nil:
		return nil, status.Error(codes.InvalidArgument, "nil logger not allowed")
	}

	api := &bankAPIServer{
		sqlDB:  opt.SQLDB,
		logger: opt.Logger,
	}

	// Auto migration
	err := api.sqlDB.AutoMigrate(&Bank{}).Error
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed to automigrate table")
	}

	return api, nil
}

func (api *bankAPIServer) CreateBank(
	ctx context.Context, createReq *bank.CreateBankRequest,
) (*bank.Bank, error) {
	// Request must not be nil
	if createReq == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	// Validation
	switch {
	case createReq.Bank == nil:
		return nil, status.Error(codes.InvalidArgument, "missing bank")
	case createReq.Bank.Name == "":
		return nil, status.Error(codes.InvalidArgument, "missing bank name")
	case createReq.Bank.Address == "":
		return nil, status.Error(codes.InvalidArgument, "missing bank address")
	}

	// Get bank model
	bankDB, err := GetBankDB(createReq.Bank)
	if err != nil {
		api.logger.Errorln(err)
		return nil, status.Error(codes.Internal, "Failed to get model for bank")
	}

	// Save the object to database
	err = api.sqlDB.Create(bankDB).Error
	if err != nil {
		api.logger.Errorln(err)
		return nil, status.Error(codes.Internal, "Failed save bank")
	}

	createReq.Bank.BankId = fmt.Sprint(bankDB.ID)

	return createReq.Bank, nil
}

func (api *bankAPIServer) GetBank(
	ctx context.Context, getReq *bank.GetBankRequest,
) (*bank.Bank, error) {
	// Request must not be nil
	if getReq == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	// Validation
	switch {
	case getReq.BankId == "":
		return nil, status.Error(codes.InvalidArgument, "missing bank id")
	}

	bankDB := &Bank{}
	err := api.sqlDB.First(bankDB, "id=?", getReq.BankId).Error
	if err != nil {
		api.logger.Errorln(err)
		return nil, status.Error(codes.Internal, "Failed to get bank")
	}

	// Get bank protobuf message
	bankPB, err := GetBankPB(bankDB)
	if err != nil {
		api.logger.Errorln(err)
		return nil, status.Error(codes.Internal, "Failed to convert bank to message")
	}

	return bankPB, nil
}
