package userService

import (
	"UserTransferService/src/db_service"
	"fmt"
)

type UserInfo struct {
	UserData    db_service.User    `json:"user_data"`
	UserBalance db_service.Balance `json:"user_balance"`
}

// Creates new user in db
func (uInf UserInfo) CreateNewUser() (err error) {
	return db_service.DB.CreateUser(uInf.UserData)
}

// Get user's balance in db
func (uInf UserInfo) GetUserBalance() (balance db_service.Balance, err error) {
	ver, err := uInf.verifyUser()
	if !ver {
		return balance, err
	}

	balance, err = db_service.DB.GetBalance(uInf.UserData)
	if err != nil {
		return balance, fmt.Errorf("failed to get user balance %s", err.Error())
	}
	return
}

// Transfers money from sender to receiver in db
func (uInf UserInfo) TransferUsersBalance(receiver UserInfo, tr int) (err error) {
	ver, err := uInf.verifyUser()
	if !ver {
		return err
	}

	_, err = db_service.DB.GetUser(receiver.UserData)
	if err != nil {
		return fmt.Errorf("failed to get reciver info %s", err.Error())
	}

	sBalance, err := db_service.DB.GetBalance(uInf.UserData)
	if err != nil {
		return fmt.Errorf("failed to get user balance %s", err.Error())
	}

	rBalance, err := db_service.DB.GetBalance(receiver.UserData)
	if err != nil {
		return fmt.Errorf("failed to get reciever balance %s", err.Error())
	}

	if tr > sBalance.Balance {
		return fmt.Errorf("transaction amount exceeds balance")
	}

	err = db_service.DB.UpdateBalance(db_service.Balance{Username: uInf.UserBalance.Username, Balance: sBalance.Balance - tr}, db_service.Balance{Username: receiver.UserBalance.Username, Balance: rBalance.Balance + tr})
	if err != nil {
		return fmt.Errorf("failed to update user balance %s", err.Error())
	}
	return
}

// Updates user creds
func (uInf UserInfo) UpdateUser(newUInf UserInfo) (err error) {
	ver, err := uInf.verifyUser()
	if !ver {
		return err
	}

	err = db_service.DB.UpdateUser(newUInf.UserData)
	if err != nil {
		return fmt.Errorf("failed to update user info %s", err.Error())
	}
	return
}

// Delets user from db
func (uInf UserInfo) DeleteUser() (err error) {
	ver, err := uInf.verifyUser()
	if !ver {
		return err
	}

	err = db_service.DB.DeleteUser(uInf.UserData)
	if err != nil {
		return fmt.Errorf("failed to update user info %s", err.Error())
	}
	return
}

// Verifies user in db
func (uInf UserInfo) verifyUser() (vf bool, err error) {
	user, err := db_service.DB.GetUser(uInf.UserData)
	if err != nil {
		return false, fmt.Errorf("failed to get user info %s", err.Error())
	}

	if user != uInf.UserData {
		return false, fmt.Errorf("failed to verify user %s", "wrong credentials")
	} else {
		return true, err
	}
}