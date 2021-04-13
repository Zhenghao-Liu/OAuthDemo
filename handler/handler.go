package handler

import (
	"github.com/Zhenghao-Liu/OAuth_demo/common"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

type Handler interface {
	Register(ginInstance *gin.Engine)
}

var handlers []Handler

func Init() {
	handlers = make([]Handler, 0)
	handlers = append(handlers, NewUserInfoHandler())
	handlers = append(handlers, NewOAuthInfoHandler())
}

func RegisterHandlers(ginInstance *gin.Engine) {
	ginInstance.LoadHTMLGlob("html/*")
	for _, h := range handlers {
		h.Register(ginInstance)
	}
}

type defaultHandlerFunc func(*gin.Context) (interface{}, int, error)

func JSONWrapper(handlerFunc defaultHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ans, code, err := handlerFunc(c)
		if err != nil {
			c.IndentedJSON(code, map[string]interface{}{
				common.Final: err.Error(),
			})
		} else {
			c.IndentedJSON(code, ans)
		}
	}
}

type redirectHandlerFunc func(*gin.Context) (map[string]string, string, string)

func RediectWrapper(handlerFunc redirectHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ans, tar, err := handlerFunc(c)
		if err != "" {
			c.Redirect(http.StatusFound, tar+"?error="+err)
		} else {
			tarUrl := []byte(tar)
			tarUrl = append(tarUrl, '?')
			for k, v := range ans {
				vv, _ := url.QueryUnescape(v)
				tarUrl = append(tarUrl, k+"="+vv+"&"...)
			}
			tarUrl = tarUrl[:len(tarUrl)-1]
			c.Redirect(http.StatusFound, string(tarUrl))
		}
	}
}
