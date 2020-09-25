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
		Name:          req.Name,
		Mobile:        req.Mobile,
		AddressDetail: req.Address,
		CodeProv:      req.CodeProv,
		CodeCity:      req.CodeCity,
		CodeCoun:      req.CodeCoun,
		CodeTown:      req.CodeTown,
		RoomNumber:    req.RoomNumber,
		IsDefault:     int32(req.IsDefault),
		Longitude:     req.Longitude,
		Latitude:      req.Latitude,
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
	
	if info.MemberId != req.MemberId {
		return nil, fmt.Errorf("different member")
	}
	
	if req.IsDefault == memberpb.AddressIsDefault_Used {
		// 把其他的默认设置为0
		db.Conn.Table(address.GetTableName()).
			Where("member_id = ? and is_default = 1", info.MemberId).
			Update("is_default", 0)
	}
	
	aul := map[string]interface{}{
		"name":           req.Name,
		"mobile":         req.Mobile,
		"address_detail": req.Address,
		"room_number":    req.RoomNumber,
		"code_prov":      req.CodeProv,
		"code_city":      req.CodeCity,
		"code_coun":      req.CodeCoun,
		"code_town":      req.CodeTown,
		"is_default":     int32(req.IsDefault),
		"longitude":      req.Longitude,
		"latitude":       req.Latitude,
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
	
	list := make([]*memberpb.Address, 0, len(rows))
	for k := range rows {
		item, _ := jsonLib.Marshal(rows[k])
		buf := &memberpb.Address{}
		_ = jsonLib.Unmarshal(item, buf)
		buf.Address = rows[k].AddressDetail
		list = append(list, buf)
	}
	
	return &memberpb.ListAddressRes{
		Total:     total,
		Addresses: list,
	}, nil
}

func (s *Address) GetAddressDetail(ctx context.Context, req *basepb.GetOneReq) (*memberpb.Address, error) {
	row, err := address.GetOneByAddressId(req.Id)
	if err != nil {
		return nil, err
	}
	if utils.IsCancelled(ctx) {
		return nil, fmt.Errorf("client cancelled ")
	}
	
	item, _ := jsonLib.Marshal(row)
	buf := &memberpb.Address{}
	_ = jsonLib.Unmarshal(item, buf)
	buf.Address = row.AddressDetail
	
	return buf, nil
}
