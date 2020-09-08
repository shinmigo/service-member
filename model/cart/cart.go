package cart

import (
	"fmt"
	
	"github.com/jinzhu/gorm"
	"goshop/service-member/pkg/db"
	"goshop/service-member/pkg/utils"
)

type Cart struct {
	CartId        uint64         `json:"cart_id" gorm:"PRIMARY_KEY"`
	MemberId      uint64         `json:"member_id"`
	ProductId     uint64         `json:"product_id"`
	ProductSpecId uint64         `json:"product_spec_id"`
	IsSelect      int32          `json:"is_select"`
	Nums          uint64         `json:"nums"`
	CreatedAt     utils.JSONTime `json:"created_at"`
	UpdatedAt     utils.JSONTime `json:"updated_at"`
}

func GetTableName() string {
	return "cart"
}

func GetField() []string {
	return []string{
		"cart_id", "member_id", "product_id", "product_spec_id", "is_select", "nums", "created_at", "updated_at",
	}
}

func GetOneByCartId(cartId uint64) (*Cart, error) {
	if cartId == 0 {
		return nil, fmt.Errorf("cart_id is null")
	}
	row := &Cart{}
	err := db.Conn.Table(GetTableName()).
		Select(GetField()).
		Where("cart_id = ?", cartId).
		First(row).Error
	
	if err != nil {
		return nil, fmt.Errorf("err: %v", err)
	}
	return row, nil
}

func SetCartNums(row *Cart, isSelect int32, nums int64, isPlus bool) (cartId uint64, err error) {
	buf := make(map[string]interface{})
	buf["is_select"] = isSelect
	if isPlus {
		if int64(row.Nums)+nums < 0 { // 5+(-6)=-1
			buf["nums"] = 0
		} else {
			buf["nums"] = gorm.Expr("nums + ?", nums)
		}
	} else {
		if nums < 0 {
			nums = 0
		}
		buf["nums"] = nums
	}
	
	if err = db.Conn.Table(GetTableName()).Model(row).Updates(buf).Error; err != nil {
		return
	}
	
	cartId = row.CartId
	return
}
