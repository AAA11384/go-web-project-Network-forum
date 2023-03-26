package mysql

import "errors"

var (
	ErrorUserNotExist  = errors.New("用户未存在")
	ErrorUserExist     = errors.New("用户已存在")
	ErrorWrongPassword = errors.New("密码错误")
	ErrorInvalidParam  = errors.New("无效的参数")
)
