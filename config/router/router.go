package router

import (
	admincontroller "backend/app/controller/AdminController"
	msgcontroller "backend/app/controller/MsgController"
	teamcontroller "backend/app/controller/TeamController"
	usercontroller "backend/app/controller/UserController"

	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {
	const pre = "/api"
	api := r.Group(pre)
	api.POST("/user/register", usercontroller.Register)
	api.POST("/user/login", usercontroller.Login)

	user := api.Group("/user")
	{
		user.Use(usercontroller.AuthMiddleware())
		user.DELETE("/delete", usercontroller.DeleteUser)
		user.POST("/uploadavataraurl", usercontroller.UpLoadAvatara)
		user.POST("/updateprofile", usercontroller.UpdateUserProfile)
		user.GET("/getprofile", usercontroller.GetUser)
	}
	team := api.Group("/team")
	{
		team.Use(usercontroller.AuthMiddleware())
		team.POST("/create", teamcontroller.CreateTeam)
		team.POST("/join", teamcontroller.AddMember)
		team.DELETE("/quit", teamcontroller.QuitMember)
		team.GET("/get", teamcontroller.GetTeam)
		team.GET("/geteammember", teamcontroller.GetTeamMember)
		team.DELETE("/del", teamcontroller.DelTeam)
		team.POST("/updateprofile", teamcontroller.UpdateTeamProfile)
		team.POST("/submit", teamcontroller.Submit)
		team.POST("/cancel", teamcontroller.Cancel)
		team.POST("/uploadavataraurl", teamcontroller.UpLoadTeamAvatara)
	}
	msg := api.Group("/msg")
	{
		msg.Use(usercontroller.AuthMiddleware())
		msg.POST("/updateunread", msgcontroller.UpdateUnread)
		msg.GET("/getmsg", msgcontroller.GetMsg)
	}
	admin := api.Group("/admin")
	{
		admin.Use(usercontroller.AuthMiddleware())
		admin.GET("/getalluser", admincontroller.GetAllUser)
		admin.GET("/getallteam", admincontroller.GetAllTeam)
		admin.DELETE("/delteam", admincontroller.DelTeam_Admin)
		admin.DELETE("/deluser", admincontroller.DelUser_Admin)
	}
}
