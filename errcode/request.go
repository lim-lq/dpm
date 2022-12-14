package errcode

const (
	// 解析参数失败
	PARSE_PARAMETER_ERROR = 10001
	// 获取项目列表失败
	GET_PROJECT_LIST_ERROR = 10002
	// 项目未找到
	PROJECT_NOT_FOUND = 10003
	// 获取应用列表失败
	GET_APPLICTION_LIST_ERROR = 10004
	// 用户未找到
	USER_NOT_FOUND = 10005
	// 登录失败
	LOGIN_INFO_ERROR = 10006
	// 获取token失败
	GEN_TOKEN_FAILURE = 10007
	// 无效Token
	INVALID_TOKEN = 10008
	// 创建账号失败
	CREATE_ACCOUNT_ERROR = 10009
	// 缺少参数
	PARAMETER_MISSING = 10010
	// 删除账号失败
	DELETE_ACCOUNT_ERROR = 10011
	// 获取账号失败
	GET_ACCOUNT_ERROR = 10012
	// 账号未找到
	ACCOUNT_NOT_FOUND = 10013
	// 更新账号失败
	UPDATE_ACCOUNT_ERROR = 10014
	// rsa解密失败
	DECRYPT_FAILURE = 10015
	// 创建项目失败
	CREATE_PROJECT_ERROR = 10016
	// 更新权限失败
	UPDATE_PRIVILEGE_ERROR = 10017
)
