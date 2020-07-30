package rpc

import (
	"goshop/service-member/model/account"
	"goshop/service-member/model/member"
	"goshop/service-member/pkg/db"

	"github.com/shinmigo/pb/memberpb"
	"golang.org/x/net/context"
)

type Member struct {
}

func NewMember() *Member {
	return &Member{}
}

// 获取会员列表
func (s *Member) GetList(ctx context.Context, args *memberpb.ListReq) (*memberpb.ListRes, error) {

	data := make([]*memberpb.Member, 0, 32)
	query := db.Conn.Table(member.GetTableName()).Select(member.GetDetailFields())
	if len(args.Mobile) > 0 {
		query = query.Where("mobile = ?", args.Mobile)
	}

	if args.MemberId > 0 {
		query = query.Where("member_id = ?", args.MemberId)
	}

	if args.Status > 0 {
		query = query.Where("status = ?", args.Status)
	}
	err := query.Find(&data).Error
	if err != nil {
		return nil, err

	}

	for _, item := range data {
		item.MemberLevelName = "V1"
		item.StatusText = member.StatusMap[item.Status]
	}

	return &memberpb.ListRes{Members: data}, nil
}

// 获取详情 TODO 没有记录
func (s *Member) GetDetail(ctx context.Context, args *memberpb.DetailReq) (*memberpb.Member, error) {
	response := &memberpb.Member{}
	err := db.Conn.Table(member.GetTableName()).Select(member.GetDetailFields()).Where("member_id = ?", args.MemberId).First(response).Error
	if err != nil {
		return nil, err
	}

	response.StatusText = member.StatusMap[response.Status]
	response.MemberLevelName = "V1"

	return response, nil
}

// 创建会员
func (s *Member) Create(ctx context.Context, args *memberpb.CreateReq) (*memberpb.Res, error) {
	var err error

	tr := db.Conn.Begin()
	if tr.Error != nil {
		return nil, tr.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tr.Rollback()
			panic(r)
		}

		if err != nil {
			tr.Rollback()
		}
	}()

	// 创建会员
	memberData := member.Member{
		Mobile:        args.Mobile,
		Nickname:      args.Nickname,
		Status:        args.Status,
		MemberLevelId: args.MemberLevelId,
		Gender:        args.Gender,
		Birthday:      args.Birthday,
		CreatedBy:     args.Operator,
	}

	if err = tr.Table(member.GetTableName()).Create(&memberData).Error; err != nil {
		return nil, err
	}

	//获取新增的会员Id
	var memberId []uint64
	tr.Raw("select LAST_INSERT_ID() as id").Pluck("member_id", &memberId)

	accountData := account.Account{
		MemberId:    memberId[0],
		AccountName: args.Mobile,
		Category:    1,
		Password:    "111111",
		Status:      args.Status,
	}

	if err = tr.Table(account.GetTableName()).Create(accountData).Error; err != nil {
		return nil, err
	}

	tr.Commit()

	return &memberpb.Res{Status: true}, nil
}

// 更新会员
func (s *Member) Update(ctx context.Context, args *memberpb.UpdateReq) (*memberpb.Res, error) {

	updateValue := map[string]interface{}{
		"nickname":        args.Nickname,
		"mobile":          args.Mobile,
		"gender":          args.Gender,
		"birthday":        args.Birthday,
		"member_level_id": args.MemberLevelId,
		"updated_by":      args.Operator,
	}
	err := db.Conn.Table(member.GetTableName()).Where("member_id = ?", args.MemberId).Update(updateValue).Error
	if err != nil {
		return nil, err
	}

	return &memberpb.Res{Status: true}, nil
}

// 更新会员状态 status 0 冻结 1 解冻
func (s *Member) UpdateStatus(ctx context.Context, args *memberpb.UpdateStatusReq) (*memberpb.Res, error) {

	updateValue := map[string]interface{}{
		"status":     args.Status,
		"updated_by": args.Operator,
	}
	err := db.Conn.Table(member.GetTableName()).Where("member_id = ?", args.MemberId).Update(updateValue).Error
	if err != nil {
		return nil, err
	}

	return &memberpb.Res{Status: true}, nil
}
