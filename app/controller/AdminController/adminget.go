package admincontroller

import (
	"backend/app/server"
	"backend/app/utils"

	"github.com/gin-gonic/gin"
)

func GetAllUser(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	// 鉴权
	_, existsEmail := c.Get("email")
	userType, existsType := c.Get("type")

	if !existsEmail || !existsType {
		utils.ResponseUnauthorized(c)
		return
	} else if userType == 2 {
		utils.ResponseUnauthorized(c)
		return
	}

	users, err := server.GetAllUser()
	if err != nil {
		utils.ResponseError(c, utils.OperationFailedCode, utils.OperationFailedMsg)
		return
	}
	utils.ResponseSuccess(c, users)
	return
}

func GetAllTeam(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	// 鉴权
	_, existsEmail := c.Get("email")
	userType, existsType := c.Get("type")

	if !existsEmail || !existsType {
		utils.ResponseUnauthorized(c)
		return
	} else if userType == 2 {
		utils.ResponseUnauthorized(c)
		return
	}

	teams, err := server.GetAllTeam()
	if err != nil {
		utils.ResponseError(c, utils.OperationFailedCode, utils.OperationFailedMsg)
		return
	}
	utils.ResponseSuccess(c, teams)
	return
}
