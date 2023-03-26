package models

const (
	OrderTime  = "time"
	OrderScore = "score"
)

// ParamSignUp 注册请求参数结构体
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
	Email      string `json:"email" binding:"required,email"`
}

// ParamLogin 登录请求参数结构体
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ParamVoteData 投票参数结构体
type ParamVoteData struct {
	PostID    string `json:"post_id" binding:"required"`
	Direction int8   `json:"direction,string" binding:"required,oneof=1 0 -1"`
}

// ParamPost 请求帖子参数结构体
type ParamPost struct {
	Page     int64  `json:"page" binding:"required" form:"page"`
	PageSize int64  `json:"pageSize" form:"pageSize"`
	Order    string `json:"order" binding:"required,oneof=time score" form:"order"`
}

// ParamCommunityPost 社区请求帖子参数结构体
type ParamCommunityPost struct {
	Page        int64  `json:"page" binding:"required" form:"page"`
	PageSize    int64  `json:"pageSize" form:"pageSize"`
	CommunityID int64  `json:"community_id" from:"community_id"`
	Order       string `json:"order" binding:"required,oneof=time score" form:"order"`
}
