package cliHandler

import (
	"UserTransferService/src/userService"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type consoleArgs struct {
	Args []string
}

// Checks which command was called
func CommandCheck(text string){
	var (
		ca      consoleArgs
		command bool
	)
	ca.Args = strings.Split(text, " ")
	if len(ca.Args) < 3 {
		log.Printf("> ERROR: %s", "not enough arguments")
		return
	}
	for _, arg := range ca.Args {
		switch arg {
		case "create":
			ca.createUser()
			command = true
			break
		case "balance":
			ca.getUserBalance()
			command = true
			break
		case "transfer":
			ca.transferFromBalance()
			command = true
			break
		case "update":
			ca.updateUser()
			command = true
			break
		case "delete":
			ca.deleteUser()
			command = true
			break
		}
	}
	if !command {
		fmt.Printf("> ERROR: %s", "no command provided")
	}
}

// Checks for arguments in query
func (ca consoleArgs) argsChecks() (uInf userService.UserInfo, rinf userService.UserInfo, tr int, updStr [2]string, err error) {
	for _, arg := range ca.Args {
		if strings.Contains(arg, "username") {
			subst := strings.Split(arg, ":")
			if len(subst) == 2 {
				uInf.UserData.Username = subst[1]
				uInf.UserBalance.Username = subst[1]
			}
		}
		if strings.Contains(arg, "email") && !strings.Contains(arg, "new_email") {
			subst := strings.Split(arg, ":")
			if len(subst) == 2 {
				uInf.UserData.Email = subst[1]
			}
		}
		if strings.Contains(arg, "password") && !strings.Contains(arg, "new_password") {
			subst := strings.Split(arg, ":")
			if len(subst) == 2 {
				uInf.UserData.Password = subst[1]
			}
		}
		if strings.Contains(arg, "balance") {
			subst := strings.Split(arg, ":")
			if len(subst) == 2 {
				if uInf.UserBalance.Balance, err = strconv.Atoi(subst[1]); err != nil {
					return uInf, rinf, tr, updStr, fmt.Errorf("wrong balance type, should be int")
				}
			}
		}
		if strings.Contains(arg, "receiver") {
			subst := strings.Split(arg, ":")
			if len(subst) == 2 {
				rinf.UserData.Username = subst[1]
				rinf.UserBalance.Username = subst[1]
			}
		}
		if strings.Contains(arg, "amount") {
			subst := strings.Split(arg, ":")
			if len(subst) == 2 {
				if tr, err = strconv.Atoi(subst[1]); err != nil {
					return uInf, rinf, tr, updStr, fmt.Errorf("wrong amount type, should be int")
				}
			}
		}
		if strings.Contains(arg, "new_email") {
			subst := strings.Split(arg, ":")
			if len(subst) == 2 {
				updStr[0] = subst[1]
			}
		}
		if strings.Contains(arg, "new_password") {
			subst := strings.Split(arg, ":")
			if len(subst) == 2 {
				updStr[1] = subst[1]
			}
		}
	}
	if uInf.UserData.Username == "" && uInf.UserData.Email == "" && uInf.UserData.Password == "" {
		return uInf, rinf, tr, updStr, fmt.Errorf("set all required fields")
	}
	return
}

// Creates new user
func (ca consoleArgs) createUser() {
	uInf, _, _, _, err := ca.argsChecks()
	if err != nil {
		fmt.Printf("> ERROR: %s\n", "check log file")
		return
	}

	if err = uInf.CreateNewUser(); err != nil {
		fmt.Printf("> ERROR: %s\n", "check log file")
		return
	} else {
		fmt.Printf("> Success: %s\n", "created new user")
	}
}

// Gets users balance
func (ca consoleArgs) getUserBalance() {
	uInf, _, _, _, err := ca.argsChecks()
	if err != nil {
		fmt.Printf("> ERROR: %s\n", "check log file")
		return
	}

	if balance, err := uInf.GetUserBalance(); err != nil {
		fmt.Printf("> ERROR: %s\n", "check log file")
		return
	} else {
		fmt.Printf("> Balance: %d\n", balance.Balance)
	}
}

// Makes a balance transfer between two users
func (ca consoleArgs) transferFromBalance() {
	uInf, rinf, tr, _, err := ca.argsChecks()
	if err != nil {
		fmt.Printf("> ERROR: %s\n", "check log file")
		return
	}
	if rinf.UserData.Username == "" {
		fmt.Printf("> ERROR: %s\n", "set receiver")
	}
	if tr == 0 {
		fmt.Printf("> ERROR: %s\n", "set transaction amount")
	}

	if err = uInf.TransferUsersBalance(rinf, tr); err != nil {
		fmt.Printf("> ERROR: %s\n", "check log file")
		return
	} else {
		fmt.Printf("> Success: %s\n", "made a transaction")
	}
}

// Updates users creds
func (ca consoleArgs) updateUser() {
	var newuInf = new(userService.UserInfo)
	uInf, _, _, newStr, err := ca.argsChecks()
	if err != nil {
		fmt.Printf("> ERROR: %s\n", "check log file")
		return
	}

	if newStr[0] == "" || newStr[1] == "" {
		fmt.Printf("> ERROR: %s\n", "set new_email or new_password")
		return
	}

	newuInf.UserData.Username = uInf.UserData.Username
	if newStr[0] != "" {
		newuInf.UserData.Email = newStr[0]
	}
	if newStr[1] != "" {
		newuInf.UserData.Password = newStr[1]
	}

	if (newuInf.UserData.Email == "" && newuInf.UserData.Password == "") || (newuInf.UserData.Email == uInf.UserData.Email && newuInf.UserData.Password == uInf.UserData.Password) {
		fmt.Printf("> ERROR: %s\n", "nothing to change")
		return
	}

	if err := uInf.UpdateUser(*newuInf); err != nil {
		fmt.Printf("> ERROR: %s\n", "check log file")
		return
	} else {
		fmt.Printf("> Success: %s\n", "updated user info")
	}

}

// Deletes user
func (ca consoleArgs) deleteUser() {
	uInf, _, _, _, err := ca.argsChecks()
	if err != nil {
		fmt.Printf("> ERROR: %s\n", "check log file")
		return
	}

	if err := uInf.DeleteUser(); err != nil {
		fmt.Printf("> ERROR: %s\n", "check log file")
		return
	} else {
		fmt.Printf("> Success: %s\n", "deleted user")
	}
}