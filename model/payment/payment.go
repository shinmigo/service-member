package payment

import (
	"fmt"
	
	"goshop/service-member/pkg/db"
	"goshop/service-member/pkg/utils"
)

type Payment struct {
	PaymentId   string
	Money       float64
	MemberId    uint64
	Type        int32
	Status      int32
	PaymentCode string
	Ip          string
	Params      string
	PayedMsg    string
	TradeNo     string
	CreatedAt   utils.JSONTime
	UpdatedAt   utils.JSONTime
}

func GetTableName() string {
	return "payment"
}

func GetField() []string {
	return []string{
		"payment_id", "money", "member_id", "type", "status",
		"payment_code", "payed_msg", "trade_no",
	}
}

func GetOneByPaymentId(paymentId string) (*Payment, error) {
	if len(paymentId) == 0 {
		return nil, fmt.Errorf("payment_id is null")
	}
	row := &Payment{}
	err := db.Conn.Table(GetTableName()).
		Select(GetField()).
		Where("payment_id = ?", paymentId).
		First(row).Error
	
	if err != nil {
		return nil, fmt.Errorf("err: %v", err)
	}
	return row, nil
}
