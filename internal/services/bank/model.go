package bank

import (
	"fmt"

	"github.com/gidyon/banky/pkg/api/bank"
	"github.com/jinzhu/gorm"
)

// Bank is a model for a bank
type Bank struct {
	Name    string
	Address string
	gorm.Model
}

// GetBankPB creates a protobuf message from model instance
func GetBankPB(bankDB *Bank) (*bank.Bank, error) {
	return &bank.Bank{
		Name:    bankDB.Name,
		BankId:  fmt.Sprint(bankDB.ID),
		Address: bankDB.Address,
	}, nil
}

// GetBankDB creates a model from protobuf message
func GetBankDB(bankPB *bank.Bank) (*Bank, error) {
	return &Bank{
		Name:    bankPB.Name,
		Address: bankPB.Address,
	}, nil
}
