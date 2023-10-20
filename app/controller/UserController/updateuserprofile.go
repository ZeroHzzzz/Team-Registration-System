package usercontroller

import (
	"backend/app/server"
	"backend/app/utils"

	"github.com/gin-gonic/gin"
)

type UpdateUserProfileData struct {
	Email       string `json:"email" binding:"required"`
	Username    string `json:"username"`
	Sign        string `json:"sign"`
	Description string `json:"description"`
	TelePhone   int    `json:"telephone"`
}

func UpdateUserProfile(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	var updateprofile UpdateUserProfileData
	err := c.ShouldBindJSON(&updateprofile)
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
	} else if userEmail != updateprofile.Email && userType == 0 {
		utils.ResponseUnauthorized(c)
		return
	}

	err = server.UpdateUserProfile(updateprofile.Email, updateprofile)
	if err != nil {
		if err == utils.ErrCopyFail {
			utils.ResponseError(c, utils.OperationFailedCode, utils.CopyFailMsg)
			return
		} else if err == utils.ErrUserNotFound {
			utils.ResponseError(c, utils.NotFoundCode, utils.UserNotFoundMsg)
			return
		}
		utils.ResponseError(c, utils.OperationFailedCode, utils.OperationFailedMsg)
		return
	}
	utils.ResponseSuccess(c, nil)
	return
}
