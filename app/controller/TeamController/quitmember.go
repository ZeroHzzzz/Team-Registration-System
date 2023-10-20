package teamcontroller

import (
	"backend/app/server"
	"backend/app/utils"

	"github.com/gin-gonic/gin"
)

type QuitMemberData struct {
	Email  string `json:"email" binding:"required"`
	TeamID string `json:"teamid" binding:"required"`
}

func QuitMember(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	var quitmemberdata QuitMemberData
	err := c.ShouldBindJSON(&quitmemberdata)
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
	} else if userEmail != quitmemberdata.Email && userType == 0 {
		utils.ResponseUnauthorized(c)
		return
	}

	err = server.QuitTeam(quitmemberdata.Email, quitmemberdata.TeamID)
	if err != nil {
		if err == utils.ErrFormatWrong {
			utils.ResponseError(c, utils.ParameterErrorCode, utils.FormatWrongMsg)
		}
		utils.ResponseError(c, utils.OperationFailedCode, utils.OperationFailedMsg)
		return
	}
	utils.ResponseSuccess(c, nil)
	return
}
