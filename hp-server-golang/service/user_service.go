package service

import (
	"errors"
	"hp-server-lib/bean"
	"hp-server-lib/config"
	"hp-server-lib/db"
	"hp-server-lib/entity"
	"hp-server-lib/util"
	"strings"
	"time"
)

type UserService struct {
}

func (receiver *UserService) Login(login bean.ReqLogin) *bean.ResLoginUser {
	if strings.Compare(config.ConfigData.Admin.Username, login.Email) == 0 && strings.Compare(config.ConfigData.Admin.Password, login.Password) == 0 {
		return bean.NewAdminUser(login)
	} else {
		userQuery := entity.UserCustomEntity{}
		db.DB.Where("username = ? and password = ?", login.Email, login.Password).First(&userQuery)
		if userQuery.Id != nil {
			if userQuery.Status != 1 {
				return nil
			}
			return bean.NewClientUser(*userQuery.Id, userQuery.Username)
		}
	}
	return nil
}

func (receiver *UserService) Register(register bean.ReqRegister) error {
	if !config.ConfigData.System.OpenRegister {
		return errors.New("当前未开启公开注册")
	}

	if len(register.Username) < 3 || len(register.Username) > 50 {
		return errors.New("用户名长度需在3-50个字符之间")
	}

	if len(register.Password) < 6 || len(register.Password) > 50 {
		return errors.New("密码长度需在6-50个字符之间")
	}

	var count int64
	db.DB.Model(&entity.UserCustomEntity{}).Where("username = ?", register.Username).Count(&count)
	if count > 0 {
		return errors.New("用户名已存在")
	}

	if register.Email != "" {
		db.DB.Model(&entity.UserCustomEntity{}).Where("email = ?", register.Email).Count(&count)
		if count > 0 {
			return errors.New("邮箱已被使用")
		}
	}

	status := 1
	if config.ConfigData.System.RegisterReview {
		status = 0
	}

	user := entity.UserCustomEntity{
		Username:   register.Username,
		Password:   register.Password,
		Email:      register.Email,
		Desc:       register.Desc,
		Status:     status,
		CreateTime: time.Now(),
	}

	result := db.DB.Create(&user)
	if result.Error != nil {
		return result.Error
	}

	if config.ConfigData.System.Smtp.Enabled && user.Email != "" {
		go util.SendRegisterEmail(user.Email, user.Username)
	}

	return nil
}

func (receiver *UserService) UpdateStatus(userId int, status int) error {
	var user entity.UserCustomEntity
	result := db.DB.Where("id = ?", userId).First(&user)
	if result.Error != nil {
		return result.Error
	}
	if user.Id == nil {
		return errors.New("用户不存在")
	}
	if err := db.DB.Model(&entity.UserCustomEntity{}).Where("id = ?", userId).UpdateColumn("status", status).Error; err != nil {
		return err
	}
	return nil
}

func (receiver *UserService) GetPublicConfig() map[string]interface{} {
	return map[string]interface{}{
		"siteTitle":    config.ConfigData.System.SiteTitle,
		"openRegister": config.ConfigData.System.OpenRegister,
	}
}

func (receiver *UserService) GetSystemConfig() *config.SystemConfig {
	return &config.ConfigData.System
}

func (receiver *UserService) UpdateSystemConfig(sysConfig config.SystemConfig) error {
	config.ConfigData.System = sysConfig
	return config.SaveConfig()
}
