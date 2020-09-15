package rpc

import (
	"context"
	"fmt"
	
	"github.com/shinmigo/pb/basepb"
	"github.com/shinmigo/pb/memberpb"
	"goshop/service-member/model/cart"
	"goshop/service-member/pkg/db"
	"goshop/service-member/pkg/utils"
)

type Cart struct {
}

func NewCart() *Cart {
	return &Cart{}
}

func (c *Cart) AddCart(ctx context.Context, req *memberpb.AddCartReq) (*basepb.AnyRes, error) {
	row := &cart.Cart{}
	exist := db.Conn.Table(cart.GetTableName()).Select(cart.GetField()).
		Where("member_id = ? and product_id = ? and product_spec_id = ?", req.MemberId, req.ProductId, req.ProductSpecId).
		First(row).RowsAffected
	
	// 连续点击增加购物车，则更新数量
	if exist > 0 {
		cartId, err := cart.SetCartNums(row, int32(req.IsSelect), req.Nums, req.IsPlus)
		if err != nil {
			return nil, err
		}
		
		if utils.IsCancelled(ctx) {
			return nil, fmt.Errorf("client cancelled ")
		}
		
		return &basepb.AnyRes{
			Id:    cartId,
			State: 1,
		}, nil
	}
	
	// 新增
	if req.Nums < 0 {
		req.Nums = 0
	}
	
	aul := cart.Cart{
		MemberId:      req.MemberId,
		ProductId:     req.ProductId,
		ProductSpecId: req.ProductSpecId,
		IsSelect:      int32(req.IsSelect),
		Nums:          uint64(req.Nums),
	}
	
	if err := db.Conn.Table(cart.GetTableName()).Create(&aul).Error; err != nil {
		return nil, err
	}
	
	if utils.IsCancelled(ctx) {
		return nil, fmt.Errorf("client cancelled ")
	}
	
	return &basepb.AnyRes{
		Id:    aul.CartId,
		State: 1,
	}, nil
}

func (c *Cart) DelCart(ctx context.Context, req *memberpb.DelCartReq) (*basepb.AnyRes, error) {
	if req.IsAll == 1 {
		if len(req.CartIds) <= 0 {
			return nil, fmt.Errorf("请选择商品 ")
		}
		
		if err := db.Conn.Table(cart.GetTableName()).
			Where("cart_id in (?) and member_id = ?", req.CartIds, req.MemberId).
			Delete(cart.Cart{}).Error; err != nil {
			return nil, err
		}
	} else {
		// 清空购物车
		if err := db.Conn.Table(cart.GetTableName()).
			Where("member_id = ?", req.CartIds, req.MemberId).
			Delete(cart.Cart{}).Error; err != nil {
			return nil, err
		}
	}
	
	if utils.IsCancelled(ctx) {
		return nil, fmt.Errorf("client cancelled ")
	}
	
	// 批量删除返回0
	return &basepb.AnyRes{
		Id:    0,
		State: 1,
	}, nil
}

func (c *Cart) GetCartListByMemberId(ctx context.Context, req *memberpb.ListCartReq) (*memberpb.ListCartRes, error) {
	rows := make([]*cart.Cart, 0, 32)
	if err := db.Conn.Table(cart.GetTableName()).Select(cart.GetField()).
		Where("member_id = ?", req.MemberId).Find(&rows).Error; err != nil {
		return nil, err
	}
	
	if utils.IsCancelled(ctx) {
		return nil, fmt.Errorf("client cancelled ")
	}
	
	list := make([]*memberpb.CartDetail, 0, len(rows))
	for k := range rows {
		temp, _ := jsonLib.Marshal(rows[k])
		buf := &memberpb.CartDetail{}
		_ = jsonLib.Unmarshal(temp, buf)
		
		list = append(list, buf)
	}
	return &memberpb.ListCartRes{
		Carts: list,
	}, nil
}

func (c *Cart) SelectCart(ctx context.Context, req *memberpb.SelectCartReq) (*basepb.AnyRes, error) {
	total := len(req.SelectCart)
	if total <= 0 {
		return nil, fmt.Errorf("请选择商品 ")
	}
	
	checkedTrue := make([]uint64, 0, total)
	checkedFalse := make([]uint64, 0, total)
	for k := range req.SelectCart {
		if req.SelectCart[k].IsSelect == memberpb.CartIsSelect_Select {
			checkedTrue = append(checkedTrue, req.SelectCart[k].CartId)
		} else {
			checkedFalse = append(checkedFalse, req.SelectCart[k].CartId)
		}
	}
	
	if len(checkedTrue) > 0 {
		if err := db.Conn.Table(cart.GetTableName()).
			Where("cart_id in (?) and member_id = ?", checkedTrue, req.MemberId).
			Update("is_select", 1).Error; err != nil {
			return nil, err
		}
	}
	
	if len(checkedFalse) > 0 {
		if err := db.Conn.Table(cart.GetTableName()).
			Where("cart_id in (?) and member_id = ?", checkedFalse, req.MemberId).
			Update("is_select", 0).Error; err != nil {
			return nil, err
		}
	}
	
	if utils.IsCancelled(ctx) {
		return nil, fmt.Errorf("client cancelled ")
	}
	
	// 批量更新返回0
	return &basepb.AnyRes{
		Id:    0,
		State: 1,
	}, nil
}
