package rpc

import (
	"fmt"
	
	"github.com/shinmigo/pb/basepb"
	"github.com/shinmigo/pb/memberpb"
	"golang.org/x/net/context"
	"goshop/service-member/model/member"
	"goshop/service-member/pkg/db"
	"goshop/service-member/pkg/utils"
)

type Member struct {
}

func NewMember() *Member {
	return &Member{}
}

func (s *Member) AddMember(ctx context.Context, req *memberpb.Member) (*basepb.AnyRes, error) {
	aul := member.Member{
		Nickname:      req.Nickname,
		Mobile:        req.Mobile,
		Name:          req.Name,
		Gender:        int32(req.Gender),
		IdCard:        req.IdCard,
		Birthday:      req.Birthday,
		Avatar:        req.Avatar,
		Email:         req.Email,
		Status:        int32(req.Status),
		Remark:        req.Remark,
		MemberLevelId: req.MemberLevelId,
		Point:         req.Point,
		Balance:       req.Balance,
		CreatedBy:     req.AdminId,
		UpdatedBy:     req.AdminId,
	}
	if err := db.Conn.Table(member.GetTableName()).Create(&aul).Error; err != nil {
		return nil, err
	}
	
	if utils.IsCancelled(ctx) {
		return nil, fmt.Errorf("client cancelled ")
	}
	
	return &basepb.AnyRes{
		Id:    aul.MemberId,
		State: 1,
	}, nil
}

func (s *Member) EditMember(ctx context.Context, req *memberpb.Member) (*basepb.AnyRes, error) {
	if _, err := member.GetOneByMemberId(req.MemberId); err != nil {
		return nil, err
	}
	
	aul := map[string]interface{}{
		"nickname":        req.Nickname,
		"mobile":          req.Mobile,
		"name":            req.Name,
		"gender":          int32(req.Gender),
		"id_card":         req.IdCard,
		"birthday":        req.Birthday,
		"avatar":          req.Avatar,
		"email":           req.Email,
		"status":          int32(req.Status),
		"remark":          req.Remark,
		"member_level_id": req.MemberLevelId,
		"point":           req.Point,
		"balance":         req.Balance,
		"updated_by":      req.AdminId,
	}
	
	if err := db.Conn.Table(member.GetTableName()).Model(&member.Member{MemberId: req.MemberId}).Updates(aul).Error; err != nil {
		return nil, err
	}
	
	if utils.IsCancelled(ctx) {
		return nil, fmt.Errorf("client cancelled ")
	}
	
	return &basepb.AnyRes{
		Id:    req.MemberId,
		State: 1,
	}, nil
}

func (s *Member) EditMemberStatus(ctx context.Context, req *basepb.EditStatusReq) (*basepb.AnyRes, error) {
	if _, ok := memberpb.MemberStatus_name[req.Status]; !ok {
		return nil, fmt.Errorf("params is err")
	}
	
	if _, err := member.GetOneByMemberId(req.Id); err != nil {
		return nil, err
	}
	
	aul := map[string]interface{}{
		"status":     memberpb.MemberStatus(req.Status),
		"updated_by": req.AdminId,
	}
	
	if err := db.Conn.Table(member.GetTableName()).Model(&member.Member{MemberId: req.Id}).Updates(aul).Error; err != nil {
		return nil, err
	}
	
	if utils.IsCancelled(ctx) {
		return nil, fmt.Errorf("client cancelled ")
	}
	
	return &basepb.AnyRes{
		Id:    req.Id,
		State: 1,
	}, nil
}

func (s *Member) GetMemberList(ctx context.Context, req *memberpb.GetMemberReq) (*memberpb.ListMemberRes, error) {
	var page uint64 = 1
	if req.Page > 0 {
		page = req.Page
	}
	
	var pageSize uint64 = 10
	if req.PageSize > 0 {
		pageSize = req.PageSize
	}
	
	rows, total, err := member.GetMemberList(req.MemberId, req.Status, req.Mobile, page, pageSize)
	if err != nil {
		return nil, err
	}
	
	if utils.IsCancelled(ctx) {
		return nil, fmt.Errorf("client cancelled ")
	}
	
	list := make([]*memberpb.MemberDetail, 0, len(rows))
	for k := range rows {
		list = append(list, &memberpb.MemberDetail{
			MemberId:      rows[k].MemberId,
			Nickname:      rows[k].Nickname,
			Mobile:        rows[k].Mobile,
			Name:          rows[k].Name,
			Gender:        memberpb.MemberGender(rows[k].Gender),
			IdCard:        rows[k].IdCard,
			Birthday:      rows[k].Birthday,
			Avatar:        rows[k].Avatar,
			Email:         rows[k].Email,
			Status:        memberpb.MemberStatus(rows[k].Status),
			Remark:        rows[k].Remark,
			MemberLevelId: rows[k].MemberLevelId,
			Point:         rows[k].Point,
			Balance:       rows[k].Balance,
			CreatedBy:     rows[k].CreatedBy,
			UpdatedBy:     rows[k].UpdatedBy,
			CreatedAt:     rows[k].CreatedAt.Format(utils.TIME_STD_FORMART),
			UpdatedAt:     rows[k].UpdatedAt.Format(utils.TIME_STD_FORMART),
		})
	}
	
	return &memberpb.ListMemberRes{
		Total:   total,
		Members: list,
	}, nil
}
