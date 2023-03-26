package models

import "time"

type Community struct {
	CommunityID int    `json:"id" db:"community_id"`
	Name        string `json:"name" db:"community_name"`
}

type Detail struct {
	ID           int       `json:"id" db:"community_id"`
	Name         string    `json:"name" db:"community_name"`
	Introduction string    `json:"introduction,omitempty" db:"introduction"`
	CreateTime   time.Time `json:"create_time" db:"create_time"`
}
