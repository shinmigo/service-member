package member

var StatusMap map[uint32]string

func init() {
	StatusMap = make(map[uint32]string, 2)

	StatusMap[0] = "冻结"
	StatusMap[1] = "正常"
}

type Member struct {
	MemberId      uint64 `json:"member_id"`
	Nickname      string `json:"nickname"`
	Status        uint32 `json:"status"`
	Mobile        string `json:"mobile"`
	MemberLevelId string `json:"member_level_id"`
	Gender        uint32 `json:"gender"`
	Birthday      string `json:"birthday"`
	CreatedBy     string `json:"created_by"`
	UpdatedBy     string `json:"updated_by"`
}

func GetTableName() string {
	return "member"
}

func GetDetailFields() []string {
	return []string{
		"member_id", "nickname", "mobile", "register_ip", "name", "gender", "id_card", "birthday", "avatar", "email", "status", "register_time", "last_login_time", "remark", "member_level_id", "point", "balance", "created_at", "created_by", "updated_at", "updated_by",
	}
}
