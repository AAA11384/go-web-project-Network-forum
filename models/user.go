package models

type User struct {
	UserId     int64  `db:"user_id,string"`
	UserName   string `db:"username"`
	Password   string `db:"password"`
	Email      string `db:"email"`
	Id         int    `db:"id"`
	Gender     int    `db:"gender"`
	CreateTime string `db:"create_time"`
	UpdateTime string `db:"update_time"`
	Token      string
}
