package server

import (
	"backend/app/model"
	"backend/app/utils"
	"backend/config/database"
)

func Check_Admin(AdminID int, AdminPwd string, Key string) error {
	var admin model.Admin
	err := database.DB.Where("admin_id = ? AND admin_pwd = ?", AdminID, AdminPwd).First(&admin).Error
	if err != nil {
		return err
	}
	if Key != admin.Key {
		return utils.ErrOperationFailed
	}
	return nil
}
