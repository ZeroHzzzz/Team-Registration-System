package server

import (
	"backend/app/model"
	"backend/app/utils"
	"backend/config/database"

	"github.com/jinzhu/copier"
)

// 存在性
func CheckTeamExist(TeamID string) error {
	result := database.DB.Where("team_id = ?", TeamID).First(&model.Team{})
	return result.Error
}

// 获取团队信息
func GetTeam(TeamID string) (*model.Team, error) {
	var team model.Team
	err := database.DB.Preload("Users").Where("team_id=?", TeamID).First(&team)
	if err.Error != nil {
		return nil, err.Error
	}
	return &team, nil
}

func CheckTeam(TeamID string, TeamPassword string) error {
	err := database.DB.Where("team_id = ? AND team_password = ?", TeamID, TeamPassword).First(&model.Team{})
	return err.Error
}

func IsInTeam(Email string) bool {
	User, err := GetUser(Email)
	if err != nil {
		return false
	}
	return User.TeamID != nil
}

func TeamRegister(data interface{}) (string, error) {
	var registerdata model.Team
	db := database.DB
	err := copier.Copy(&registerdata, data)
	if err != nil {
		return "", utils.ErrCopyFail
	} else if VerifyEmailFormat(registerdata.LeaderID) {
		return "", utils.ErrFormatWrong
	} else if !CheckUserExistByEmail(registerdata.LeaderID) {
		return "", utils.ErrUserNotFound
	} else if IsInTeam(registerdata.LeaderID) {
		return "", utils.ErrHaveInTeam
	} else if err = db.Create(&registerdata).Error; err != nil {
		return "", utils.ErrOperationFailed
	} else if err = AddMember(registerdata.LeaderID, registerdata.TeamID, registerdata.TeamPassword); err != nil {
		DelTeam(registerdata.TeamID, registerdata.TeamPassword)
		return "", utils.ErrAddMemberOperationFailed
	}

	user, _ := GetUser(registerdata.LeaderID)
	err = db.Model(&user).Update("type", 1).Error
	if err != nil {
		return "", utils.ErrOperationFailed
	}

	return registerdata.TeamID, nil
}

func AddMember(Email string, TeamID string, Password string) error { // Password 是团队的
	db := database.DB
	if VerifyEmailFormat(Email) {
		return utils.ErrFormatWrong
	}
	member, err := GetUser(Email)
	if err != nil {
		return utils.ErrUserNotFound
	} else if IsInTeam(Email) {
		return utils.ErrHaveInTeam
	}
	team, err := GetTeam(TeamID)
	if err != nil {
		return utils.ErrTeamNotFound
	}
	err = db.Model(&team).Association("Users").Append(member)
	db.Save(&team)
	if err != nil {
		return utils.ErrOperationFailed
	}
	msgcontent := member.Username + "加入团队"
	err = CreateMsg(TeamID, msgcontent)
	return err
}

func QuitTeam(Email string, TeamID string) error {
	db := database.DB
	var team model.Team
	var member model.Usermodel
	if VerifyEmailFormat(Email) {
		return utils.ErrFormatWrong
	}
	if err := db.Where("team_id = ?", TeamID).First(&team).Error; err != nil {
		return err
	}
	if err := db.Where("email = ? AND team_id = ?", Email, TeamID).First(&member).Error; err != nil {
		return err
	}
	if member.Type == 1 {
		err := DelTeam(TeamID, team.TeamPassword)
		if err != nil {
			return utils.ErrOperationFailed
		}
	}
	err := db.Model(&member).Update("type", 0).Error
	if err != nil {
		return err
	}
	err = db.Model(&team).Association("Users").Delete(&member)
	return err
}

func DelTeam(TeamID string, TeamPassword string) error {
	db := database.DB
	err := CheckTeam(TeamID, TeamPassword)
	if err != nil {
		return utils.ErrTeamNotFound
	}
	var team model.Team
	db.Preload("Users").Where("team_id=?", TeamID).First(&team, TeamID)
	// fmt.Println(team.LeaderID)
	leader, _ := GetUser(team.LeaderID)
	err = db.Model(&leader).Update("type", 0).Error
	if err != nil {
		return utils.ErrOperationFailed
	}

	db.Model(&team).Association("Users").Delete(&team.Users)
	err = db.Delete(&model.Team{}, "team_id=?", TeamID).Error
	if err != nil {
		return utils.ErrDelTeamFailed
	}
	msgcontent := "团队解散"
	err = CreateMsg(TeamID, msgcontent)
	return err
}

func DelTeam_Admin(TeamID string) error {
	db := database.DB
	err := CheckTeamExist(TeamID)
	if err != nil {
		return utils.ErrTeamNotFound
	}
	var team model.Team
	db.Preload("Users").Where("team_id=?", TeamID).First(&team, TeamID)

	leader, _ := GetUser(team.LeaderID)
	err = db.Model(&leader).Update("type", 0).Error
	if err != nil {
		return utils.ErrOperationFailed
	}

	db.Model(&team).Association("Users").Delete(&team.Users)
	err = db.Delete(&model.Team{}, "team_id=?", TeamID).Error
	if err != nil {
		return utils.ErrDelTeamFailed
	}
	msgcontent := "团队已被管理员解散"
	err = CreateMsg(TeamID, msgcontent)
	return err
}

func GetTeamMember(TeamID string) ([]model.Usermodel, error) {
	db := database.DB
	var team model.Team
	var users []model.Usermodel

	err := db.Preload("Users").First(&team, TeamID).Error
	if err != nil {
		return nil, err
	}

	users = team.Users
	return users, nil
}

func Submit(TeamID string, TeamPassword string) error {
	err := CheckTeam(TeamID, TeamPassword)
	if err != nil {
		return utils.ErrTeamNotFound
	}
	team, _ := GetTeam(TeamID)
	if len(team.Users) > 6 || len(team.Users) < 4 {
		return utils.ErrSubmitFailed
	}
	db := database.DB
	err = db.Model(&team).Update("state", true).Error
	if err != nil {
		return utils.ErrOperationFailed
	}
	msgcontent := "队长已提交报名"
	err = CreateMsg(TeamID, msgcontent)
	return err
}

func Cancel(TeamID string) error {
	err := CheckTeamExist(TeamID)
	if err != nil {
		return utils.ErrTeamNotFound
	}
	team, _ := GetTeam(TeamID)
	db := database.DB
	err = db.Model(&team).Update("state", false).Error
	if err != nil {
		return utils.ErrOperationFailed
	}
	msgcontent := "队长撤销报名"
	err = CreateMsg(TeamID, msgcontent)
	return err
}

type NewTeamfile struct {
	TeamName        string
	TeamDescription string
}

func UpdateTeamProfile(TeamID string, data interface{}) error {
	var newfile NewTeamfile
	if err := copier.Copy(&newfile, data); err != nil {
		return utils.ErrCopyFail
	}
	if err := CheckTeamExist(TeamID); err != nil {
		return utils.ErrTeamNotFound
	}
	db := database.DB
	err := db.Model(&model.Team{}).Where("team_id = ?", TeamID).Updates(newfile).Error
	return err
}

func GetAllTeam() ([]model.Team, error) {
	db := database.DB
	var teams []model.Team
	err := db.Find(&teams).Error
	if err != nil {
		return nil, err
	}
	return teams, nil
}

func UpLoadTeamAvatara(TeamID string, AvataraUrl string) error {
	team, err := GetTeam(TeamID)
	if err != nil {
		return utils.ErrTeamNotFound
	}
	db := database.DB
	err = db.Model(&team).Update("AvataraUrl", AvataraUrl).Error
	return err
}
