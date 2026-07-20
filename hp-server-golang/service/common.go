package service

import (
	"hp-server-lib/db"
	"hp-server-lib/entity"

	"gorm.io/gorm"
)

func PaginateQuery(model interface{}, userId int, page int, pageSize int, results interface{}) (int64, error) {
	var total int64
	query := db.DB.Model(model)
	if userId >= 0 {
		query = query.Where("user_id = ?", userId)
	}
	err := query.Order("id desc").Count(&total).Offset((page - 1) * pageSize).Limit(pageSize).Find(results).Error
	return total, err
}

func PaginateWithQuery(query *gorm.DB, page int, pageSize int, results interface{}) (int64, error) {
	var total int64
	err := query.Order("id desc").Count(&total).Offset((page - 1) * pageSize).Limit(pageSize).Find(results).Error
	return total, err
}

func GetUserMap(userIds []int) map[int]*entity.UserCustomEntity {
	userMap := make(map[int]*entity.UserCustomEntity)
	if len(userIds) == 0 {
		return userMap
	}
	var users []*entity.UserCustomEntity
	if err := db.DB.Model(&entity.UserCustomEntity{}).Where("id IN ?", userIds).Find(&users).Error; err == nil {
		for _, user := range users {
			if user.Id != nil {
				userMap[*user.Id] = user
			}
		}
	}
	return userMap
}

func GetConfigMap(configIds []int) map[int]*entity.UserConfigEntity {
	configMap := make(map[int]*entity.UserConfigEntity)
	if len(configIds) == 0 {
		return configMap
	}
	var configItems []*entity.UserConfigEntity
	if err := db.DB.Model(&entity.UserConfigEntity{}).Where("id IN ?", configIds).Find(&configItems).Error; err == nil {
		for _, conf := range configItems {
			if conf.Id != nil {
				configMap[*conf.Id] = conf
			}
		}
	}
	return configMap
}
