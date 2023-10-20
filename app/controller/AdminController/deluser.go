package admincontroller

import (
	"backend/app/server"
	"backend/app/utils"

	"github.com/gin-gonic/gin"
)

type DeleteUserData_Admin struct {
	Email string `json:"email" binding:"required"`
}

func DelUser_Admin(c *gin.Context) {
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

	var deleteuserdata DeleteUserData_Admin
	err := c.ShouldBindJSON(&deleteuserdata)
	if err != nil {
		utils.ResponseError(c, utils.ParameterErrorCode, utils.ParameterErrorMsg)
		return
	}
	err = server.DeleteUser_Admin(deleteuserdata.Email)
	if err != nil {
		// fmt.Println(err)
		utils.ResponseError(c, utils.OperationFailedCode, utils.OperationFailedMsg)
		return
	} else {
		utils.ResponseSuccess(c, nil)
		return
	}
}
