package models

import (
	u "go-card/utils"

	"github.com/jinzhu/gorm"
)

type Action struct {
	gorm.Model
	Type          string `json:"type"`
	TransactionId int    `json:"transactionId"`
	Amount        int    `json:"amount"`
}

func (action *Action) GetTransaction() (*Transaction, map[string]interface{}) {
	transaction := &Transaction{}
	err := GetDB().Table("transactions").Where("id = ?", action.TransactionId).First(transaction).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, u.Message(false, "Connection error. Please retry")
	}

	return transaction, nil
}

func (action *Action) Capture() map[string]interface{} {
	action.Type = "CAPTURE"
	transaction, err := action.GetTransaction()
	if err != nil {
		resp := u.Message(false, "Transaction Not Found")
		return resp
	}

	err = transaction.CaptureFunds(action.Amount, action.TransactionId)
	if err != nil {
		resp := u.Message(false, "Transaction already paid")
		return resp
	}

	GetDB().Create(action)

	resp := u.Message(true, "success")
	resp["transaction"] = transaction
	return resp
}

func (action *Action) ChangeAuthorised() map[string]interface{} {
	action.Type = "CHANGE"
	transaction, err := action.GetTransaction()
	if err != nil {
		resp := u.Message(false, "Transaction Not Found")
		return resp
	}

	err = transaction.ChangeAuthorisedFunds(action.Amount, action.TransactionId)
	if err != nil {
		resp := u.Message(false, "Transaction already paid")
		return resp
	}

	GetDB().Create(action)

	resp := u.Message(true, "success")
	return resp
}

func (action *Action) Refund() map[string]interface{} {
	action.Type = "REFUND"
	transaction, err := action.GetTransaction()
	if err != nil {
		resp := u.Message(false, "Transaction Not Found")
		return resp
	}

	if action.Amount <= transaction.Captured {
		account, err := transaction.GetAccount()
		if err != nil {
			resp := u.Message(false, "Transaction Not Found")
			return resp
		}
		transaction.Refund(action.Amount)
		account.AdjustFunds(action.Amount)
	} else {
		resp := u.Message(false, "Can Not Refund More Than Captured")
		return resp
	}

	GetDB().Create(action)

	resp := u.Message(true, "Refunded")
	return resp
}
