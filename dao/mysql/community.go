package mysql

import (
	"database/sql"
	"fmt"
	"qimiproject/models"

	"go.uber.org/zap"
)

func GetCommunityList() (data []*models.Community, err error) {
	sqlstr := "select community_id, community_name from community"
	if err := db.Select(&data, sqlstr); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no community in sql")
			err = nil
			return nil, ErrorInvalidParam
		} else {
			zap.L().Error("mysql.GetCommunity failed")
			return nil, err
		}
	}
	return data, nil
}

func GetCommunityListByID(id int64) (detail *models.Detail, err error) {
	detail = new(models.Detail)
	sqlstr := "select community_id, community_name, introduction, create_time from community where community_id = ?"
	if err := db.Get(detail, sqlstr, id); err != nil {
		fmt.Println("error happened in GetCommunityListByID")
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no community detail in sql")
			return nil, ErrorInvalidParam
		} else {
			zap.L().Warn("mysql.GetCommunityListByID failed")
			return nil, err
		}
	}
	return
}
