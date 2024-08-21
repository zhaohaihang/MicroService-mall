package utils

import "gorm.io/gorm"

// Paginate 分页查询
func Paginate(page int, pageSize int) func(db *gorm.DB) *gorm.DB {  // 返回类型必须是func(*gorm.DB) *gorm.DB
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
