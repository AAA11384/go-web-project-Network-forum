package controllers

import (
	"errors"
	"fmt"
	"qimiproject/dao/mysql"
	"qimiproject/logic"
	"qimiproject/models"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func SignUpHandler(c *gin.Context) {
	//参数校验
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("SignUpHandler.ShouldBindJSON error", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParam, err.Error())
		return
	}
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("error happened in logic.SignUp", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		} else {
			ResponseErrorWithMsg(c, CodeServerBusy, err.Error())
			return
		}
	}
	ResponseSuccess(c, "success")
	return
}

func LoginHandler(c *gin.Context) {
	//绑定Post表单参数
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("LoginHandler.ShouldBindJSON error", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParam, err.Error())
		return
	}
	//调用logic层login方法
	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("logic.Login error", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		if errors.Is(err, mysql.ErrorWrongPassword) {
			ResponseError(c, CodeWrongPassword)
			return
		}
	}
	ResponseSuccess(c, gin.H{
		"user_id":   fmt.Sprintf("%d", user.UserId),
		"user_name": user.UserName,
		"token":     user.Token,
	})
	return
}
