package server

import (
	"backend/app/model"
	"backend/app/utils"
	"backend/config/database"
	"regexp"

	"github.com/jinzhu/copier"
)

// 注册

// Format  格式错误会返回一个true，正确返回false
func VerifyEmailFormat(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return !reg.MatchString(email)
}

// 查询用户密码
func Check(Email string, Password string) (*model.Usermodel, error) {
	var user model.Usermodel
	err := database.DB.Where("email = ? AND password = ?", Email, Password).First(&user).Error
	return &user, err
}

// 账号存在性
func CheckUserExistByEmail(Email string) bool {
	err := database.DB.Where("email = ?", Email).First(&model.Usermodel{}).Error
	return err == nil
}

// 创建用户
func CreateNewUser(username string, password string, email string) bool {
	db := database.DB
	newUser := &model.Usermodel{
		Username: username,
		Password: password,
		Email:    email,
		// TeamID:   20230000,
	}
	err := db.Create(newUser)
	return err.Error != nil
}

// 查询用户
func GetUser(Email string) (*model.Usermodel, error) {
	var user model.Usermodel
	err := database.DB.Where("email=?", Email).First(&user)
	if err.Error != nil {
		// fmt.Println("a")
		return nil, err.Error
	}
	return &user, nil
}

// 注销账号
func DeleteUser(Email string, Password string) error {
	db := database.DB
	err := db.Delete(&model.Usermodel{}, "email=? AND password = ?", Email, Password)
	return err.Error
}

func DeleteUser_Admin(Email string) error {
	db := database.DB
	err := db.Delete(&model.Usermodel{}, "email=?", Email)
	return err.Error
}

// 注册

func Register(data interface{}) error {
	var registerdata model.Usermodel
	db := database.DB
	err := copier.Copy(&registerdata, data)

	if err != nil {
		return utils.ErrCopyFail
	} else if VerifyEmailFormat(registerdata.Email) {
		return utils.ErrFormatWrong
	} else if CheckUserExistByEmail(registerdata.Email) {
		return utils.ErrUserHaveExist
	} else if err := db.Create(&registerdata).Error; err != nil {
		return utils.ErrOperationFailed
	}
	return nil
}

// 上传头像
func UpLoadAvatara(Email string, AvataraUrl string) error {
	if VerifyEmailFormat(Email) {
		return utils.ErrFormatWrong
	}
	user, err := GetUser(Email)
	if err != nil {
		return utils.ErrUserNotFound
	}
	db := database.DB
	err = db.Model(&user).Update("AvataraUrl", AvataraUrl).Error
	return err
}

// 更新资料
type NewUserfile struct {
	Username    string
	Sign        string
	Description string
}

func UpdateUserProfile(Email string, data interface{}) error {
	var newfile NewUserfile
	if err := copier.Copy(&newfile, data); err != nil {
		return utils.ErrCopyFail
	}
	if !CheckUserExistByEmail(Email) {
		return utils.ErrUserNotFound
	}
	db := database.DB
	err := db.Model(&model.Usermodel{}).Where("email = ?", Email).Updates(newfile).Error
	return err
}

func GetAllUser() ([]model.Usermodel, error) {
	db := database.DB
	var users []model.Usermodel
	err := db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
