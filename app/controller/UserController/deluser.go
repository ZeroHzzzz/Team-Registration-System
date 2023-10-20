package usercontroller

import (
	"backend/app/server"
	"backend/app/utils"

	"github.com/gin-gonic/gin"
)

type DeleteUserData struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func DeleteUser(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	var deleteuserdata DeleteUserData
	err := c.ShouldBindJSON(&deleteuserdata)
	if err != nil {
		utils.ResponseError(c, utils.ParameterErrorCode, utils.ParameterErrorMsg)
		return
	}

	// 鉴权
	userEmail, existsEmail := c.Get("email")
	userType, existsType := c.Get("type")

	if !existsEmail || !existsType {
		utils.ResponseUnauthorized(c)
		return
	} else if userEmail != deleteuserdata.Email && userType == 0 {
		utils.ResponseUnauthorized(c)
		return
	}

	err = server.DeleteUser(deleteuserdata.Email, deleteuserdata.Password)
	if err != nil {
		utils.ResponseError(c, utils.OperationFailedCode, utils.OperationFailedMsg)
		return
	} else {
		utils.ResponseSuccess(c, nil)
		return
	}
}
