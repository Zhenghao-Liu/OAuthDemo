package main

import (
	"fmt"
	"github.com/Zhenghao-Liu/OAuth_demo/common"
	"github.com/Zhenghao-Liu/OAuth_demo/config"
	"github.com/Zhenghao-Liu/OAuth_demo/handler"
	"github.com/Zhenghao-Liu/OAuth_demo/model"
	"github.com/Zhenghao-Liu/OAuth_demo/service"
	"github.com/Zhenghao-Liu/OAuth_demo/utils"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
)

var (
	ginInstance *gin.Engine
)

func main() {
	initUtils()
	initConfig()
	initModel()
	initGin()
	initHandler()
	initService()
	ginInstance.Run(common.HomePage)
}

func initUtils() {
	utils.Init()
}

func initConfig() {
	if err := config.NewConfigInstance(); err != nil {
		panic(err)
	}
}

func initModel() {
	model.InitDatabase()
	model.InitRedis()
}

func initService() {
	service.Init()
}

func CORS() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 允许 Origin 字段中的域发送请求
		context.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		// 设置预验请求有效期为 86400 秒
		context.Writer.Header().Set("Access-Control-Max-Age", "86400")
		// 设置允许请求的方法
		context.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		// 设置允许请求的 Header
		context.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		// 设置拿到除基本字段外的其他字段，如上面的Apitoken, 这里通过引用Access-Control-Expose-Headers，进行配置，效果是一样的。
		context.Writer.Header().Set("Access-Control-Expose-Headers", "*")
		// 配置是否可以带认证信息
		context.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		// OPTIONS请求返回200
		if context.Request.Method == "OPTIONS" {
			context.AbortWithStatus(200)
		} else {
			context.Next()
		}
	}
}

func initUse() {
	logFile, err := os.Create(common.LogFile)
	if err != nil {
		panic(err)
	}
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)
	ginInstance.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - %s \"%s %s %s %d %s \"%s\" %s\"\n",
			param.TimeStamp.Format("2006-01-02 03:04:05"),
			param.ClientIP,
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	ginInstance.Use(CORS())
}

func initGin() {
	ginInstance = gin.Default()
	initUse()
	ginInstance.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
}

func initHandler() {
	handler.Init()
	handler.RegisterHandlers(ginInstance)
}
