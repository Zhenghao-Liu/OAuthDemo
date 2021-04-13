package service

import (
	"github.com/Zhenghao-Liu/OAuth_demo/common"
	"github.com/Zhenghao-Liu/OAuth_demo/model"
	"github.com/Zhenghao-Liu/OAuth_demo/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"sync"
)

type UserInfoService struct {
	oauthDemoCache *redis.Client
	userInfoDao    *model.UserInfoDao
}

var (
	userInfoServiceOnce     sync.Once
	userInfoServiceInstance *UserInfoService
)

func UserInfoServiceInstance() *UserInfoService {
	userInfoServiceOnce.Do(func() {
		userInfoServiceInstance = &UserInfoService{
			userInfoDao:    model.UserInfoDaoOnceInstance(),
			oauthDemoCache: model.OAuthDemoCache,
		}
	})
	return userInfoServiceInstance
}

func (s *UserInfoService) GetByAccount(ctx *gin.Context, account string) (*model.UserInfo, error) {
	return s.userInfoDao.GetByAccount(ctx, account)
}

func (s *UserInfoService) CreateUserInfo(ctx *gin.Context, instance *model.UserInfo) error {
	return s.userInfoDao.CreateUserInfo(ctx, instance)
}

func (s *UserInfoService) Cancel(ctx *gin.Context, account string) {
	key := common.RedisAccount + account
	for _, i := range s.oauthDemoCache.HVals(key).Val() {
		s.delRedisKey(i)
	}
	s.oauthDemoCache.Del(key)
}

func (s *UserInfoService) UpdateUserInfo(ctx *gin.Context, instance *model.UserInfo) error {
	if err := s.userInfoDao.UpdateUserInfo(ctx, instance); err != nil {
		return err
	}
	s.Cancel(ctx, instance.Account)
	return nil
}

func (s *UserInfoService) delRedisKey(i string) {
	ii := utils.Decode(i)
	if len(ii) == common.StringUpper {
		s.oauthDemoCache.Del(common.RedisCode + i)
	} else if len(ii) == common.StringUpper*2 {
		l := utils.Encode(ii[:common.StringUpper])
		r := utils.Encode(ii[common.StringUpper:])
		s.oauthDemoCache.Del(common.RedisToken + l)
		s.oauthDemoCache.Del(common.RedisRefresh + r)
	}
}
