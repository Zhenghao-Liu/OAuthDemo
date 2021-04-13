package common

const (
	StatusSuccess   = "操作成功"
	StatusDBError   = "操作失败:数据库异常"
	StatusParamErr  = "操作失败:参数错误"
	StatusMissParam = "操作失败:缺少必要参数"
)

const (
	UserInfoNotFound = "操作失败:缺少用户信息"
	AccountRepeated  = "操作失败:账号已经存在"
	PasswordError    = "操作失败:密码错误"
)

const (
	OAuthInfoNotFound = "操作失败:缺少OAuth信息"
	AppSecretError    = "操作失败:客户端验证失败"
	CodeError         = "操作失败:授权码验证失败"
	RefreshError      = "操作失败:刷新令牌验证失败"
	TokenError        = "操作失败:访问令牌验失败"
)

const (
	InvalidRequest          = "invalid_request"
	UnauthorizedClient      = "unauthorized_client"
	UnsupportedResponseType = "unsupported_response_type"
	InvalidScope            = "invalid_scope"
	ServerError             = "server_error"
	TemporarilyUnavailable  = "temporarily_unavailable"
	LogInError              = "login_error"
)
