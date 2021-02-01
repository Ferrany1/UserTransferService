package netHandler

import (
	"UserTransferService/src/userService"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)
// Creates new user
func CreateUser(c *gin.Context) {
	var uInf = new(userService.UserInfo)
	if err := c.Bind(&uInf); err != nil || uInf.UserData.Username == "" || uInf.UserData.Email == "" || uInf.UserData.Password == "" {
		c.JSON(http.StatusBadRequest, map[string]string{"message": "set all required fields"})
		return
	}

	if err := uInf.CreateNewUser(); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, map[string]string{"status": "success"})
	}
}

// Gets users balance
func GetUserBalance(c *gin.Context) {
	var uInf = new(userService.UserInfo)
	if err := c.BindJSON(&uInf); err != nil || uInf.UserData.Username == "" || uInf.UserData.Email == "" || uInf.UserData.Password == "" {
		log.Println(uInf)
		c.JSON(http.StatusBadRequest, map[string]string{"message": "set all required fields"})
		return
	}

	if balance, err := uInf.GetUserBalance(); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, map[string]int{"balance": balance.Balance})
	}
}

// Makes a balance transfer between two users
func TransferFromBalance(c *gin.Context) {
	var (
		uInf     = new(userService.UserInfo)
		receiver userService.UserInfo
		tr       int
		err      error
	)
	if err = c.Bind(&uInf); err != nil || uInf.UserData.Username == "" || uInf.UserData.Email == "" || uInf.UserData.Password == "" {
		c.JSON(http.StatusBadRequest, map[string]string{"message": "set all required fields"})
		return
	}

	if rec := c.Param("receiver"); rec != "" {
		receiver.UserData.Username = rec
		receiver.UserBalance.Username = rec
	} else {
		c.JSON(http.StatusBadRequest, map[string]string{"message": "set receiver"})
		return
	}

	if amount := c.Param("amount"); amount != "" {
		if tr, err = strconv.Atoi(amount); err != nil {

			c.JSON(http.StatusBadRequest, map[string]string{"message": "wrong transaction amount value"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, map[string]string{"message": "set transaction amount"})
		return
	}

	if err = uInf.TransferUsersBalance(receiver, tr); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, map[string]string{"status": "success"})
	}
}

// Updates users creds
func UpdateUser(c *gin.Context) {
	var (
		uInf    = new(userService.UserInfo)
		newuInf = new(userService.UserInfo)
	)
	if err := c.BindJSON(&uInf); err != nil || uInf.UserData.Username == "" || uInf.UserData.Email == "" || uInf.UserData.Password == "" {
		c.JSON(http.StatusBadRequest, map[string]string{"message": "set all required fields"})
		return
	}
	newuInf.UserData.Username = uInf.UserData.Username
	newuInf.UserData.Email 		= c.Query("email")
	newuInf.UserData.Password 	= c.Query("password")

	if (newuInf.UserData.Email == "" && newuInf.UserData.Password == "") || (newuInf.UserData.Email == uInf.UserData.Email && newuInf.UserData.Password == uInf.UserData.Password) {
		c.JSON(http.StatusBadRequest, map[string]string{"message": "nothing to change"})
		return
	}

	if err := uInf.UpdateUser(*newuInf); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, map[string]string{"status": "success"})
	}
}

// Deletes user
func DeleteUser(c *gin.Context) {
	var uInf = new(userService.UserInfo)
	if err := c.Bind(&uInf); err != nil || uInf.UserData.Username == "" || uInf.UserData.Email == "" || uInf.UserData.Password == "" {
		c.JSON(http.StatusBadRequest, map[string]string{"message": "set all required fields"})
		return
	}

	if err := uInf.DeleteUser(); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, map[string]string{"status": "success"})
	}
}