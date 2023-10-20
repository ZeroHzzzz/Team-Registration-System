package teamcontroller

import (
	"backend/app/model"
	"backend/app/server"
	"backend/app/utils"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GetTeamData struct {
	TeamID string `form:"teamid"`
}

func GetTeam(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	var geteamdata GetTeamData
	err := c.ShouldBindQuery(&geteamdata)
	if err != nil {
		utils.ResponseError(c, utils.ParameterErrorCode, utils.ParameterErrorMsg)
		return
	}
	Team, err := server.GetTeam(geteamdata.TeamID)
	// 可能会出bug
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ResponseError(c, utils.NotFoundCode, utils.NotFoundMsg)
			return
		}
		utils.ResponseError(c, utils.OperationFailedCode, utils.OperationFailedMsg)
		return
	}

	utils.ResponseSuccess(c, Team)
	return
}

func GetTeamMember(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	var users []model.Usermodel
	var geteammemberdata GetTeamData
	err1 := c.ShouldBindQuery(&geteammemberdata)
	if err1 != nil {
		utils.ResponseError(c, utils.ParameterErrorCode, utils.ParameterErrorMsg)
		return
	}
	users, err1 = server.GetTeamMember(geteammemberdata.TeamID)
	if err1 != nil {
		utils.ResponseError(c, utils.OperationFailedCode, utils.OperationFailedMsg)
		return
	} else {
		utils.ResponseSuccess(c, users)
		return
	}
}
