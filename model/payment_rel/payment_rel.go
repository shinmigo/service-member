package payment_rel

import (
	"bytes"
	"fmt"
	
	"github.com/jinzhu/gorm"
	"goshop/service-member/pkg/db"
)

// 支持批量支付
type PaymentRel struct {
	PaymentId string
	SourceId  string  `json:"source_id"`
	Money     float64 `json:"money"`
}

func BatchInsert(db *gorm.DB, rel []*PaymentRel) error {
	var buf bytes.Buffer
	sql := "INSERT INTO payment_rel (payment_id, source_id, money) VALUES "
	if _, err := buf.WriteString(sql); err != nil {
		return err
	}
	
	for k := range rel {
		if k == len(rel)-1 {
			buf.WriteString(fmt.Sprintf("('%s', '%s', %f);",
				rel[k].PaymentId,
				rel[k].SourceId,
				rel[k].Money,
			))
		} else {
			buf.WriteString(fmt.Sprintf("('%s', '%s', %f),",
				rel[k].PaymentId,
				rel[k].SourceId,
				rel[k].Money,
			))
		}
	}
	return db.Exec(buf.String()).Error
}

func GetAllByPaymentId(paymentId string) ([]*PaymentRel, error) {
	if len(paymentId) == 0 {
		return nil, fmt.Errorf("payment_id is null")
	}
	rows := make([]*PaymentRel, 0, 8)
	err := db.Conn.Table("payment_rel").
		Select([]string{"source_id", "money"}).
		Where("payment_id = ?", paymentId).
		Find(&rows).Error
	
	if err != nil {
		return nil, fmt.Errorf("err: %v", err)
	}
	return rows, nil
}
