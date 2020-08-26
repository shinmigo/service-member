package rpc

import (
	"context"
	"fmt"
	
	"github.com/shinmigo/pb/basepb"
	"github.com/shinmigo/pb/memberpb"
	"goshop/service-member/model/address"
	"goshop/service-member/pkg/db"
	"goshop/service-member/pkg/utils"
)

type Address struct {
}

func NewAddress() *Address {
	return &Address{}
}

func (s *Address) AddAddress(ctx context.Context, req *memberpb.Address) (*basepb.AnyRes, error) {
	aul := address.Address{
		MemberId:      req.MemberId,
		Category:      int32(req.Category),
		Name:          req.Name,
		Mobile:        req.Mobile,
		Tel:           req.Tel,
		Postcode:      req.Postcode,
		AddressDetail: req.Address,
		CodeProv:      req.CodeProv,
		CodeCity:      req.CodeCity,
		CodeCoun:      req.CodeCoun,
		CodeTown:      req.CodeTown,
		IsDefault:     int32(req.IsDefault),
		Longitude:     req.Longitude,
		Latitude:      req.Latitude,
		CreatedBy:     req.AdminId,
		UpdatedBy:     req.AdminId,
	}
	
	if err := db.Conn.Table(address.GetTableName()).Create(&aul).Error; err != nil {
		return nil, err
	}
	
	if utils.IsCancelled(ctx) {
		return nil, fmt.Errorf("client cancelled ")
	}
	
	return &basepb.AnyRes{
		Id:    aul.AddressId,
		State: 1,
	}, nil
}

func (s *Address) EditAddress(ctx context.Context, req *memberpb.Address) (*basepb.AnyRes, error) {
	info, err := address.GetOneByAddressId(req.AddressId)
	if err != nil {
		return nil, err
	}
	
	if req.IsDefault == memberpb.AddressIsDefault_Used {
		// 把其他的默认设置为0
		db.Conn.Table(address.GetTableName()).
			Where("member_id = ? and is_default = 1", info.MemberId).
			Update("is_default", 0)
	}
	
	aul := map[string]interface{}{
		"category":       int32(req.Category),
		"name":           req.Name,
		"mobile":         req.Mobile,
		"tel":            req.Tel,
		"postcode":       req.Postcode,
		"address_detail": req.Address,
		"code_prov":      req.CodeProv,
		"code_city":      req.CodeCity,
		"code_coun":      req.CodeCoun,
		"code_town":      req.CodeTown,
		"is_default":     int32(req.IsDefault),
		"longitude":      req.Longitude,
		"latitude":       req.Latitude,
		"updated_by":     req.AdminId,
	}
	
	if err := db.Conn.Table(address.GetTableName()).Model(&address.Address{AddressId: req.AddressId}).Updates(aul).Error; err != nil {
		return nil, err
	}
	
	if utils.IsCancelled(ctx) {
		return nil, fmt.Errorf("client cancelled ")
	}
	
	return &basepb.AnyRes{
		Id:    req.AddressId,
		State: 1,
	}, nil
}

func (s *Address) DelAddress(ctx context.Context, req *basepb.DelReq) (*basepb.AnyRes, error) {
	if _, err := address.GetOneByAddressId(req.Id); err != nil {
		return nil, err
	}
	
	if err := db.Conn.Table(address.GetTableName()).Delete(&address.Address{AddressId: req.Id}).Error; err != nil {
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

func (s *Address) GetAddressListByMemberId(ctx context.Context, req *memberpb.ListAddressReq) (*memberpb.ListAddressRes, error) {
	var page uint64 = 1
	if req.Page > 0 {
		page = req.Page
	}
	
	var pageSize uint64 = 10
	if req.PageSize > 0 {
		pageSize = req.PageSize
	}
	rows, total, err := address.GetAddressListByMemberId(req.MemberId, page, pageSize)
	if err != nil {
		return nil, err
	}
	
	if utils.IsCancelled(ctx) {
		return nil, fmt.Errorf("client cancelled ")
	}
	
	list := make([]*memberpb.AddressDetail, 0, len(rows))
	for k := range rows {
		list = append(list, &memberpb.AddressDetail{
			AddressId: rows[k].AddressId,
			MemberId:  rows[k].MemberId,
			Category:  memberpb.AddressCategory(rows[k].Category),
			Name:      rows[k].Name,
			Mobile:    rows[k].Mobile,
			Tel:       rows[k].Tel,
			Postcode:  rows[k].Postcode,
			Address:   rows[k].AddressDetail,
			CodeProv:  rows[k].CodeProv,
			CodeCity:  rows[k].CodeCity,
			CodeCoun:  rows[k].CodeCoun,
			CodeTown:  rows[k].CodeTown,
			IsDefault: memberpb.AddressIsDefault(rows[k].IsDefault),
			Longitude: rows[k].Longitude,
			Latitude:  rows[k].Latitude,
			CreatedBy: rows[k].CreatedBy,
			UpdatedBy: rows[k].UpdatedBy,
			CreatedAt: rows[k].CreatedAt.Format(utils.TIME_STD_FORMART),
			UpdatedAt: rows[k].UpdatedAt.Format(utils.TIME_STD_FORMART),
		})
	}
	
	return &memberpb.ListAddressRes{
		Total:     total,
		Addresses: list,
	}, nil
}

func (s *Address) GetAddressDetail(ctx context.Context, req *basepb.GetOneReq) (*memberpb.AddressDetail, error) {
	row, err := address.GetOneByAddressId(req.Id)
	if err != nil {
		return nil, err
	}
	if utils.IsCancelled(ctx) {
		return nil, fmt.Errorf("client cancelled ")
	}
	
	return &memberpb.AddressDetail{
		AddressId: row.AddressId,
		MemberId:  row.MemberId,
		Category:  memberpb.AddressCategory(row.Category),
		Name:      row.Name,
		Mobile:    row.Mobile,
		Tel:       row.Tel,
		Postcode:  row.Postcode,
		Address:   row.AddressDetail,
		CodeProv:  row.CodeProv,
		CodeCity:  row.CodeCity,
		CodeCoun:  row.CodeCoun,
		CodeTown:  row.CodeTown,
		IsDefault: memberpb.AddressIsDefault(row.IsDefault),
		Longitude: row.Longitude,
		Latitude:  row.Latitude,
		CreatedBy: row.CreatedBy,
		UpdatedBy: row.UpdatedBy,
		CreatedAt: row.CreatedAt.Format(utils.TIME_STD_FORMART),
		UpdatedAt: row.UpdatedAt.Format(utils.TIME_STD_FORMART),
	}, nil
}
