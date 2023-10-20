package teamcontroller

import (
	"backend/app/server"
	"backend/app/utils"

	"github.com/gin-gonic/gin"
)

type SubmitData struct {
	TeamID       string `json:"teamid" binding:"required"`
	TeamPassword string `json:"password" binding:"required"`
}

func Submit(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	var submitdata SubmitData
	err := c.ShouldBindJSON(&submitdata)
	if err != nil {
		utils.ResponseError(c, utils.ParameterErrorCode, utils.ParameterErrorMsg)
		return
	}

	// 鉴权
	team, err := server.GetTeam(submitdata.TeamID)
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

	err = server.Submit(submitdata.TeamID, submitdata.TeamPassword)
	if err != nil {
		if err == utils.ErrTeamNotFound {
			utils.ResponseError(c, utils.NotFoundCode, utils.NotFoundMsg)
			return
		} else if err == utils.ErrSubmitFailed {
			utils.ResponseError(c, utils.OperationFailedCode, utils.SubmitFailedMsg)
			return
		}
		utils.ResponseError(c, utils.OperationFailedCode, utils.OperationFailedMsg)
		return
	}
	utils.ResponseSuccess(c, nil)
	return
}
