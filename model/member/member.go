package member

import (
	"fmt"
	
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"goshop/service-member/pkg/db"
	"goshop/service-member/pkg/utils"
)

type Member struct {
	MemberId      uint64         `json:"member_id" gorm:"PRIMARY_KEY"`
	Password      string         `json:"password"`
	Nickname      string         `json:"nickname"`
	Mobile        string         `json:"mobile"`
	RegisterIp    string         `json:"register_ip"`
	Name          string         `json:"name"`
	Gender        int32          `json:"gender"`
	IdCard        string         `json:"id_card"`
	Birthday      string         `json:"birthday"`
	Avatar        string         `json:"avatar"`
	Email         string         `json:"email"`
	Status        int32          `json:"status"`
	LastLoginTime utils.JSONTime `json:"last_login_time"`
	Remark        string         `json:"remark"`
	MemberLevelId uint64         `json:"member_level_id"`
	Point         int64          `json:"point"`
	Balance       float64        `json:"balance"`
	CreatedBy     uint64         `json:"created_by"`
	UpdatedBy     uint64         `json:"updated_by"`
	CreatedAt     utils.JSONTime `json:"created_at"`
	UpdatedAt     utils.JSONTime `json:"updated_at"`
}

func GetTableName() string {
	return "member"
}

func GetField() []string {
	return []string{
		"member_id", "password", "nickname", "mobile", "register_ip", "name", "gender", "id_card", "birthday",
		"avatar", "email", "status", "last_login_time", "remark", "member_level_id", "point", "balance",
		"created_by", "updated_by", "created_at", "updated_at",
	}
}

func (m *Member) BeforeSave(scope *gorm.Scope) (err error) {
	if len(m.Password) > 0 {
		if pw, err := bcrypt.GenerateFromPassword([]byte(m.Password), 0); err == nil {
			scope.SetColumn("password", pw)
		}
	}
	return
}

func GetOneByMemberId(MemberId uint64) (*Member, error) {
	if MemberId == 0 {
		return nil, fmt.Errorf("member_id is null")
	}
	row := &Member{}
	err := db.Conn.Table(GetTableName()).
		Select(GetField()).
		Where("member_id = ?", MemberId).
		First(row).Error
	
	if err != nil {
		return nil, fmt.Errorf("err: %v", err)
	}
	return row, nil
}

func GetOneByMobile(mobile string) (*Member, error) {
	if len(mobile) == 0 {
		return nil, fmt.Errorf("mobile is null")
	}
	
	row := &Member{}
	err := db.Conn.Table(GetTableName()).
		Select(GetField()).
		Where("mobile = ?", mobile).
		First(row).Error
	
	if err != nil {
		return nil, fmt.Errorf("err: %v", err)
	}
	return row, nil
}

func GetMemberList(memberId uint64, status int32, mobile string, page, pageSize uint64, nickname string) ([]*Member, uint64, error) {
	var total uint64
	
	rows := make([]*Member, 0, pageSize)
	
	query := db.Conn.Table(GetTableName()).Select(GetField())
	if memberId > 0 {
		query = query.Where("member_id = ?", memberId)
	}
	
	if status > 0 {
		query = query.Where("status = ?", status)
	}
	
	if mobile != "" {
		query = query.Where("mobile = ?", mobile)
	}
	
	if len(nickname) > 0 {
		query = query.Where("nickname like ?", nickname+"%")
	}
	
	err := query.Order("member_id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&rows).Error
	if err != nil {
		return nil, total, err
	}
	
	query.Count(&total)
	
	return rows, total, nil
}
