package logic

import (
	"qimiproject/dao/mysql"
	"qimiproject/models"
	"qimiproject/pkg/jwt"
	"qimiproject/pkg/snowFlake"
)

func SignUp(p *models.ParamSignUp) (err error) {
	//查询用户是否存在
	err = mysql.CheckUserExist(p.Username)
	if err != nil {
		return err
	}
	//生成UUID
	userid := snowFlake.GenId()
	//构建结构体
	u := &models.User{
		UserId:   userid,
		UserName: p.Username,
		Email:    p.Email,
		Password: p.Password,
	}
	//生成uuid，存入用户
	if err := mysql.InsertUser(u); err != nil {
		return err
	}
	return nil
}

func Login(p *models.ParamLogin) (user *models.User, err error) {
	//链接数据库查询密码
	u := &models.User{
		UserName: p.Username,
		Password: p.Password,
	}
	err = mysql.CheckPassword(u)
	if err != nil {
		return nil, err
	}
	u.Token, err = jwt.GenToken(u.UserId, u.UserName)
	if err != nil {
		return nil, err
	}
	return u, err
}
