package usercontroller

import (
	"backend/app/server"
	"backend/app/utils"

	"github.com/gin-gonic/gin"
)

func UpLoadAvatara(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	Email := c.PostForm("Email")
	file, err := c.FormFile("Avatar")
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
	} else if userEmail != Email && userType == 0 {
		utils.ResponseUnauthorized(c)
		return
	}

	AvataraUrl := "UserUploads/" + file.Filename
	if err := c.SaveUploadedFile(file, AvataraUrl); err != nil {
		utils.ResponseError(c, utils.OperationFailedCode, utils.FileSaveFailedMsg)
		return
	}
	err = server.UpLoadAvatara(Email, AvataraUrl)
	if err != nil {
		if err == utils.ErrUserNotFound {
			utils.ResponseError(c, utils.NotFoundCode, utils.UserNotFoundMsg)
			return
		}
		utils.ResponseError(c, utils.OperationFailedCode, utils.OperationFailedMsg)
		return
	}
	AvataraUrl = "http://127.0.0.1:8080/" + AvataraUrl
	utils.ResponseSuccess(c, AvataraUrl)
	return
}
