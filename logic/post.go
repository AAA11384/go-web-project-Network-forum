package logic

import (
	"errors"
	"qimiproject/dao/mysql"
	"qimiproject/dao/redis"
	"qimiproject/models"
	"qimiproject/pkg/snowFlake"

	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	p.PostID = snowFlake.GenId()
	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}
	err = redis.CreatePost(p.PostID, p.CommunityID)
	return
}

// GetPostDetail 根据帖子ID查询帖子详情
func GetPostDetail(pid int64) (*models.ApiPostDetail, error) {
	post, err := mysql.GetPostByID(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostByID(pid) : ", zap.Error(err))
		return nil, err
	}
	userName, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserByID:", zap.Error(err))
		return nil, err
	}
	community, err := mysql.GetCommunityListByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityByID:", zap.Error(err))
		return nil, err
	}
	//ApiData := &models.ApiPostDetail{
	//	AuthorName: userName,
	//	Post:       post,
	//	Detail:     community,
	ApiData := new(models.ApiPostDetail)
	ApiData.Post = post
	ApiData.Detail = community
	ApiData.AuthorName = userName
	return ApiData, nil
}

// GetPostList 根据页码，页大小查询帖子信息
func GetPostList(page, pageSize int64) ([]*models.ApiPostDetail, error) {
	posts, err := mysql.GetPostList(page, pageSize)
	if err != nil {
		zap.L().Error("mysql.GetPostList:", zap.Error(err))
		return nil, err
	}
	data := make([]*models.ApiPostDetail, 0, len(posts))

	for _, post := range posts {
		userName, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID:", zap.Error(err))
			return nil, err
		}
		community, err := mysql.GetCommunityListByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityByID:", zap.Error(err))
			return nil, err
		}
		ApiData := new(models.ApiPostDetail)
		ApiData.Post = post
		ApiData.Detail = community
		ApiData.AuthorName = userName

		data = append(data, ApiData)
	}
	return data, err
}

// GetPostList2 根据页码，页大小，顺序查询帖子信息
func GetPostList2(post *models.ParamCommunityPost) ([]*models.ApiPostDetail, error) {
	//在redis中查询id列表
	ids, err := redis.GetPostIDInOrder(post)
	if err != nil {
		zap.L().Error("GetPostIDInOrder errors", zap.Error(err))
		return nil, err
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDInOrder success but return 0 lines")
		return nil, errors.New("no data in redis.GetPostIDInOrder")
	}
	//通过id列表，在数据库中查询帖子详细信息
	posts, err := mysql.GetListPostsByIDs(ids)
	if err != nil {
		zap.L().Error("mysql.GetListPostsByIDs errors", zap.Error(err))
		return nil, err
	}
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return nil, err
	}

	data := make([]*models.ApiPostDetail, 0, len(posts))
	for idx, post := range posts {
		userName, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID:", zap.Error(err))
			return nil, err
		}
		community, err := mysql.GetCommunityListByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityByID:", zap.Error(err))
			return nil, err
		}
		ApiData := new(models.ApiPostDetail)
		ApiData.Post = &post
		ApiData.VoteNum = voteData[idx]
		ApiData.Detail = community
		ApiData.AuthorName = userName

		data = append(data, ApiData)
	}
	return data, nil
}

// GetCommunityPostList 根据页码，页大小，Community和顺序查询帖子信息
func GetCommunityPostList(post *models.ParamCommunityPost) ([]*models.ApiPostDetail, error) {
	//在redis中查询id列表
	ids, err := redis.GetCommunityPostIDInOrder(post)
	if err != nil {
		zap.L().Error("GetPostIDInOrder errors", zap.Error(err))
		return nil, err
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDInOrder success but return 0 lines")
		return nil, errors.New("no data in redis.GetPostIDInOrder")
	}
	//通过id列表，在数据库中查询帖子详细信息
	posts, err := mysql.GetListPostsByIDs(ids)
	if err != nil {
		zap.L().Error("mysql.GetListPostsByIDs errors", zap.Error(err))
		return nil, err
	}
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return nil, err
	}

	data := make([]*models.ApiPostDetail, 0, len(posts))
	for idx, post := range posts {
		userName, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID:", zap.Error(err))
			return nil, err
		}
		community, err := mysql.GetCommunityListByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityByID:", zap.Error(err))
			return nil, err
		}
		ApiData := new(models.ApiPostDetail)
		ApiData.Post = &post
		ApiData.VoteNum = voteData[idx]
		ApiData.Detail = community
		ApiData.AuthorName = userName
		data = append(data, ApiData)
	}
	return data, nil
}

// GetPostListByOrderSwitcher 如果ParamCommunityPost中community字段为0
// 调用GetPostList2，否则调用GetCommunityPostList
func GetPostListByOrderSwitcher(post *models.ParamCommunityPost) (data []*models.ApiPostDetail, err error) {
	if post.CommunityID == 0 {
		data, err = GetPostList2(post)
	} else {
		data, err = GetCommunityPostList(post)
	}
	if err != nil {
		zap.L().Error("errors happened in GetPostListByOrderSwitcher", zap.Error(err))
		return nil, err
	}
	return data, nil
}
