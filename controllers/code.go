package controllers

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CodeServerBusy
	CodeWrongPassword
	CodeVoteFail

	CodeNeedLogin ResCode = 2000 + iota
	CodeErrorAuth
	CodeInvalidToken
	CodeRequestFrequently
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:           "success",
	CodeInvalidParam:      "请求参数错误",
	CodeUserExist:         "用户名已存在",
	CodeUserNotExist:      "用户名未存在",
	CodeServerBusy:        "服务器繁忙",
	CodeWrongPassword:     "密码错误",
	CodeVoteFail:          "投票失败",
	CodeNeedLogin:         "需要登陆",
	CodeErrorAuth:         "请求头中auth格式有误",
	CodeInvalidToken:      "无效的Token",
	CodeRequestFrequently: "请求过于频繁",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = "服务繁忙"
	}
	return msg
}
