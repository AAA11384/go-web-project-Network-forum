package mysql

import (
	"qimiproject/models"
	"strings"

	"github.com/jmoiron/sqlx"

	"go.uber.org/zap"
)

func CreatePost(p *models.Post) error {
	sqlStr := "insert into post (post_id, title, content, author_id, community_id) values (?, ?, ?, ?, ?)"
	if _, err := db.Exec(sqlStr, p.PostID, p.Title, p.Content, p.AuthorID, p.CommunityID); err != nil {
		return err
	}
	return nil
}

func GetPostByID(pid int64) (data *models.Post, err error) {
	data = new(models.Post)
	sqlStr := `select 
	post_id, title, content, author_id, community_id, create_time 
	from post 
	where post_id = ?
	`
	if err := db.Get(data, sqlStr, pid); err != nil {
		zap.L().Error("GetPostByID : ", zap.Error(err))
		return nil, err
	}
	return data, nil
}

func GetUserByID(id int64) (userName string, err error) {
	sqlStr := `select username from user where user_id = ?`
	err = db.Get(&userName, sqlStr, id)
	return
}

func GetPostList(page, pageSize int64) (posts []*models.Post, err error) {
	sqlStr := `select
	post_id, author_id, community_id, title, content, create_time
	from post
	order by create_time desc 
	limit ?, ?`
	posts = make([]*models.Post, 0, pageSize)

	err = db.Select(&posts, sqlStr, (page-1)*pageSize, pageSize)
	return
}

func GetListPostsByIDs(ids []string) (postList []models.Post, err error) {
	sqlStr := `select
	post_id, author_id, community_id, title, content, create_time
	from post
	where post_id in (?)
	order by FIND_IN_SET(post_id, ?)
	`
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}

	query = db.Rebind(query)

	err = db.Select(&postList, query, args...)
	return
}
