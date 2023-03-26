package redis

import (
	"errors"
	"math"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/go-redis/redis/v8"
)

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePreVote     = 432
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeat     = errors.New("禁止重复投票")
)

func CreatePost(postID, communityID int64) error {
	pipeline := client.TxPipeline()
	pipeline.ZAdd(ctx, getRedisKey(KeyPostTime), &redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	pipeline.ZAdd(ctx, getRedisKey(KeyPostScore), &redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	key := getRedisKey(KetCommunitySetPF + strconv.Itoa(int(communityID)))
	if client.Exists(ctx, key).Val() == 0 {
		pipeline.SAdd(ctx, key, postID)
	}
	_, err := pipeline.Exec(ctx)
	return err
}

func VoteForPost(userID, postID string, newValue float64) error {
	//查询帖子是否过期
	postTime := client.ZScore(ctx, getRedisKey(KeyPostTime), postID).Val()
	zap.L().Debug("client.ZScore(ctx, getRedisKey(KeyPostTime), postID).Val()", zap.Float64("val", postTime))
	if time.Now().Unix()-int64(postTime) > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}
	//查询之前是否有投票纪录
	oldValue := client.ZScore(ctx, getRedisKey(KeyPostVotedPrefix+postID), userID).Val()
	if newValue == oldValue {
		return ErrVoteRepeat
	}
	diff := math.Abs(oldValue - newValue)
	var dir float64
	if newValue > oldValue {
		dir = 1
	} else {
		dir = -1
	}
	//更改帖子票数
	pipeline := client.TxPipeline()
	pipeline.ZIncrBy(ctx, getRedisKey(KeyPostScore), dir*diff*scorePreVote, postID)
	//记录用户投票
	if newValue == 0 {
		pipeline.ZRem(ctx, getRedisKey(KeyPostVotedPrefix+postID), userID)
	} else {
		pipeline.ZAdd(ctx, getRedisKey(KeyPostVotedPrefix+postID), &redis.Z{
			Score:  newValue,
			Member: userID,
		})
	}
	_, err := pipeline.Exec(ctx)
	return err
}
