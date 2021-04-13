package model

import (
	"github.com/Zhenghao-Liu/OAuth_demo/common"
	"github.com/gin-gonic/gin"
	"sync"
)

type OAuthInfo struct {
	AppName     string `gorm:"column:app_name" json:"app_name"`        // 应用名称
	Homepage    string `gorm:"column:homepage" json:"homepage"`        // 首页
	Description string `gorm:"column:description;" json:"description"` // 描述简介
	Callback    string `gorm:"column:callback" json:"callback"`        // 重定向URL
	AppID       string `gorm:"column:app_id" json:"app_id"`            // 客户端标识
	AppSecret   string `gorm:"column:app_secret" json:"app_secret"`    // 客户端密钥
	BaseInfo
}

type OAuthInfoDao struct{}

var oauthInfoDaoOnce sync.Once
var oauthInfoDaoInstance *OAuthInfoDao

func OAuthInfoDaoInstance() *OAuthInfoDao {
	oauthInfoDaoOnce.Do(func() {
		oauthInfoDaoInstance = &OAuthInfoDao{}
	})
	return oauthInfoDaoInstance
}

func (*OAuthInfo) TableName() string {
	return common.OAuthInfoTableName
}

func (d *OAuthInfoDao) CreateOAuthInfo(ctx *gin.Context, instance *OAuthInfo) error {
	return OAuthDemoDB.Table(common.OAuthInfoTableName).
		Create(instance).Error
}

func (d *OAuthInfoDao) GetByAppID(ctx *gin.Context, appID string) (*OAuthInfo, error) {
	instance := &OAuthInfo{}
	err := OAuthDemoDB.Table(common.OAuthInfoTableName).
		Where("app_id = ?", appID).Where("is_delete = 0").
		First(&instance).Error
	if err != nil {
		return nil, err
	}
	return instance, nil
}

func (d *OAuthInfoDao) UpdateOAuthInfo(ctx *gin.Context, instance *OAuthInfo) error {
	return OAuthDemoDB.Table(common.OAuthInfoTableName).
		Where("app_id = ?", instance.AppID).Where("is_delete = 0").
		Select([]string{"app_name", "homepage", "description", "callback"}).
		Updates(instance).Error
}
