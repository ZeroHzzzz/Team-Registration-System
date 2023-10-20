package usercontroller

import (
	"backend/app/model"
	"backend/app/server"
	"backend/app/utils"
	"math/rand"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

// 解析信息
type LoginData struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	var logindata LoginData
	err := c.ShouldBindJSON(&logindata)
	if err != nil {
		utils.ResponseError(c, utils.ParameterErrorCode, utils.ParameterErrorMsg)
		return
	}
	returnfile, token, err := LoginHandler(logindata)
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

	responseData := map[string]interface{}{
		"returnfile": returnfile,
		"token":      token,
	}

	utils.ResponseSuccess(c, responseData)
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
	Email string `json:"email"`
	Type  int    `json:"type"`
	jwt.StandardClaims
}

func LoginHandler(data interface{}) (model.Usermodel, string, error) {
	var logindata LoginData
	if err := copier.Copy(&logindata, data); err != nil {
		return model.Usermodel{}, "", utils.ErrCopyFail
	}
	if server.VerifyEmailFormat(logindata.Email) {
		return model.Usermodel{}, "", utils.ErrFormatWrong
	}
	userfile, err := server.Check(logindata.Email, logindata.Password)
	if err != nil {
		if err == utils.ErrUserNotFound {
			return model.Usermodel{}, "", utils.ErrUserNotFound
		}
		return model.Usermodel{}, "", utils.ErrOperationFailed
	}

	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		Email: logindata.Email,
		Type:  userfile.Type,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return *userfile, "", err
	}
	return *userfile, tokenString, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		// Parse and validate the token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			utils.ResponseUnauthorized(c)
			c.Abort()
			return
		}

		// Set the user ID in the context for further use
		claims, accepted := token.Claims.(*Claims)
		if !accepted {
			utils.ResponseUnauthorized(c)
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("type", claims.Type)
		c.Next()
	}
}

// 验证
// JWT中间件
// func AuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		tokenString := c.GetHeader("Authorization")

// 		// Parse and validate the token
// 		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
// 			return secretKey, nil
// 		})

// 		if err != nil || !token.Valid {
// 			utils.ResponseUnauthorized(c)
// 			c.Abort()
// 			return
// 		}
// 		if !token.Valid {
// 			// Set the user ID in the context for further use
// 			claims, accepted := token.Claims.(*Claims)
// 			if !accepted {
// 				utils.ResponseUnauthorized(c)
// 				c.Abort()
// 				return
// 			}

// 			currentTime := time.Now().Unix()
// 			if currentTime > claims.ExpiresAt {
// 				// 令牌过期，尝试刷新
// 				refreshedToken, refreshErr := RefreshToken(claims.Email, claims.Type)

// 				if refreshErr != nil {
// 					utils.ResponseUnauthorized(c)
// 					c.Abort()
// 					return
// 				}
// 				c.Header("Authorization", refreshedToken)
// 			} else {
// 				utils.ResponseUnauthorized(c)
// 				c.Abort()
// 				return
// 			}
// 			c.Set("email", claims.Email)
// 			c.Set("type", claims.Type)
// 			c.Next()
// 			return
// 		}

// 	}
// }

// func RefreshToken(userID string, Type int) (string, error) {
// 	newToken, err := GenerateNewToken(userID, Type)
// 	if err != nil {
// 		return "", err
// 	}

// 	return newToken, nil
// }

// func GenerateNewToken(userID string, Type int) (string, error) {
// 	user, _ := server.GetUser(userID)
// 	expirationTime := time.Now().Add(1 * time.Hour)
// 	newClaims := &Claims{
// 		Email: userID,
// 		Type:  user.Type,
// 		StandardClaims: jwt.StandardClaims{
// 			ExpiresAt: expirationTime.Unix(),
// 			IssuedAt:  time.Now().Unix(),
// 		},
// 	}
// 	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
// 	tokenString, err := newToken.SignedString(secretKey)

// 	if err != nil {
// 		return "", err
// 	}

// 	return tokenString, nil
// }

// func AuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		tokenString := c.GetHeader("Authorization")

// 		// Parse and validate the token
// 		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
// 			return secretKey, nil
// 		})

// 		if err != nil {
// 			utils.ResponseUnauthorized(c)
// 			c.Abort()
// 			return
// 		}

// 		if token.Valid {
// 			// Token is valid, set user information in the context
// 			claims, accepted := token.Claims.(*Claims)
// 			if !accepted {
// 				utils.ResponseUnauthorized(c)
// 				c.Abort()
// 				return
// 			}

// 			c.Set("email", claims.Email)
// 			c.Set("type", claims.Type)
// 			c.Next()
// 			return
// 		}

