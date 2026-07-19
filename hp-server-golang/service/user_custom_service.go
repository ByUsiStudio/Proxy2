package service

import (
	"hp-server-lib/bean"
	"hp-server-lib/db"
	"hp-server-lib/entity"
	"time"
)

type UserCustomService struct {
}

func (receiver *UserCustomService) AddData(custom entity.UserCustomEntity) error {
	if custom.Id == nil {
		custom.CreateTime = time.Now()
		if err := db.DB.Create(&custom).Error; err != nil {
			return err
		}
	} else {
		if custom.Password == "" {
			if err := db.DB.Model(&entity.UserCustomEntity{}).Where("id = ?", custom.Id).Updates(map[string]interface{}{
				"username": custom.Username,
				"email":    custom.Email,
				"desc":     custom.Desc,
				"status":   custom.Status,
			}).Error; err != nil {
				return err
			}
		} else {
			if err := db.DB.Save(&custom).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func (receiver *UserCustomService) ListData(page int, pageSize int) *bean.ResPage {
	var results []entity.UserCustomEntity
	var total int64
	// 计算总记录数并执行分页查询
	db.DB.Model(&entity.UserCustomEntity{}).Order("id desc").Count(&total).Offset((page - 1) * pageSize).Limit(pageSize).Find(&results)
	return bean.PageOk(total, results)
}

func (receiver *UserCustomService) RemoveData(id int) {
	db.DB.Delete(&entity.UserCustomEntity{Id: &id})
}
