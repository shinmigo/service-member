package account

type Account struct {
	AccountId   uint64 `json:"account_id"`
	MemberId    uint64 `json:"member_id"`
	AccountName string `json:"account_name"`
	Category    uint8  `json:"category"`
	Password    string `json:"password"`
	Status      uint32 `json:"status"`
}

func GetTableName() string {
	return "account"
}
