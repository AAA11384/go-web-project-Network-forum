package redis

import (
	"qimiproject/models"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

// GetPostIDInOrder 根据期望的类型(time or score)获取帖子的ID
func GetPostIDInOrder(p *models.ParamCommunityPost) ([]string, error) {
	//从redis获取id
	start := (p.Page - 1) * p.PageSize
	end := start + p.PageSize - 1
	key := ""
	if p.Order == models.OrderTime {
		key = getRedisKey(KeyPostTime)
	} else {
		key = getRedisKey(KeyPostScore)
	}
	return client.ZRevRange(ctx, key, start, end).Result()
}

// GetPostVoteData 按照id查询帖子的赞成票
func GetPostVoteData(ids []string) ([]int64, error) {
	res := make([]int64, 0, len(ids))
	pipeline := client.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedPrefix + id)
		pipeline.ZCount(ctx, key, "1", "1")
	}
	cmders, err := pipeline.Exec(ctx)
	if err != nil {
		return nil, err
	}
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		res = append(res, v)
	}
	return res, nil
}

// GetCommunityPostIDInOrder 根据社区获取ids
// 通过community 和 post/order取交集获得结果集
// 若存在则直接返回ids，不存在则新建并缓存60s
func GetCommunityPostIDInOrder(post *models.ParamCommunityPost) ([]string, error) {
	//根据order以及Community取数据，将两个表取交集 zinterstore
	//post:order 和 community:id 取交集,得到该社区下帖子按照特定order排序的信息
	//新生成的表叫做 community/id/order
	//同时利用缓存key减少zinterstore的执行次数
	newZSetKey := getRedisKey(KetCommunitySetPF + strconv.Itoa(int(post.CommunityID)) + ":" + post.Order)
	//不存在或已过期，重新计算
	if client.Exists(ctx, newZSetKey).Val() == 0 {
		communityKey := getRedisKey(KetCommunitySetPF + strconv.Itoa(int(post.CommunityID)))
		postOrderKey := ""
		if post.Order == "time" {
			postOrderKey = getRedisKey(KeyPostTime)
		} else {
			postOrderKey = getRedisKey(KeyPostScore)
		}
		pipeline := client.Pipeline()
		pipeline.ZInterStore(ctx, newZSetKey, &redis.ZStore{
			Keys:      []string{communityKey, postOrderKey},
			Aggregate: "MAX",
		})
		pipeline.Expire(ctx, newZSetKey, time.Second*60)
		_, err := pipeline.Exec(ctx)
		if err != nil {
			return nil, err
		}
	}
	//从redis获取id
	start := (post.Page - 1) * post.PageSize
	end := start + post.PageSize - 1
	return client.ZRevRange(ctx, newZSetKey, start, end).Result()
}
