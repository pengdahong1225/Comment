package mysql

import "time"

type RoomComment struct {
	ID        int64     `gorm:"primary_key"`
	CreateAt  time.Time `gorm:"column:create_at" json:"create_at"`
	UpdateAt  time.Time `gorm:"column:update_at" json:"update_at"`
	DeletedAt time.Time `gorm:"column:delete_at" json:"delete_at"`

	RoomId    int64  `gorm:"column:room_id" json:"room_id"`
	UserId    int64  `gorm:"column:user_id" json:"user_id"`
	UserName  string `gorm:"column:user_name" json:"user_name"`
	Content   string `gorm:"column:content" json:"content"`
	PubStamp  int64  `gorm:"column:pub_stamp" json:"pub_stamp"`
	PubRegion string `gorm:"column:pub_region" json:"pub_region"`
}

func (receiver *RoomComment) TableName() string {
	return "room_comment"
}
