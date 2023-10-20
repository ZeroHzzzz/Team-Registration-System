package admincontroller

import (
	"backend/app/server"
	"backend/app/utils"

	"github.com/gin-gonic/gin"
)

type DelTeamData_Admin struct {
	TeamID string `json:"teamid" binding:"required"`
}

func DelTeam_Admin(c *gin.Context) {
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

	var delteamdata DelTeamData_Admin
	err := c.ShouldBindJSON(&delteamdata)
	if err != nil {
		utils.ResponseError(c, utils.ParameterErrorCode, utils.ParameterErrorMsg)
		return
	}
	err = server.DelTeam_Admin(delteamdata.TeamID)
	if err != nil {
		// fmt.Println(err)
		if err == utils.ErrTeamNotFound {
			utils.ResponseError(c, utils.NotFoundCode, utils.NotFoundMsg)
			return
		} else if err == utils.ErrDelTeamFailed {
			utils.ResponseError(c, utils.OperationFailedCode, utils.DelTeamFailedMsg)
			return
		} else if err == utils.ErrCreateMsgFailed {
			utils.ResponseError(c, utils.OperationFailedCode, utils.CreateMsgFailedMsg)
			return
		}
		utils.ResponseError(c, utils.OperationFailedCode, utils.OperationFailedMsg)
		return
	}
	utils.ResponseSuccess(c, nil)
	return
}
