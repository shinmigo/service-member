package rpc

import (
	"fmt"

	"goshop/service-member/model/member"
	"goshop/service-member/pkg/db"
	"goshop/service-member/pkg/utils"

	"golang.org/x/crypto/bcrypt"

	"github.com/shinmigo/pb/basepb"
	"github.com/shinmigo/pb/memberpb"
	"golang.org/x/net/context"
)

type Member struct {
}

func NewMember() *Member {
	return &Member{}
}

func (s *Member) RegisterByMobile(ctx context.Context, req *memberpb.MobilePasswdReq) (*memberpb.LoginRes, error) {
	row, _ := member.GetOneByMobile(req.Mobile)
	if row != nil {
		return nil, fmt.Errorf("手机号已被注册")
	}

	aul := member.Member{
		Mobile:     req.Mobile,
		Password:   req.Password,
		RegisterIp: req.Ip,
		Status:     int32(memberpb.MemberStatus_Normal),
	}

	if err := db.Conn.Create(&aul).Error; err != nil {
		return nil, err
	}

	if utils.IsCancelled(ctx) {
		return nil, fmt.Errorf("client cancelled ")
	}

	return login(&aul)
}

func (s *Member) LoginByMobile(ctx context.Context, req *memberpb.MobilePasswdReq) (*memberpb.LoginRes, error) {
	row, _ := member.GetOneByMobile(req.Mobile)
	if row == nil {
		return nil, fmt.Errorf("手机号不存在，请注册")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(row.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("密码不正确")
	}

	if memberpb.MemberStatus(row.Status) != memberpb.MemberStatus_Normal {
		return nil, fmt.Errorf("该账号不可使用")
	}

	if utils.IsCancelled(ctx) {
		return nil, fmt.Errorf("client cancelled ")
	}

	return login(row)
}

func (s *Member) GetMemberForLogin(ctx context.Context, req *memberpb.MemberIdReq) (*memberpb.LoginRes, error) {
	row, err := member.GetOneByMemberId(req.MemberId)
	if err != nil {
		return nil, err
	}

	return login(row)
}

func login(m *member.Member) (*memberpb.LoginRes, error) {
	if m == nil {
		return nil, nil
	}

	if memberpb.MemberStatus(m.Status) != memberpb.MemberStatus_Normal {
		return nil, fmt.Errorf("该账号不可使用")
	}

	return &memberpb.LoginRes{
		MemberId:      m.MemberId,
		Nickname:      m.Nickname,
		Mobile:        m.Mobile,
		Name:          m.Name,
		Gender:        memberpb.MemberGender(m.Gender),
		IdCard:        m.IdCard,
		Birthday:      m.Birthday,
		Avatar:        m.Avatar,
		Email:         m.Email,
		MemberLevelId: m.MemberLevelId,
		Point:         m.Point,
		Balance:       m.Balance,
	}, nil
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

	/*aul := map[string]interface{}{
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
	}*/

	if err := db.Conn.Table(member.GetTableName()).Model(&member.Member{MemberId: req.MemberId}).Updates(&req).Error; err != nil {
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

	aul := map[string]interface{}{
		"status":     memberpb.MemberStatus(req.Status),
		"updated_by": req.AdminId,
	}

	if err := db.Conn.Table(member.GetTableName()).Where("member_id in (?)", req.Id).Updates(aul).Error; err != nil {
		return nil, err
	}

	if utils.IsCancelled(ctx) {
		return nil, fmt.Errorf("client cancelled ")
	}

	return &basepb.AnyRes{
		Id:    req.Id[0],
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

	rows, total, err := member.GetMemberList(req.MemberId, req.Status, req.Mobile, page, pageSize, req.Nickname)
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

func (s *Member) GetMemberDetail(ctx context.Context, req *basepb.GetOneReq) (*memberpb.MemberDetail, error) {
	row, err := member.GetOneByMemberId(req.Id)
	if err != nil {
		return nil, err
	}
	if utils.IsCancelled(ctx) {
		return nil, fmt.Errorf("client cancelled ")
	}

	return &memberpb.MemberDetail{
		MemberId:      row.MemberId,
		Nickname:      row.Nickname,
		Mobile:        row.Mobile,
		Name:          row.Name,
		Gender:        memberpb.MemberGender(row.Gender),
		IdCard:        row.IdCard,
		Birthday:      row.Birthday,
		Avatar:        row.Avatar,
		Email:         row.Email,
		Status:        memberpb.MemberStatus(row.Status),
		Remark:        row.Remark,
		MemberLevelId: row.MemberLevelId,
		Point:         row.Point,
		Balance:       row.Balance,
		CreatedBy:     row.CreatedBy,
		UpdatedBy:     row.UpdatedBy,
		CreatedAt:     row.CreatedAt.Format(utils.TIME_STD_FORMART),
		UpdatedAt:     row.UpdatedAt.Format(utils.TIME_STD_FORMART),
	}, nil
}
