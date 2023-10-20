package teamcontroller

import (
	"backend/app/server"
	"backend/app/utils"

	"github.com/gin-gonic/gin"
)

type CreateTeamData struct {
	TeamName        string `json:"teamname" binding:"required"`
	TeamPassword    string `json:"password" binding:"required"`
	TeamDescription string `json:"teamdescription"`
	LeaderID        string `json:"leaderid" binding:"required"` // email
}

func CreateTeam(c *gin.Context) {
	// 鉴权
	_, exists := c.Get("email")
	if !exists {
		utils.ResponseUnauthorized(c)
		return
	}
	_, exists = c.Get("type")
	if !exists {
		utils.ResponseUnauthorized(c)
		return
	}

	var createteamdata CreateTeamData
	err := c.ShouldBindJSON(&createteamdata)
	if err != nil {
		utils.ResponseError(c, utils.ParameterErrorCode, utils.ParameterErrorMsg)
		return
	}
	newteamid, err := server.TeamRegister(createteamdata)
	if err != nil || newteamid == "" {
		if err == utils.ErrCopyFail {
			utils.ResponseError(c, utils.OperationFailedCode, utils.CopyFailMsg)
			return
		} else if err == utils.ErrFormatWrong {
			utils.ResponseError(c, utils.ParameterErrorCode, utils.FormatWrongMsg)
			return
		} else if err == utils.ErrUserNotFound {
			utils.ResponseError(c, utils.NotFoundCode, utils.UserNotFoundMsg)
			return
		} else if err == utils.ErrHaveInTeam {
			utils.ResponseError(c, utils.HaveExistCode, utils.HaveInTeamMsg)
			return
		} else if err == utils.ErrAddMemberOperationFailed {
			utils.ResponseError(c, utils.OperationFailedCode, utils.AddMemberFailMsg)
			return
		} else if err == utils.ErrOperationFailed {
			utils.ResponseError(c, utils.OperationFailedCode, utils.OperationFailedMsg)
			return
		}
		utils.ResponseError(c, utils.OperationFailedCode, utils.OperationFailedMsg)
		return
	}
	utils.ResponseSuccess(c, newteamid)
	c.Set("type", 1)
	return
}
