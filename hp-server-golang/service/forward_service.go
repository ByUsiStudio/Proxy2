package service

import (
	"hp-server-lib/bean"
	"hp-server-lib/db"
	"hp-server-lib/entity"
	"hp-server-lib/ext"
	"hp-server-lib/ext/forward"
	"hp-server-lib/log"
	"sync"
)

var FORWARD_CACHE = sync.Map{}

type ForwardService struct {
}

func InitForward() {
	page := 1
	pageSize := 100
	for {
		var results []*entity.UserFwdEntity
		tx := db.DB.Model(&entity.UserFwdEntity{}).
			Offset((page - 1) * pageSize).
			Limit(pageSize).
			Find(&results)
		if tx.Error != nil {
			break
		}
		// 如果本页没有数据，说明结束
		if len(results) == 0 {
			break
		}
		// 放入缓存
		for _, r := range results {
			start(*r)
		}
		// 下一页
		page++
	}
}

func (receiver *ForwardService) AddData(custom entity.UserFwdEntity) error {
	if custom.Id != nil {
		value, ok := FORWARD_CACHE.Load(*custom.Id)
		if ok {
			proxy := value.(forward.ForwardProxy)
			proxy.Stop()
		}
	}
	tx := db.DB.Save(&custom)
	if tx.Error != nil {
		return tx.Error
	}
	start(custom)
	return nil
}

func start(custom entity.UserFwdEntity) {
	if *custom.Type == "1" && *custom.Status == "1" {
		server := ext.NewHttpFwdServer(*custom.Port, *custom.User, *custom.Pwd)
		start := server.Start(func() {
			FORWARD_CACHE.Delete(*custom.Id)
		})
		if start {
			FORWARD_CACHE.Store(*custom.Id, server)
		}
	}
	if *custom.Type == "2" && *custom.Status == "1" {
		server := ext.NewSocks(*custom.Port, *custom.User, *custom.Pwd)
		start := server.Start(func() {
			FORWARD_CACHE.Delete(*custom.Id)
		})
		if start {
			FORWARD_CACHE.Store(*custom.Id, server)
		}
	}
}

func (receiver *ForwardService) ListData(userId int, page int, pageSize int) *bean.ResPage {
	var results []*entity.UserFwdEntity
	total, _ := PaginateQuery(&entity.UserFwdEntity{}, userId, page, pageSize, &results)
	for _, item := range results {
		_, ok := FORWARD_CACHE.Load(*item.Id)
		if ok {
			item.Tips = "正常"
		} else {
			item.Tips = "已停止"
		}
	}
	if userId < 0 {
		var userIds []int
		for _, item := range results {
			userIds = append(userIds, *item.UserId)
		}
		userMap := GetUserMap(userIds)
		for _, item := range results {
			customEntity := userMap[*item.UserId]
			if customEntity != nil {
				item.Username = customEntity.Username
				item.UserDesc = customEntity.Desc
			}

		}
	}
	return bean.PageOk(total, results)
}

func (receiver *ForwardService) RemoveData(id int) {
	if err := db.DB.Delete(&entity.UserFwdEntity{Id: &id}).Error; err != nil {
		log.Errorf("删除转发配置失败: %v", err)
	}
	value, ok := FORWARD_CACHE.Load(id)
	if ok {
		proxy := value.(forward.ForwardProxy)
		proxy.Stop()
	}
}
