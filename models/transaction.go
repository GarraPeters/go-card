package models

import (
	u "go-card/utils"

	"github.com/jinzhu/gorm"
)

type Transaction struct {
	gorm.Model
	CardNo   string `json:"cardNo"`
	Merchant string `json:"merchant"`
	Amount   int    `json:"amount"`
	Captured int    `json:"captured"`
	Refunded int    `json:"refunded"`
}

func (transaction *Transaction) Validate() (map[string]interface{}, bool) {

	if transaction.CardNo == "" {
		return u.Message(false, "CardNo should be on the payload"), false
	}

	if transaction.Amount < 0 {
		return u.Message(false, "Amount should be on the payload"), false
	}

	if transaction.Merchant == "" {
		return u.Message(false, "Merchant is not recognized"), false
	}

	//All the required parameters are present
	return u.Message(true, "success"), true
}

func (transaction *Transaction) GetAccount() (*Account, map[string]interface{}) {
	account := &Account{}
	err := GetDB().Table("accounts").Where("card_no = ?", transaction.CardNo).First(account).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, u.Message(false, "Connection error. Please retry")
	}

	return account, nil
}

func (transaction *Transaction) GetTransaction(transactionId int) (*Transaction, map[string]interface{}) {
	err := GetDB().Table("transactions").Where("id = ?", transactionId).First(transaction).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, u.Message(false, "Connection error. Please retry")
	}

	return transaction, nil
}

func (transaction *Transaction) Create() map[string]interface{} {

	if resp, ok := transaction.Validate(); !ok {
		return resp
	}

	account, err := transaction.GetAccount()
	if err != nil {
		resp := u.Message(false, "Card Not Found")
		return resp
	}

	if account.AvailableFunds() >= transaction.Amount {
		account.AdjustBlocked(transaction.Amount)
	} else {
		resp := u.Message(false, "insufficient funds")
		return resp
	}

	GetDB().Create(transaction)
	// GetDB().Save(&account)

	resp := u.Message(true, "success")
	resp["transaction"] = transaction
	return resp
}

func (transaction *Transaction) AuthorisedUncaptured() int {
	return transaction.Amount - transaction.Captured
}

func (transaction *Transaction) CaptureFunds(amount int, transactionId int) map[string]interface{} {

	transaction.GetTransaction(transactionId)

	if transaction.AuthorisedUncaptured() >= amount && transaction.AuthorisedUncaptured() > 0 {
		transaction.Captured += amount
		account, _ := transaction.GetAccount()
		account.AdjustBlocked(amount - amount*2)
		account.AdjustFunds(amount - amount*2)
	} else {
		resp := u.Message(false, "Not enough authorised funds")
		return resp
	}

	GetDB().Save(&transaction)
	return nil
}

func (transaction *Transaction) ChangeAuthorisedFunds(amount int, transactionId int) map[string]interface{} {

	transaction.GetTransaction(transactionId)

	if transaction.AuthorisedUncaptured() >= amount && transaction.AuthorisedUncaptured() > 0 {
		transaction.Amount -= amount
		GetDB().Save(&transaction)

		account, _ := transaction.GetAccount()
		account.AdjustBlocked(amount - amount*2)
		return nil
	} else {
		resp := u.Message(false, "Not enough authorised funds")
		return resp
	}

}

func (transaction *Transaction) Refund(amount int) {
	transaction.Refunded += amount
	GetDB().Save(&transaction)
	return
}
