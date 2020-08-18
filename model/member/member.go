package member

import (
	"goshop/service-member/pkg/db"
	"goshop/service-member/pkg/utils"

	"github.com/shinmigo/pb/memberpb"
)

var StatusMap map[uint32]string

func init() {
	StatusMap = make(map[uint32]string, 2)

	StatusMap[2] = "冻结"
	StatusMap[1] = "正常"
}

type Member struct {
	MemberId      uint64 `gorm:"column:id;PRIMARY_KEY"`
	Nickname      string `json:"nickname"`
	Status        uint32 `json:"status"`
	Mobile        string `json:"mobile"`
	MemberLevelId string `json:"member_level_id"`
	Gender        uint32 `json:"gender"`
	Birthday      string `json:"birthday"`
	CreatedBy     string `json:"created_by"`
	UpdatedBy     string `json:"updated_by"`
}

func GetTableName() string {
	return "member"
}

func GetInfoFields() []string {
	return []string{
		"member_id", "nickname", "mobile", "register_ip", "name", "gender", "id_card", "birthday", "avatar", "email", "status", "register_time", "last_login_time", "remark", "member_level_id", "point", "balance", "created_at", "created_by", "updated_at", "updated_by",
	}
}

func GetList(args *memberpb.ListReq) ([]*memberpb.Member, uint64, error) {
	query := db.Conn.Table(GetTableName()).Select(GetInfoFields())
	if len(args.Mobile) > 0 {
		query = query.Where("mobile = ?", args.Mobile)
	}

	if args.MemberId > 0 {
		query = query.Where("member_id = ?", args.MemberId)
	}

	if args.Status > 0 {
		query = query.Where("status = ?", args.Status)
	}

	var total uint64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if args.PageSize == 0 {
		args.PageSize = 20
	}

	data := make([]*memberpb.Member, 0, args.PageSize)

	query = query.Limit(args.PageSize).Offset(utils.GetPageOffset(args.PageSize, args.Page)).Order("member_id desc")
	if err := query.Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}
