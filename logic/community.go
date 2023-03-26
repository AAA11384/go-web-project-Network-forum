package logic

import (
	"qimiproject/dao/mysql"
	"qimiproject/models"
)

func GetCommunityList() (data []*models.Community, err error) {
	return mysql.GetCommunityList()
}

func GetCommunityDetail(CommunityID int64) (*models.Detail, error) {
	return mysql.GetCommunityListByID(CommunityID)
}
