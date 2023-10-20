package admincontroller

import (
	"backend/app/server"
	"backend/app/utils"
	"math/rand"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type Admin_LoginData struct {
	AdminID int    `json:"adminid" binding:"required"`
	AdminPw string `json:"adminpw" binding:"required"`
	Key     string `json:"key" binding:"required"`
}

// 要改
func Admin_Login(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	var adminlogin Admin_LoginData
	err := c.ShouldBindJSON(&adminlogin)
	if err != nil {
		utils.ResponseError(c, utils.ParameterErrorCode, utils.ParameterErrorMsg)
		return
	}
	token, err := Admin_LoginHandler(adminlogin)
	if err != nil {
		if err == utils.ErrFormatWrong {
			utils.ResponseError(c, utils.ParameterErrorCode, utils.FormatWrongMsg)
			return
		} else if err == utils.ErrUserNotFound {
			utils.ResponseError(c, utils.NotFoundCode, utils.UserNotFoundMsg)
			return
		} else if err == utils.ErrCopyFail {
			utils.ResponseError(c, utils.OperationFailedCode, utils.CopyFailMsg)
			return
		} else if err == utils.ErrOperationFailed {
			utils.ResponseError(c, utils.OperationFailedCode, utils.OperationFailedMsg)
			return
		}
		utils.ResponseError(c, utils.OperationFailedCode, utils.OperationFailedMsg)
		return
	}
	utils.ResponseSuccess(c, token)
	return
}

// 登录
// 生成随机密钥
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRandomString(length int) string {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rng.Intn(len(charset))]
	}

	return string(result)
}

var secretKey = []byte(generateRandomString(10))

type Claims struct {
	AdminID int `json:"adminid"`
	Type    int `json:"type"`
	jwt.StandardClaims
}

func Admin_LoginHandler(data interface{}) (string, error) {
	var logindata Admin_LoginData
	if err := copier.Copy(&logindata, data); err != nil {
		return "", utils.ErrCopyFail
	}
	err := server.Check_Admin(logindata.AdminID, logindata.AdminPw, logindata.Key)
	if err != nil {
		if err == utils.ErrUserNotFound {
			return "", utils.ErrUserNotFound
		}
		return "", utils.ErrOperationFailed
	}

	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		AdminID: logindata.AdminID,
		Type:    2,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
