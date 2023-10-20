package usercontroller

import (
	"backend/app/server"
	"backend/app/utils"

	"github.com/gin-gonic/gin"
)

// 解析信息
type RegisterData struct {
	Email    string `jion:"email" binding:"required"`
	Username string `json:"username"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	var registerdata RegisterData
	err := c.ShouldBindJSON(&registerdata)
	if err != nil {
		utils.ResponseError(c, utils.ParameterErrorCode, utils.ParameterErrorMsg)
		return
	}

	err = server.Register(registerdata)
	if err != nil {
		if err == utils.ErrCopyFail {
			utils.ResponseError(c, utils.OperationFailedCode, utils.CopyFailMsg)
			return
		} else if err == utils.ErrFormatWrong {
			utils.ResponseError(c, utils.ParameterErrorCode, utils.FormatWrongMsg)
			return
		} else if err == utils.ErrUserHaveExist {
			utils.ResponseError(c, utils.HaveExistCode, utils.UserHaveExistMsg)
			return
		} else if err == utils.ErrOperationFailed {
			utils.ResponseError(c, utils.OperationFailedCode, utils.OperationFailedMsg)
			return
		}
		utils.ResponseError(c, utils.OperationFailedCode, utils.OperationFailedMsg)
		return
	}
	utils.ResponseSuccess(c, nil)
	return
}
