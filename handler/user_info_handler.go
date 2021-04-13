package handler

import (
	"errors"
	"github.com/Zhenghao-Liu/OAuth_demo/common"
	"github.com/Zhenghao-Liu/OAuth_demo/model"
	"github.com/Zhenghao-Liu/OAuth_demo/service"
	"github.com/Zhenghao-Liu/OAuth_demo/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type UserInfoHandler struct {
	userInfoService *service.UserInfoService
}

func NewUserInfoHandler() *UserInfoHandler {
	return &UserInfoHandler{
		userInfoService: service.UserInfoServiceInstance(),
	}
}

func (h *UserInfoHandler) Register(ginInstance *gin.Engine) {
	gp := ginInstance.Group("/user")
	gp.POST("/create", JSONWrapper(h.CreateUserInfo))
	gp.POST("/update", JSONWrapper(h.UpdateUserInfo))
	gp.POST("/cancel", JSONWrapper(h.Cancel))
}

func (h *UserInfoHandler) CreateUserInfo(ctx *gin.Context) (interface{}, int, error) {
	account := ctx.Request.Header.Get("account")
	password := ctx.Request.Header.Get("password")
	resource1 := ctx.DefaultPostForm("resource1", "")
	resource2 := ctx.DefaultPostForm("resource2", "")
	resource3 := ctx.DefaultPostForm("resource3", "")
	if account == "" || password == "" || resource1 == "" || resource2 == "" || resource3 == "" {
		return nil, http.StatusBadRequest, errors.New(common.StatusMissParam)
	}
	_, err := h.userInfoService.GetByAccount(ctx, account)
	if err == nil {
		return nil, http.StatusBadRequest, errors.New(common.AccountRepeated)
	} else if err != gorm.ErrRecordNotFound {
		return nil, http.StatusInternalServerError, err
	}
	instance := &model.UserInfo{
		Account:   account,
		Password:  utils.Encode(password),
		Resource1: resource1,
		Resource2: resource2,
		Resource3: resource3,
		BaseInfo:  model.BaseInfo{},
	}
	err = h.userInfoService.CreateUserInfo(ctx, instance)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	mp := make(map[string]interface{})
	mp[common.Final] = common.StatusSuccess
	return mp, http.StatusOK, nil
}

func (h *UserInfoHandler) UpdateUserInfo(ctx *gin.Context) (interface{}, int, error) {
	account := ctx.Request.Header.Get("account")
	password := ctx.Request.Header.Get("password")
	resource1 := ctx.DefaultPostForm("resource1", "")
	resource2 := ctx.DefaultPostForm("resource2", "")
	resource3 := ctx.DefaultPostForm("resource3", "")
	if account == "" || password == "" || resource1 == "" || resource2 == "" || resource3 == "" {
		return nil, http.StatusBadRequest, errors.New(common.StatusMissParam)
	}
	instance, err := h.userInfoService.GetByAccount(ctx, account)
	if err == gorm.ErrRecordNotFound {
		return nil, http.StatusBadRequest, errors.New(common.UserInfoNotFound)
	} else if err != nil {
		return nil, http.StatusInternalServerError, err
	} else if utils.Encode(password) != instance.Password {
		return nil, http.StatusBadRequest, errors.New(common.PasswordError)
	}
	newInstance := &model.UserInfo{
		Account:   instance.Account,
		Password:  instance.Password,
		Resource1: resource1,
		Resource2: resource2,
		Resource3: resource3,
	}
	err = h.userInfoService.UpdateUserInfo(ctx, newInstance)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	mp := make(map[string]interface{})
	mp[common.Final] = common.StatusSuccess
	return mp, http.StatusOK, nil
}

func (h *UserInfoHandler) Cancel(ctx *gin.Context) (interface{}, int, error) {
	account := ctx.Request.Header.Get("account")
	password := ctx.Request.Header.Get("password")
	if account == "" || password == "" {
		return nil, http.StatusBadRequest, errors.New(common.StatusMissParam)
	}
	instance, err := h.userInfoService.GetByAccount(ctx, account)
	if err == gorm.ErrRecordNotFound {
		return nil, http.StatusBadRequest, errors.New(common.UserInfoNotFound)
	} else if err != nil {
		return nil, http.StatusInternalServerError, err
	} else if utils.Encode(password) != instance.Password {
		return nil, http.StatusBadRequest, errors.New(common.PasswordError)
	}
	h.userInfoService.Cancel(ctx, account)
	mp := make(map[string]interface{})
	mp[common.Final] = common.StatusSuccess
	return mp, http.StatusOK, nil
}
