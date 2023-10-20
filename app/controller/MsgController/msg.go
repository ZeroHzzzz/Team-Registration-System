package msgcontroller

import (
	"backend/app/server"
	"backend/app/utils"

	"github.com/gin-gonic/gin"
)

type UpdateMsgData struct {
	Email string `json:"email" binding:"required"`
}

func UpdateUnread(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	var updatemsgdata UpdateMsgData
	err := c.ShouldBindJSON(&updatemsgdata)
	if err != nil {
		utils.ResponseError(c, utils.ParameterErrorCode, utils.ParameterErrorMsg)
		return
	}
	err = server.UpdateUnread(updatemsgdata.Email)
	if err != nil {
		if err == utils.ErrUserNotFound {
			utils.ResponseError(c, utils.OperationFailedCode, utils.NotFoundMsg)
			return
		} else if err == utils.ErrFormatWrong {
			utils.ResponseError(c, utils.ParameterErrorCode, utils.FormatWrongMsg)
			return
		}
		utils.ResponseError(c, utils.OperationFailedCode, utils.OperationFailedMsg)
		return
	}
	utils.ResponseSuccess(c, nil)
	return
}

type GetMsgData struct {
	Email string `form:"Email"`
}

func GetMsg(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	var updatemsgdata GetMsgData
	err := c.ShouldBindQuery(&updatemsgdata)
	if err != nil {
		utils.ResponseError(c, utils.ParameterErrorCode, utils.ParameterErrorMsg)
		return
	}
	// fmt.Println(updatemsgdata.Email)
	msgs, err := server.GetMsg(updatemsgdata.Email)
	if err != nil {
		if err == utils.ErrFormatWrong {
			utils.ResponseError(c, utils.ParameterErrorCode, utils.FormatWrongMsg)
			return
		}
		utils.ResponseError(c, utils.OperationFailedCode, utils.OperationFailedMsg)
		return
	}

	utils.ResponseSuccess(c, msgs)
	return
}
