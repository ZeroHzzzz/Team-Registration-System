package teamcontroller

import (
	"backend/app/server"
	"backend/app/utils"

	"github.com/gin-gonic/gin"
)

type AddMemberData struct {
	Email        string `json:"email" binding:"required"`
	TeamID       string `json:"teamid" binding:"required"`
	TeamPassword string `json:"password" binding:"required"`
}

func AddMember(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	var addmemberdata AddMemberData
	err := c.ShouldBindJSON(&addmemberdata)
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
	} else if userEmail != addmemberdata.Email && userType == 0 {
		utils.ResponseUnauthorized(c)
		return
	}

	err = server.AddMember(addmemberdata.Email, addmemberdata.TeamID, addmemberdata.TeamPassword)
	if err != nil {
		if err == utils.ErrFormatWrong {
			utils.ResponseError(c, utils.ParameterErrorCode, utils.FormatWrongMsg)
			return
		} else if err == utils.ErrUserNotFound {
			utils.ResponseError(c, utils.NotFoundCode, utils.UserNotFoundMsg)
			return
		} else if err == utils.ErrTeamNotFound {
			utils.ResponseError(c, utils.NotFoundCode, utils.TeamNotFoundMsg)
			return
		} else if err == utils.ErrCreateMsgFailed {
			utils.ResponseError(c, utils.OperationFailedCode, utils.CreateMsgFailedMsg)
			return
		}
		// fmt.Println(err)
		utils.ResponseError(c, utils.OperationFailedCode, utils.AddMemberFailMsg)
		return
	}
	utils.ResponseSuccess(c, nil)
	return
}
