package server

import (
	"backend/app/model"
	"backend/app/utils"
	"backend/config/database"
)

// 由于信息只由团队产生，因此只根据团队创建信息
func CreateMsg(TeamID string, Content string) error {
	db := database.DB
	info := model.Message{
		Content: Content,
	}
	err := db.Create(&info).Error
	if err != nil {
		return utils.ErrCreateMsgFailed
	}
	users, _ := GetTeamMember(TeamID)
	for _, user := range users {
		db.Model(&user).Association("Message").Append(&info)
		user.Unread++
		db.Save(&user)
	}
	return nil
}

func GetMessage(MsgID int) (*model.Message, error) {
	var msg model.Message
	err := database.DB.Where("id = ?", MsgID).First(&msg)
	if err.Error != nil {
		return nil, err.Error
	}
	return &msg, nil
}
func UpdateUnread(Email string) error {
	db := database.DB
	if VerifyEmailFormat(Email) {
		return utils.ErrFormatWrong
	}
	user, err := GetUser(Email)
	if err != nil {
		if err == utils.ErrUserNotFound {
			return utils.ErrUserNotFound
		}
		return utils.ErrOperationFailed
	}
	user.Unread = 0
	db.Save(&user)
	return nil
}

func GetMsg(Email string) ([]model.Message, error) {
	db := database.DB
	var user model.Usermodel
	var msgs []model.Message
	if VerifyEmailFormat(Email) {
		return nil, utils.ErrFormatWrong
	}
	err := db.Preload("Message").Where("email = ?", Email).First(&user).Error
	if err != nil {
		return nil, err
	}
	msgs = user.Message
	return msgs, nil
}
