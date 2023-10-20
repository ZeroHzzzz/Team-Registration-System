package teamcontroller

import (
	"backend/app/server"
	"backend/app/utils"

	"github.com/gin-gonic/gin"
)

func UpLoadTeamAvatara(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	TeamID := c.PostForm("TeamID")
	file, err := c.FormFile("Avatar")
	if err != nil {
		utils.ResponseError(c, utils.ParameterErrorCode, utils.ParameterErrorMsg)
		return
	}
	// 鉴权
	team, err := server.GetTeam(TeamID)
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

	AvataraUrl := "UserUploads/" + file.Filename
	if err := c.SaveUploadedFile(file, AvataraUrl); err != nil {
		utils.ResponseError(c, utils.OperationFailedCode, utils.FileSaveFailedMsg)
		return
	}
	err = server.UpLoadTeamAvatara(TeamID, AvataraUrl)
	if err != nil {
		if err == utils.ErrTeamNotFound {
			utils.ResponseError(c, utils.NotFoundCode, utils.TeamNotFoundMsg)
			return
		}
		// fmt.Println(err)
		utils.ResponseError(c, utils.OperationFailedCode, utils.OperationFailedMsg)
		return
	}
	AvataraUrl = "http://127.0.0.1:8080/" + AvataraUrl
	utils.ResponseSuccess(c, AvataraUrl)
	return
}
