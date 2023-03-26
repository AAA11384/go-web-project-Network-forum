package redis

const (
	KeyPrefix          = "bluebell:"
	KeyPostTime        = "post:time"
	KeyPostScore       = "post:score"
	KeyPostVotedPrefix = "post:Voted:"
	KetCommunitySetPF  = "community:"
)

// getRedisKey 返回项目名称+后面内容拼接的key
func getRedisKey(key string) string {
	return KeyPrefix + key
}
