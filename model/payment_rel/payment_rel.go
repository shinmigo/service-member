package payment_rel

import (
	"bytes"
	"fmt"
	
	"github.com/jinzhu/gorm"
)

// 支持批量支付
type PaymentRel struct {
	PaymentId string
	SourceId  string
	Money     float64
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
