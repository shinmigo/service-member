package address

import (
	"fmt"
	
	"goshop/service-member/pkg/db"
	"goshop/service-member/pkg/utils"
)

type Address struct {
	AddressId     uint64         `json:"address_id" gorm:"PRIMARY_KEY"`
	MemberId      uint64         `json:"member_id"`
	Name          string         `json:"name"`
	Mobile        string         `json:"mobile"`
	AddressDetail string         `json:"address_detail"`
	RoomNumber    string         `json:"room_number"`
	CodeProv      uint64         `json:"code_prov"`
	CodeCity      uint64         `json:"code_city"`
	CodeCoun      uint64         `json:"code_coun"`
	CodeTown      uint64         `json:"code_town"`
	IsDefault     int32          `json:"is_default"`
	Longitude     string         `json:"longitude"`
	Latitude      string         `json:"latitude"`
	CreatedAt     utils.JSONTime `json:"created_at"`
	UpdatedAt     utils.JSONTime `json:"updated_at"`
}

func GetTableName() string {
	return "address"
}

func GetField() []string {
	return []string{
		"address_id", "member_id", "name", "mobile", "address_detail", "room_number",
		"code_prov", "code_city", "code_coun", "code_town",
		"is_default", "longitude", "latitude",
	}
}

func GetOneByAddressId(AddressId uint64) (*Address, error) {
	if AddressId == 0 {
		return nil, fmt.Errorf("address_id is null")
	}
	row := &Address{}
	err := db.Conn.Table(GetTableName()).
		Select(GetField()).
		Where("address_id = ?", AddressId).
		First(row).Error
	
	if err != nil {
		return nil, fmt.Errorf("err: %v", err)
	}
	return row, nil
}

func GetAddressListByMemberId(memberId uint64, page, pageSize uint64) ([]*Address, uint64, error) {
	var total uint64
	rows := make([]*Address, 0, pageSize)
	
	query := db.Conn.Table(GetTableName()).Select(GetField()).Where("member_id = ?", memberId)
	
	err := query.Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("is_default desc").
		Order("address_id desc").
		Find(&rows).Error
	
	if err != nil {
		return nil, total, err
	}
	
	query.Count(&total)
	
	return rows, total, nil
}
