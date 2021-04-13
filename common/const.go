package common

// URL
const (
	HomePage  = "127.0.0.1:22222"
	OAuthPage = "http://" + HomePage + "/oauth/authorize"
)

// FILE
const (
	ConfFile = "conf/OAuth_demo.json"
	LogFile  = "output/log/gin.log"
)

// CONST
const (
	StringUpper       = 64
	CodeAliveTime     = 300
	TokenAliveTime    = 600
	RefreshAliveTime  = 3600
	StringAll         = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	ResponseType      = "code"
	Resource1         = "resource1"
	Resource2         = "resource2"
	Resource3         = "resource3"
	AuthorizationCode = "authorization_code"
	RefreshToken      = "refresh_token"
	TokenType         = "bearer"
	Final             = "result"
)

// model
const (
	UserInfoTableName  = "user_info"
	OAuthInfoTableName = "oauth_info"
	RedisAppID         = "app_id:"
	RedisAccount       = "account:"
	RedisCode          = "code:"
	RedisToken         = "token:"
	RedisRefresh       = "refresh:"
)

var (
	Resource = map[string]bool{Resource1: true, Resource2: true, Resource3: true}
)
