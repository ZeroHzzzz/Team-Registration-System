package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Team struct {
	//ID       int    `gorm:"primaryKey"`
	TeamName        string `gorm:"not null type:varchar(25)"`
	TeamPassword    string `json:"-"`
	State           bool
	TeamDescription string `gorm:"type:varchar(25)"`
	LeaderID        string // Email
	TeamID          string `gorm:"primarykey;unique"`
	AvataraUrl      string
	Users           []Usermodel `gorm:"foreignKey:TeamID"`
}

// 修改id规则

func (team *Team) BeforeCreate(tx *gorm.DB) error {
	team.TeamID = "Team_" + fmt.Sprint(time.Now().Unix())
	return nil

}

// func generateCustomTeamID(tx *gorm.DB) int {
// 	var count int64
// 	tx.Model(&Team{}).Count(&count)
// 	return int(20230000 + count + 1)
// }

func (team *Team) BeforeDelete(tx *gorm.DB) error {
	teamID := team.TeamID
	if err := tx.Model(&Usermodel{}).Where("team_id=?", teamID).Update("team_id", nil).Error; err != nil {
		return err
	}
	return nil
}
