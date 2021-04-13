package service

import (
	"errors"
	"github.com/Zhenghao-Liu/OAuth_demo/common"
	"github.com/Zhenghao-Liu/OAuth_demo/model"
	"github.com/Zhenghao-Liu/OAuth_demo/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"sync"
	"time"
)

type OAuthInfoService struct {
	oauthInfoDao   *model.OAuthInfoDao
	userInfoDao    *model.UserInfoDao
	oauthDemoCache *redis.Client
}

var (
	oauthInfoServiceOnce     sync.Once
	oauthInfoServiceInstance *OAuthInfoService
)

func OAuthInfoServiceInstance() *OAuthInfoService {
	oauthInfoServiceOnce.Do(func() {
		oauthInfoServiceInstance = &OAuthInfoService{
			oauthInfoDao:   model.OAuthInfoDaoInstance(),
			userInfoDao:    model.UserInfoDaoOnceInstance(),
			oauthDemoCache: model.OAuthDemoCache,
		}
	})
	return oauthInfoServiceInstance
}

func (s *OAuthInfoService) CreateOAuthInfo(ctx *gin.Context, instance *model.OAuthInfo) error {
	return s.oauthInfoDao.CreateOAuthInfo(ctx, instance)
}

func (s *OAuthInfoService) GetByAppID(ctx *gin.Context, appID string) (*model.OAuthInfo, error) {
	return s.oauthInfoDao.GetByAppID(ctx, appID)
}

func (s *OAuthInfoService) UpdateOAuthInfo(ctx *gin.Context, instance *model.OAuthInfo) error {
	if err := s.oauthInfoDao.UpdateOAuthInfo(ctx, instance); err != nil {
		return err
	}
	key := common.RedisAppID + instance.AppID
	for _, i := range s.oauthDemoCache.HVals(key).Val() {
		s.delRedisKey(i)
	}
	s.oauthDemoCache.Del(key)
	return nil
}

func (s *OAuthInfoService) delRedisKey(i string) {
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

func (s *OAuthInfoService) Authorize(ctx *gin.Context, appID, account, scope string) (string, error) {
	oldKey := s.oauthDemoCache.HGet(common.RedisAppID+appID, account).Val()
	s.delRedisKey(oldKey)
	oldKey = s.oauthDemoCache.HGet(common.RedisAccount+account, appID).Val()
	s.delRedisKey(oldKey)
	code := utils.GenString()
	key := utils.Encode(code)
	val := appID + "_" + account + "_" + scope
	if err := s.oauthDemoCache.HSet(common.RedisAppID+appID, account, key).Err(); err != nil {
		return "", err
	} else if err := s.oauthDemoCache.HSet(common.RedisAccount+account, appID, key).Err(); err != nil {
		return "", err
	} else if err := s.oauthDemoCache.Set(common.RedisCode+key, val, common.CodeAliveTime*time.Second).Err(); err != nil {
		return "", err
	}
	return code, nil
}

func (s *OAuthInfoService) CheckCode(ctx *gin.Context, code string) bool {
	key := utils.Encode(code)
	return s.oauthDemoCache.Get(common.RedisCode+key).Val() != ""
}

func (s *OAuthInfoService) CheckToken(ctx *gin.Context, token string) bool {
	key := utils.Encode(token)
	return s.oauthDemoCache.Get(common.RedisToken+key).Val() != ""
}

func (s *OAuthInfoService) CheckRefresh(ctx *gin.Context, refresh string) bool {
	key := utils.Encode(refresh)
	return s.oauthDemoCache.Get(common.RedisRefresh+key).Val() != ""
}

func (s *OAuthInfoService) parseVal(ctx *gin.Context, val string) (string, string, string, error) {
	pos := make([]int, 0)
	for i, v := range val {
		if v == '_' {
			pos = append(pos, i)
		}
	}
	if len(pos) < 2 {
		return "", "", "", errors.New(common.StatusDBError)
	}
	return val[:pos[0]], val[pos[0]+1 : pos[1]], val[pos[1]+1:], nil
}

func (s *OAuthInfoService) genTokenRefresh(ctx *gin.Context, val string) (string, string, error) {
	appID, account, _, err := s.parseVal(ctx, val)
	if err != nil {
		return "", "", err
	}
	token := utils.GenString()
	refresh := utils.GenString()
	key := utils.Encode(token + refresh)
	if err = s.oauthDemoCache.HSet(common.RedisAppID+appID, account, key).Err(); err != nil {
		return "", "", err
	} else if err = s.oauthDemoCache.HSet(common.RedisAccount+account, appID, key).Err(); err != nil {
		return "", "", err
	} else if err = s.oauthDemoCache.Set(common.RedisToken+utils.Encode(token), val, common.TokenAliveTime*time.Second).Err(); err != nil {
		return "", "", err
	} else if err = s.oauthDemoCache.Set(common.RedisRefresh+utils.Encode(refresh), val, common.RefreshAliveTime*time.Second).Err(); err != nil {
		return "", "", err
	}
	return token, refresh, nil
}

func (s *OAuthInfoService) Token(ctx *gin.Context, code string) (string, string, error) {
	key := utils.Encode(code)
	val := s.oauthDemoCache.Get(common.RedisCode + key).Val()
	s.delRedisKey(key)
	return s.genTokenRefresh(ctx, val)
}

func (s *OAuthInfoService) Refresh(ctx *gin.Context, refresh string) (string, string, error) {
	key := utils.Encode(refresh)
	val := s.oauthDemoCache.Get(common.RedisRefresh + key).Val()
	appID, account, _, err := s.parseVal(ctx, val)
	if err != nil {
		return "", "", nil
	}
	oldKey := s.oauthDemoCache.HGet(common.RedisAppID+appID, account).Val()
	s.delRedisKey(oldKey)
	oldKey = s.oauthDemoCache.HGet(common.RedisAccount+account, appID).Val()
	s.delRedisKey(oldKey)
	return s.genTokenRefresh(ctx, val)
}

func (s *OAuthInfoService) Resource(ctx *gin.Context, token string) (*model.UserInfo, error) {
	key := utils.Encode(token)
	val := s.oauthDemoCache.Get(common.RedisToken + key).Val()
	_, account, scope, err := s.parseVal(ctx, val)
	if err != nil {
		return nil, err
	}
	return s.userInfoDao.GetResource(ctx, account, scope)
}
