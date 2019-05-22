package models

import (
	u "go-card/utils"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

//a struct to rep user account
type Account struct {
	gorm.Model
	CardNo   string `json:"cardNo"`
	Password string `json:"password"`
	Funds    int    `json:"funds";sql:"DEFAULT:0"`
	Blocked  int    `json:"blocked";sql:"DEFAULT:0"`
}

type Balance struct {
	Funds     int `json:"funds"`
	Blocked   int `json:"blocked"`
	Available int `json:"available"`
}

//Validate incoming user details...
func (account *Account) Validate() (map[string]interface{}, bool) {

	if len(account.Password) < 6 {
		return u.Message(false, "Password is required"), false
	}

	//CardNo must be unique
	temp := &Account{}

	//check for errors and duplicate CardNo
	err := GetDB().Table("accounts").Where("card_no = ?", account.CardNo).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry"), false
	}
	if temp.CardNo != "" {
		return u.Message(false, "CardNo already in use by another user."), false
	}

	return u.Message(false, "Requirement passed"), true
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {

		return false
	}
	return true
}

func (account *Account) Create() map[string]interface{} {

	if resp, ok := account.Validate(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	GetDB().Create(account)

	if account.ID <= 0 {
		return u.Message(false, "Failed to create account, connection error.")
	}

	account.Password = "" //delete password

	response := u.Message(true, "Account has been created")
	response["account"] = account
	return response
}

func (account *Account) GetBalance() map[string]interface{} {
	password := account.Password
	cardNum := account.CardNo

	err := GetDB().Table("accounts").Where("card_no = ?", cardNum).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "cardNo not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	if CheckPasswordHash(password, account.Password) == true {
		var balance Balance
		balance.Funds = account.Funds
		balance.Blocked = account.Blocked
		balance.Available = account.AvailableFunds()

		response := u.Message(true, "Account Balance")
		response["balance"] = balance
		return response
	}

	response := u.Message(true, "Auth Error")
	return response
}

func (account *Account) AddFunds() map[string]interface{} {

	cardNum := account.CardNo
	toAdd := account.Funds

	err := GetDB().Table("accounts").Where("card_no = ?", cardNum).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "cardNo not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}
	account.Funds += toAdd
	GetDB().Save(&account)

	response := u.Message(true, "Account Funds Updated")
	return response
}

func (account *Account) AvailableFunds() int {
	response := account.Funds - account.Blocked
	return response
}

func (account *Account) AdjustBlocked(amount int) {
	account.Blocked += amount
	GetDB().Save(&account)
	return
}

func (account *Account) AdjustFunds(amount int) {
	account.Funds += amount
	GetDB().Save(&account)
	return
}
