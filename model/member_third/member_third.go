package member_third

import (
	"fmt"
	
	"goshop/service-member/pkg/db"
	"goshop/service-member/pkg/utils"
)

type MemberThird struct {
	Id         uint64         `json:"id" gorm:"PRIMARY_KEY"`
	MemberId   uint64         `json:"member_id"`
	Type       uint32         `json:"type"`
	OpenId     string         `json:"open_id"`
	SessionKey string         `json:"session_key"`
	Unionid    string         `json:"unionid"`
	CreatedAt  utils.JSONTime `json:"created_at"`
	UpdatedAt  utils.JSONTime `json:"updated_at"`
}

func GetTableName() string {
	return "member_third"
}

func GetField() []string {
	return []string{
		"id", "member_id", "type", "open_id", "session_key", "unionid", "created_at", "updated_at",
	}
}

func GetOneByOpenId(openId string, typeStr uint32) (*MemberThird, error) {
	if len(openId) == 0 {
		return nil, fmt.Errorf("open_id is null")
	}
	row := &MemberThird{}
	err := db.Conn.Table(GetTableName()).
		Select(GetField()).
		Where("open_id = ? and type = ?", openId, typeStr).
		First(row).Error
	
	if err != nil {
		return nil, fmt.Errorf("err: %v", err)
	}
	return row, nil
}
