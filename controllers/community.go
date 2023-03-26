package controllers

import (
	"qimiproject/logic"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CommunityHandler(c *gin.Context) {
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunity Failed")
		ResponseError(c, CodeServerBusy)
	}
	ResponseSuccess(c, data)
}

func CommunityDetailHandler(c *gin.Context) {
	CommunityID := c.Param("id")
	id, err := strconv.ParseInt(CommunityID, 10, 64)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityDetail Failed")
		ResponseError(c, CodeInvalidParam)
		return
	}
	ResponseSuccess(c, data)
}
