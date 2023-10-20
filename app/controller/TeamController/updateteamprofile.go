package teamcontroller

import (
	"backend/app/server"
	"backend/app/utils"

	"github.com/gin-gonic/gin"
)

type UpdateTeamProfileData struct {
	TeamID      string `json:"teamid" binding:"required"`
	TeamName    string `json:"teamname"`
	Description string `json:"description"`
}

func UpdateTeamProfile(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	var updateprofile UpdateTeamProfileData
	err := c.ShouldBindJSON(&updateprofile)
	if err != nil {
		utils.ResponseError(c, utils.ParameterErrorCode, utils.ParameterErrorMsg)
		return
	}

	// 鉴权
	team, err := server.GetTeam(updateprofile.TeamID)
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

	err = server.UpdateTeamProfile(updateprofile.TeamID, updateprofile)
	if err != nil {
		if err == utils.ErrCopyFail {
			utils.ResponseError(c, utils.OperationFailedCode, utils.CopyFailMsg)
			return
		} else if err == utils.ErrTeamNotFound {
			utils.ResponseError(c, utils.NotFoundCode, utils.TeamNotFoundMsg)
			return
		}
		utils.ResponseError(c, utils.OperationFailedCode, utils.OperationFailedMsg)
		return
	}
	utils.ResponseSuccess(c, nil)
	return
}
