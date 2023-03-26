package controllers

import (
	"qimiproject/logic"
	"qimiproject/models"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func CreatePostHandler(c *gin.Context) {
	//参数校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Debug("c.ShouldBindJSON(p)", zap.Any("err", err))
		zap.L().Debug("CreatePostHandler : Invalid Param")
		ResponseError(c, CodeInvalidParam)
		return
	}
	//logic函数具体操作
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePostHandler failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}

// PostDetailHandler 查询单个帖子详情
func PostDetailHandler(c *gin.Context) {
	//获取参数
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("PostDetailHandler: Invalid id param", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//调用logic
	data, err := logic.GetPostDetail(pid)
	if err != nil {
		zap.L().Error("logic.GetPostDetail : error", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

// PostListHandler 获得帖子列表
func PostListHandler(c *gin.Context) {
	page, pageSize := getPageInfo(c)
	data, err := logic.GetPostList(page, pageSize)
	if err != nil {
		zap.L().Error("logic.GetPostList() : ", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

// CommunityPostsOrderHandler 根据排序和社区id获取帖子列表
func CommunityPostsOrderHandler(c *gin.Context) {
	p := &models.ParamCommunityPost{
		Page:        1,
		PageSize:    10,
		Order:       "time",
		CommunityID: 1,
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("CommunityPostsOrderHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	zap.L().Debug("debugger 1", zap.Any("p:", p))
	data, err := logic.GetPostListByOrderSwitcher(p)
	if err != nil {
		zap.L().Error("logic.GetPostList() : ", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}
