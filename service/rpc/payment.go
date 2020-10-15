package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	
	"github.com/shinmigo/pb/basepb"
	"github.com/shinmigo/pb/memberpb"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"goshop/service-member/model/payment"
	"goshop/service-member/model/payment_rel"
	"goshop/service-member/pkg/db"
	"goshop/service-member/pkg/utils"
)

type Payment struct {
}

func NewPayment() *Payment {
	return &Payment{}
}

func (p *Payment) AddPay(ctx context.Context, req *memberpb.ToAdd) (res *basepb.AnyRes, err error) {
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
	
	return &basepb.AnyRes{
		Id:    paymentId,
		State: 1,
	}, nil
}

func (p *Payment) EditPay(ctx context.Context, req *memberpb.ToEdit) (res *basepb.AnyRes, err error) {
	var info *payment.Payment
	info, err = payment.GetOneByPaymentId(req.PaymentId)
	if err != nil {
		return
	}
	
	if req.Money != info.Money {
		err = fmt.Errorf("支付金额不匹配")
		return
	}
	
	if info.Status != int32(memberpb.PaymentStatus_Unpaid) {
		err = fmt.Errorf("支付状态出错")
		return
	}
	
	aul := payment.Payment{
		Status:   int32(memberpb.PaymentStatus_PaySuccess),
		PayedMsg: req.PayedMsg,
		TradeNo:  req.TradeNo,
	}
	
	if err = db.Conn.Table(payment.GetTableName()).Where("payment_id = ?", req.PaymentId).Updates(aul).Error; err != nil {
		return
	}
	
	if ctx.Err() == context.Canceled {
		err = status.Errorf(codes.Canceled, "The client canceled the request")
		return
	}
	
	paymentId, _ := strconv.ParseUint(req.PaymentId, 10, 64)
	return &basepb.AnyRes{
		Id:    paymentId,
		State: 1,
	}, nil
}