// 		if ve, ok := err.(*jwt.ValidationError); ok {
// 			if ve.Errors&jwt.ValidationErrorExpired != 0 {
// 				// Token is expired, attempt to refresh it
// 				refreshToken := c.GetHeader("Refresh-Token")
// 				newToken, err := refreshAccessToken(refreshToken)
// 				if err != nil {
// 					utils.ResponseUnauthorized(c)
// 					c.Abort()
// 					return
// 				}

// 				// Set the new token in the response header
// 				c.Header("Authorization", "Bearer "+newToken)

// 				// Parse and validate the new token
// 				newTokenClaims := &Claims{}
// 				newToken, err := jwt.ParseWithClaims(newToken, newTokenClaims, func(token *jwt.Token) (interface{}, error) {
// 					return secretKey, nil
// 				})

// 				if err != nil || !newToken.Valid {
// 					utils.ResponseUnauthorized(c)
// 					c.Abort()
// 					return
// 				}

// 				// Token is now valid, set user information in the context
// 				c.Set("email", newTokenClaims.Email)
// 				c.Set("type", newTokenClaims.Type)
// 				c.Next()
// 				return
// 			}
// 		}

// 		utils.ResponseUnauthorized(c)
// 		c.Abort()
// 	}
// }

// func AuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		tokenString := c.GetHeader("Authorization")

// 		// Parse and validate the token
// 		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
// 			return secretKey, nil
// 		})

// 		if err != nil {
// 			utils.ResponseUnauthorized(c)
// 			c.Abort()
// 			return
// 		}

// 		if token.Valid {
// 			// Token is valid, set user information in the context
// 			claims, accepted := token.Claims.(*Claims)
// 			if !accepted {
// 				utils.ResponseUnauthorized(c)
// 				c.Abort()
// 				return
// 			}

// 			c.Set("email", claims.Email)
// 			c.Set("type", claims.Type)
// 			c.Next()
// 			return
// 		}

// 		if ve, ok := err.(*jwt.ValidationError); ok {
// 			if ve.Errors&jwt.ValidationErrorExpired != 0 {
// 				// Token is expired, attempt to refresh it
// 				refreshToken := c.GetHeader("Authorization")
// 				newToken, err := refreshAccessToken(refreshToken)
// 				if err != nil {
// 					utils.ResponseUnauthorized(c)
// 					c.Abort()
// 					return
// 				}

// 				// Set the new token in the response header
// 				// c.Header("Authorization", "Bearer "+newToken)

// 				// Parse and validate the new token
// 				newTokenClaims := &Claims{}
// 				parsedNewToken, err := jwt.ParseWithClaims(newToken, newTokenClaims, func(token *jwt.Token) (interface{}, error) {
// 					return secretKey, nil
// 				})

// 				if err != nil || !parsedNewToken.Valid {
// 					utils.ResponseUnauthorized(c)
// 					c.Abort()
// 					return
// 				}

// 				// Token is now valid, set user information in the context
// 				c.Set("email", newTokenClaims.Email)
// 				c.Set("type", newTokenClaims.Type)
// 				c.Next()
// 				return
// 			}
// 		}

// 		utils.ResponseUnauthorized(c)
// 		c.Abort()
// 	}
// }

// // 刷新访问令牌
// func refreshAccessToken(refreshToken string) (string, error) {
// 	// 解析刷新令牌
// 	claims := &RefreshTokenClaims{}
// 	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
// 		return secretKey, nil
// 	})

// 	if err != nil || !token.Valid {
// 		return "", errors.New("Invalid refresh token")
// 	}

// 	// TODO: 在这里添加逻辑，根据用户ID或其他信息生成新的访问令牌
// 	// 这里可能需要调用你的登录逻辑或者使用其他方式获取新的访问令牌
// 	newAccessToken, err := generateAccessToken(claims.UserID)
// 	if err != nil {
// 		return "", err
// 	}

// 	return newAccessToken, nil
// }

// func generateAccessToken(userID string) (string, error) {
// 	expirationTime := time.Now().Add(1 * time.Hour)
// 	claims := &Claims{
// 		Email: userID, // 这里可以根据你的需求设置用户标识
// 		Type:  1,      // 假设 Type 为 1，可以根据你的需求设置
// 		StandardClaims: jwt.StandardClaims{
// 			ExpiresAt: expirationTime.Unix(),
// 			IssuedAt:  time.Now().Unix(),
// 		},
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	tokenString, err := token.SignedString(secretKey)
// 	if err != nil {
// 		return "", err
// 	}

// 	return tokenString, nil
// }
