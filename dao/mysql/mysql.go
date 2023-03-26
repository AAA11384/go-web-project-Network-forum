package mysql

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"qimiproject/models"

	"go.uber.org/zap"

	"github.com/spf13/viper"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

const key = "Jialin Liang"

func Init() (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.dbname"),
	)
	// 也可以使用MustConnect连接不成功就panic
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("db conn errors", zap.Error(err))
		return
	}
	db.SetMaxOpenConns(viper.GetInt("max_open_conn"))
	db.SetMaxIdleConns(viper.GetInt("max_idle_conn"))
	return
}

func Close() {
	_ = db.Close()
}

func InsertUser(p *models.User) (err error) {
	//密码加密
	p.Password = encryptPassword(p.Password)
	//执行sql语句入库
	sqlStr := "insert into user(user_id, username, password, email) values(?,?,?,?)"
	if _, err := db.Exec(sqlStr, p.UserId, p.UserName, p.Password, p.Email); err != nil {
		return err
	}
	return
}

func CheckUserExist(username string) (err error) {
	sqlStr := `select count(*) from user where username = ?`
	var count int
	err = db.Get(&count, sqlStr, username)
	if count != 0 {
		return ErrorUserExist
	}
	return
}

func encryptPassword(data string) string {
	h := md5.New()
	h.Write([]byte(key))
	return hex.EncodeToString(h.Sum([]byte(data)))
}

func CheckPassword(user *models.User) (err error) {
	var pas = user.Password
	sqlStr := "select * from user where username = ?"
	if err = db.Get(user, sqlStr, user.UserName); err != nil {
		return ErrorUserNotExist
	}
	if encryptPassword(pas) != user.Password {
		return ErrorWrongPassword
	}
	return
}
