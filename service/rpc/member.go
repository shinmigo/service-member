package rpc

import (
	"goshop/service-member/model/account"
	"goshop/service-member/model/member"
	"goshop/service-member/pkg/db"
	"goshop/service-member/service/rpc/logic"

	"github.com/shinmigo/pb/basepb"

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
	data, total, err := member.GetList(args)
	if err != nil {
		return nil, err
	}

	for _, item := range data {
		item.MemberLevelName = "V1"
		item.StatusText = member.StatusMap[item.Status]
	}

	return &memberpb.ListRes{Total: total, Members: data}, nil
}

// 获取详情 TODO 没有记录
func (s *Member) GetInfo(ctx context.Context, args *memberpb.InfoReq) (*memberpb.Member, error) {
	response := &memberpb.Member{}
	err := db.Conn.Table(member.GetTableName()).Select(member.GetInfoFields()).Where("member_id = ?", args.MemberId).First(response).Error
	if err != nil {
		return nil, err
	}

	response.StatusText = member.StatusMap[response.Status]
	response.MemberLevelName = "V1"

	return response, nil
}

// 创建会员
// TODO 缓存会员账号密码
func (s *Member) Add(ctx context.Context, args *memberpb.AddReq) (*basepb.AnyRes, error) {
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

	password, err := logic.GeneratePassword(args.Password) // 对密码进行加密， 加密算法：bcrypt
	if err != nil {
		return &basepb.AnyRes{State: 1}, err
	}

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

	accountData := account.Account{
		MemberId:    memberData.MemberId,
		AccountName: args.Mobile,
		Category:    1,
		Password:    password,
		Status:      args.Status,
	}

	if err = tr.Table(account.GetTableName()).Create(accountData).Error; err != nil {
		return nil, err
	}

	tr.Commit()

	return &basepb.AnyRes{State: 1, Id: memberData.MemberId}, nil
}

// 更新会员
func (s *Member) Edit(ctx context.Context, args *memberpb.EditReq) (*basepb.AnyRes, error) {

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

	return &basepb.AnyRes{State: 1}, nil
}

// 更新会员状态 status 0 冻结 1 解冻
func (s *Member) EditStatus(ctx context.Context, args *memberpb.EditStatusReq) (*basepb.AnyRes, error) {

	updateValue := map[string]interface{}{
		"status":     args.Status,
		"updated_by": args.Operator,
	}
	err := db.Conn.Table(member.GetTableName()).Where("member_id = ?", args.MemberId).Update(updateValue).Error
	if err != nil {
		return nil, err
	}

	return &basepb.AnyRes{State: 1}, nil
}
