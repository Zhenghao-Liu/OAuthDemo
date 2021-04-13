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
	"net/url"
	"strconv"
)

type OAuthInfoHandler struct {
	oauthInfoService *service.OAuthInfoService
	userInfoService  *service.UserInfoService
}

func NewOAuthInfoHandler() *OAuthInfoHandler {
	return &OAuthInfoHandler{
		oauthInfoService: service.OAuthInfoServiceInstance(),
		userInfoService:  service.UserInfoServiceInstance(),
	}
}

func (h *OAuthInfoHandler) Register(ginInstance *gin.Engine) {
	gp := ginInstance.Group("/oauth")
	gp.POST("/create", JSONWrapper(h.CreateOAuthInfo))
	gp.POST("/update", JSONWrapper(h.UpdateOAuthInfo))
	gp.POST("/authorize", RediectWrapper(h.Authorize))
	gp.POST("/token", JSONWrapper(h.Token))
	gp.POST("/refresh", JSONWrapper(h.Refresh))
	gp.POST("/resource", JSONWrapper(h.Resource))
	gp.POST("/welcome", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "welcome.html", gin.H{
			"url": common.OAuthPage,
		})
	})
}

func (h *OAuthInfoHandler) CreateOAuthInfo(ctx *gin.Context) (interface{}, int, error) {
	appName := ctx.DefaultPostForm("app_name", "")
	homepage := ctx.DefaultPostForm("homepage", "")
	description := ctx.DefaultPostForm("description", "")
	callback := ctx.DefaultPostForm("callback", "")
	if appName == "" || homepage == "" || callback == "" {
		return nil, http.StatusBadRequest, errors.New(common.StatusMissParam)
	}
	appID := utils.GenString()
	appSecret := utils.GenString()
	instance := &model.OAuthInfo{
		AppName:     appName,
		Homepage:    homepage,
		Description: description,
		Callback:    callback,
		AppID:       appID,
		AppSecret:   utils.Encode(appSecret),
	}
	err := h.oauthInfoService.CreateOAuthInfo(ctx, instance)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	mp := make(map[string]interface{})
	mp["app_id"] = appID
	mp["app_secret"] = appSecret
	mp[common.Final] = common.StatusSuccess
	return mp, http.StatusOK, nil
}

func (h *OAuthInfoHandler) UpdateOAuthInfo(ctx *gin.Context) (interface{}, int, error) {
	appName := ctx.DefaultPostForm("app_name", "")
	homepage := ctx.DefaultPostForm("homepage", "")
	description := ctx.DefaultPostForm("description", "")
	callback := ctx.DefaultPostForm("callback", "")
	appID := ctx.Request.Header.Get("app_id")
	appSecret := ctx.Request.Header.Get("app_secret")
	if appName == "" || homepage == "" || callback == "" || appID == "" || appSecret == "" {
		return nil, http.StatusBadRequest, errors.New(common.StatusMissParam)
	}
	instance, err := h.oauthInfoService.GetByAppID(ctx, appID)
	if err == gorm.ErrRecordNotFound {
		return nil, http.StatusBadRequest, errors.New(common.OAuthInfoNotFound)
	} else if err != nil {
		return nil, http.StatusInternalServerError, err
	} else if utils.Encode(appSecret) != instance.AppSecret {
		return nil, http.StatusBadRequest, errors.New(common.AppSecretError)
	}
	newInstance := &model.OAuthInfo{
		AppID:       appID,
		AppName:     appName,
		Homepage:    homepage,
		Description: description,
		Callback:    callback,
	}
	err = h.oauthInfoService.UpdateOAuthInfo(ctx, newInstance)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	mp := make(map[string]interface{})
	mp[common.Final] = common.StatusSuccess
	return mp, http.StatusOK, nil
}

func (h *OAuthInfoHandler) Authorize(ctx *gin.Context) (map[string]string, string, string) {
	appID, _ := url.QueryUnescape(ctx.Request.Header.Get("app_id"))
	account := ctx.Request.Header.Get("account")
	password := ctx.Request.Header.Get("password")
	responseType := ctx.Request.Header.Get("response_type")
	callback, _ := url.QueryUnescape(ctx.Request.Header.Get("callback"))
	scope := ctx.Request.Header.Get("scope")
	state := ctx.Request.Header.Get("state")
	if appID == "" || scope == "" || account == "" || password == "" || callback == "" {
		return nil, callback, common.InvalidRequest
	}
	oauthInfo, err := h.oauthInfoService.GetByAppID(ctx, appID)
	if err == gorm.ErrRecordNotFound {
		return nil, callback, common.UnauthorizedClient
	} else if err != nil {
		return nil, callback, common.ServerError
	}
	userInfo, err := h.userInfoService.GetByAccount(ctx, account)
	if err == gorm.ErrRecordNotFound {
		return nil, callback, common.LogInError
	} else if err != nil {
		return nil, callback, common.ServerError
	} else if responseType != common.ResponseType || callback != oauthInfo.Callback {
		return nil, callback, common.InvalidRequest
	} else if utils.Encode(password) != userInfo.Password {
		return nil, callback, common.LogInError
	}
	code, err := h.oauthInfoService.Authorize(ctx, appID, account, scope)
	if err != nil {
		return nil, callback, err.Error()
	}
	mp := make(map[string]string)
	mp["code"] = url.QueryEscape(code)
	mp["expires"] = strconv.Itoa(common.CodeAliveTime)
	mp["state"] = state
	return mp, callback, ""
}

