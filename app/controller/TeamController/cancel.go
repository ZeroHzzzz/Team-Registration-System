package teamcontroller

import (
	"backend/app/server"
	"backend/app/utils"

	"github.com/gin-gonic/gin"
)

type CancelData struct {
	TeamID string `json:"teamid" binding:"required"`
}

func Cancel(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	var canceldata CancelData
	err := c.ShouldBindJSON(&canceldata)
	if err != nil {
		utils.ResponseError(c, utils.ParameterErrorCode, utils.ParameterErrorMsg)
		return
	}

	// 鉴权
	team, err := server.GetTeam(canceldata.TeamID)
	if err == utils.ErrTeamNotFound {
		utils.ResponseError(c, utils.NotFoundCode, utils.NotFoundMsg)
		return
	}

	userEmail, existsEmail := c.Get("email")
	userType, existsType := c.Get("type")

	if !existsEmail || !existsType {
		utils.ResponseUnauthorized(c)
		return
	} else if userEmail != team.LeaderID || userType == 0 {
		utils.ResponseUnauthorized(c)
		return
	}

	err = server.Cancel(canceldata.TeamID)
	if err != nil {
		if err == utils.ErrTeamNotFound {
			utils.ResponseError(c, utils.NotFoundCode, utils.NotFoundMsg)
			return
		}
		utils.ResponseError(c, utils.OperationFailedCode, utils.OperationFailedMsg)
		return
	}
	utils.ResponseSuccess(c, nil)
	return
}
