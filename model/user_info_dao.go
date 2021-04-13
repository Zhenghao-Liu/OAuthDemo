package model

import (
	"github.com/Zhenghao-Liu/OAuth_demo/common"
	"github.com/gin-gonic/gin"
	"strings"
	"sync"
)

type UserInfo struct {
	Account   string `gorm:"column:account" json:"account"`     // 用户账号
	Password  string `gorm:"column:password" json:"password"`   // 用户密码
	Resource1 string `gorm:"column:resource1" json:"resource1"` // 受保护资源1
	Resource2 string `gorm:"column:resource2" json:"resource2"` // 受保护资源2
	Resource3 string `gorm:"column:resource3" json:"resource3"` // 受保护资源3
	BaseInfo
}

type UserInfoDao struct{}

var userInfoDaoOnce sync.Once
var userInfoDaoOnceInstance *UserInfoDao

func UserInfoDaoOnceInstance() *UserInfoDao {
	userInfoDaoOnce.Do(func() {
		userInfoDaoOnceInstance = &UserInfoDao{}
	})
	return userInfoDaoOnceInstance
}

func (*UserInfo) TableName() string {
	return common.UserInfoTableName
}

func (d *UserInfoDao) GetByAccount(ctx *gin.Context, account string) (*UserInfo, error) {
	instance := &UserInfo{}
	err := OAuthDemoDB.Table(common.UserInfoTableName).
		Where("account = ?", account).Where("is_delete = 0").
		First(&instance).Error
	if err != nil {
		return nil, err
	}
	return instance, nil
}

func (d *UserInfoDao) CreateUserInfo(ctx *gin.Context, instance *UserInfo) error {
	return OAuthDemoDB.Table(common.UserInfoTableName).
		Create(instance).Error
}

func (d *UserInfoDao) UpdateUserInfo(ctx *gin.Context, instance *UserInfo) error {
	return OAuthDemoDB.Table(common.UserInfoTableName).
		Where("account = ?", instance.Account).Where("is_delete = 0").
		Updates(instance).Error
}

func (d *UserInfoDao) GetResource(ctx *gin.Context, account, scope string) (*UserInfo, error) {
	instance := &UserInfo{}
	sele := make([]string, 0)
	for _, i := range strings.Split(scope, ",") {
		if _, ok := common.Resource[i]; ok {
			sele = append(sele, i)
		}
	}
	err := OAuthDemoDB.Table(common.UserInfoTableName).
		Where("account = ?", account).Where("is_delete = 0").
		Select(sele).
		First(instance).Error
	if err != nil {
		return nil, err
	}
	return instance, nil
}
