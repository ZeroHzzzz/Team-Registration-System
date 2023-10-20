package usercontroller

import (
	"backend/app/server"
	"backend/app/utils"

	"github.com/gin-gonic/gin"
)

type GetUserData struct {
	Email string `form:"email" binding:"required"`
}

func GetUser(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	// 鉴权
	_, existsEmail := c.Get("email")
	_, existsType := c.Get("type")

	if !existsEmail || !existsType {
		utils.ResponseUnauthorized(c)
		return
	}

	var getuserdata GetUserData
	err := c.ShouldBindQuery(&getuserdata)
	if err != nil {
		utils.ResponseError(c, utils.ParameterErrorCode, utils.ParameterErrorMsg)
		return
	}
	user, err := server.GetUser(getuserdata.Email)
	if err != nil {
		if err == utils.ErrUserNotFound {
			utils.ResponseError(c, utils.NotFoundCode, utils.NotFoundMsg)
			return
		}
		utils.ResponseError(c, utils.OperationFailedCode, utils.OperationFailedMsg)
		return
	} else {
		utils.ResponseSuccess(c, user)
		return
	}
}
