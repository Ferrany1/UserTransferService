package router

import (
	"UserTransferService/src/userService/netHandler"
	"UserTransferService/src/system/config"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

var router *gin.Engine

// Launces router
func LaunchRouter() {
	createRouter()
	log.Println(router.Run(":" + strconv.Itoa(config.CF.Router.Port)))
}

// Creates router instance
func createRouter() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	v1 := r.Group("/api/v1")

	u := v1.Group("/user")
		u.POST("/create", netHandler.CreateUser)
		u.GET("/balance", netHandler.GetUserBalance)
		u.POST("/transfer/:receiver/:amount", netHandler.TransferFromBalance)
		u.PUT("/update", netHandler.UpdateUser)
		u.DELETE("/delete", netHandler.DeleteUser)

	router = r
}