func (h *OAuthInfoHandler) Token(ctx *gin.Context) (interface{}, int, error) {
	appID := ctx.Request.Header.Get("app_id")
	appSecret := ctx.Request.Header.Get("app_secret")
	code := ctx.Request.Header.Get("code")
	grantType := ctx.DefaultPostForm("grant_type", "")
	callback := ctx.DefaultPostForm("callback", "")
	if appID == "" || appSecret == "" || code == "" || grantType == "" || callback == "" {
		return nil, http.StatusBadRequest, errors.New(common.StatusMissParam)
	}
	oauthInfo, err := h.oauthInfoService.GetByAppID(ctx, appID)
	if err == gorm.ErrRecordNotFound {
		return nil, http.StatusBadRequest, errors.New(common.OAuthInfoNotFound)
	} else if err != nil {
		return nil, http.StatusInternalServerError, err
	} else if callback != oauthInfo.Callback || grantType != common.AuthorizationCode {
		return nil, http.StatusBadRequest, errors.New(common.StatusParamErr)
	} else if utils.Encode(appSecret) != oauthInfo.AppSecret {
		return nil, http.StatusBadRequest, errors.New(common.AppSecretError)
	} else if !h.oauthInfoService.CheckCode(ctx, code) {
		return nil, http.StatusBadRequest, errors.New(common.CodeError)
	}
	token, refresh, err := h.oauthInfoService.Token(ctx, code)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	mp := make(map[string]interface{})
	mp["access_token"] = token
	mp["token_type"] = common.TokenType
	mp["refresh_token"] = refresh
	mp["token_expires"] = common.TokenAliveTime
	mp["refresh_expires"] = common.RefreshAliveTime
	mp[common.Final] = common.StatusSuccess
	return mp, http.StatusOK, nil
}

func (h *OAuthInfoHandler) Refresh(ctx *gin.Context) (interface{}, int, error) {
	appID := ctx.Request.Header.Get("app_id")
	appSecret := ctx.Request.Header.Get("app_secret")
	refresh := ctx.Request.Header.Get("refresh_token")
	grantType := ctx.DefaultPostForm("grant_type", "")
	callback := ctx.DefaultPostForm("callback", "")
	if appID == "" || appSecret == "" || refresh == "" || grantType == "" || callback == "" {
		return nil, http.StatusBadRequest, errors.New(common.StatusMissParam)
	}
	oauthInfo, err := h.oauthInfoService.GetByAppID(ctx, appID)
	if err == gorm.ErrRecordNotFound {
		return nil, http.StatusBadRequest, errors.New(common.OAuthInfoNotFound)
	} else if err != nil {
		return nil, http.StatusInternalServerError, err
	} else if callback != oauthInfo.Callback || grantType != common.RefreshToken {
		return nil, http.StatusBadRequest, errors.New(common.StatusParamErr)
	} else if utils.Encode(appSecret) != oauthInfo.AppSecret {
		return nil, http.StatusBadRequest, errors.New(common.AppSecretError)
	} else if !h.oauthInfoService.CheckRefresh(ctx, refresh) {
		return nil, http.StatusBadRequest, errors.New(common.RefreshError)
	}
	token, refresh, err := h.oauthInfoService.Refresh(ctx, refresh)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	mp := make(map[string]interface{})
	mp["access_token"] = token
	mp["token_type"] = common.TokenType
	mp["refresh_token"] = refresh
	mp["token_expires"] = common.TokenAliveTime
	mp["refresh_expires"] = common.RefreshAliveTime
	mp[common.Final] = common.StatusSuccess
	return mp, http.StatusOK, nil
}

func (h *OAuthInfoHandler) Resource(ctx *gin.Context) (interface{}, int, error) {
	appID := ctx.Request.Header.Get("app_id")
	appSecret := ctx.Request.Header.Get("app_secret")
	token := ctx.Request.Header.Get("access_token")
	tokenType := ctx.DefaultPostForm("token_type", "")
	if appID == "" || appSecret == "" || token == "" || tokenType == "" {
		return nil, http.StatusBadRequest, errors.New(common.StatusMissParam)
	} else if tokenType != common.TokenType {
		return nil, http.StatusBadRequest, errors.New(common.StatusParamErr)
	}
	oauthInfo, err := h.oauthInfoService.GetByAppID(ctx, appID)
	if err == gorm.ErrRecordNotFound {
		return nil, http.StatusBadRequest, errors.New(common.OAuthInfoNotFound)
	} else if err != nil {
		return nil, http.StatusInternalServerError, err
	} else if utils.Encode(appSecret) != oauthInfo.AppSecret {
		return nil, http.StatusBadRequest, errors.New(common.AppSecretError)
	} else if !h.oauthInfoService.CheckToken(ctx, token) {
		return nil, http.StatusBadRequest, errors.New(common.TokenError)
	}
	resource, err := h.oauthInfoService.Resource(ctx, token)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	mp := make(map[string]interface{})
	mp["resource1"] = resource.Resource1
	mp["resource2"] = resource.Resource2
	mp["resource3"] = resource.Resource3
	mp[common.Final] = common.StatusSuccess
	return mp, http.StatusOK, nil
}
