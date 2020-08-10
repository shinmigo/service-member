package utils

// 获取分页的offset
func GetPageOffset(pageSize uint64, page uint64) uint64 {
	// 默认从0开始
	if page == 0 {
		return 0
	}

	return (page - 1) * pageSize
}
