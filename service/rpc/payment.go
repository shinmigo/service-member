package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	
	"github.com/shinmigo/pb/memberpb"
	"github.com/shinmigo/pb/orderpb"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"goshop/service-member/model/payment"
	"goshop/service-member/model/payment_rel"
	"goshop/service-member/pkg/db"
	"goshop/service-member/pkg/grpc/gclient"
	"goshop/service-member/pkg/utils"
)

type Payment struct {
}

func NewPayment() *Payment {
	return &Payment{}
}

func (p *Payment) GetPay(ctx context.Context, req *memberpb.PaymentIdReq) (res *memberpb.PaymentRelList, err error) {
	rows, err := payment_rel.GetAllByPaymentId(req.PaymentId)
	if err != nil {
		return nil, err
	}
	
	list := make([]*memberpb.PaymentParams, 0, len(rows))
	for k := range rows {
		item, _ := jsonLib.Marshal(rows[k])
		buf := &memberpb.PaymentParams{}
		_ = jsonLib.Unmarshal(item, buf)
		list = append(list, buf)
	}
	
	return &memberpb.PaymentRelList{List: list}, nil
}

func (p *Payment) AddPay(ctx context.Context, req *memberpb.ToAdd) (res *memberpb.PaymentRes, err error) {
	if req.Type == 0 {
		err = fmt.Errorf("支付方式错误！")
		return
	}
	
	paymentCode, ok := memberpb.PaymentCode_name[int32(req.PaymentCode)]
	if !ok {
		err = fmt.Errorf("支付类型不存在")
		return
	}
	
	if len(req.Params) == 0 {
		err = fmt.Errorf("支付错误！")
		return
	}
	
	paymentId := utils.GetUniqueId()
	
	paymentIdStr := strconv.FormatUint(paymentId, 10)
	var moeny float64 = 0
	params := make([]*payment_rel.PaymentRel, 0, len(req.Params))
	for k := range req.Params {
		moeny, _ = decimal.NewFromFloat(moeny).Add(decimal.NewFromFloat(req.Params[k].Money)).Float64()
		params = append(params, &payment_rel.PaymentRel{
			PaymentId: paymentIdStr,
			SourceId:  req.Params[k].SourceId,
			Money:     req.Params[k].Money,
		})
	}
	
	var jsonStr []byte
	jsonStr, err = json.Marshal(params)
	if err != nil {
		err = fmt.Errorf("支付错误！")
		return
	}
	
	aul := payment.Payment{
		PaymentId:   paymentIdStr,
		Money:       moeny,
		MemberId:    req.MemberId,
		Type:        int32(req.Type),
		Status:      int32(memberpb.PaymentStatus_Unpaid),
		PaymentCode: paymentCode,
		Ip:          req.Ip,
		Params:      string(jsonStr),
	}
	
	tx := db.Conn.Begin()
	if err = tx.Error; err != nil {
		return
	}
	
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
		
		if err != nil {
			tx.Rollback()
		}
	}()
	
	if err = tx.Table(payment.GetTableName()).Create(&aul).Error; err != nil {
		return
	}
	
	if err = payment_rel.BatchInsert(tx, params); err != nil {
		return
	}
	
	if err = tx.Commit().Error; err != nil {
		return
	}
	
	if ctx.Err() == context.Canceled {
		err = status.Errorf(codes.Canceled, "The client canceled the request")
		return
	}
	
	return &memberpb.PaymentRes{
		PaymentId: paymentIdStr,
		State:     1,
		Money:     moeny,
	}, nil
}

func (p *Payment) EditPay(ctx context.Context, req *memberpb.ToEdit) (*memberpb.PaymentRes, error) {
	info, err := payment.GetOneByPaymentId(req.PaymentId)
	if err != nil {
		return nil, err
	}
	
	if req.Money != info.Money {
		return nil, fmt.Errorf("支付金额不匹配")
	}
	
	if info.Status != int32(memberpb.PaymentStatus_Unpaid) {
		return nil, fmt.Errorf("支付状态出错")
	}
	
	aul := payment.Payment{
		Status:   int32(req.Status),
		PayedMsg: req.PayedMsg,
		TradeNo:  req.TradeNo,
	}
	
	if err := db.Conn.Table(payment.GetTableName()).Where("payment_id = ?", req.PaymentId).Updates(aul).Error; err != nil {
		return nil, err
	}
	
	if req.Status == memberpb.PaymentStatus_PaySuccess { //成功订单触发
		// TODO处理订单
		rows, err := payment_rel.GetAllByPaymentId(req.PaymentId)
		if err != nil {
			return nil, err
		}
		orderIds := make([]uint64, 0, len(rows))
		for k := range rows {
			orderId, _ := strconv.ParseUint(rows[k].SourceId, 0, 64)
			orderIds = append(orderIds, orderId)
		}
		
		r, err := gclient.OrderClient.PayOrder(ctx, &orderpb.PayOrderReq{OrderId: orderIds})
		if err != nil {
			return nil, err
		}
		
		if r.State != 1 {
			return nil, fmt.Errorf("更新订单失败")
		}
	}
	if ctx.Err() == context.Canceled {
		return nil, status.Errorf(codes.Canceled, "The client canceled the request")
	}
	
	return &memberpb.PaymentRes{
		PaymentId: req.PaymentId,
		State:     1,
		Money:     info.Money,
	}, nil
}
