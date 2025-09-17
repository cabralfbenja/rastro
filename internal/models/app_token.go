package models

import "time"

type AppTokenModel struct {
	BaseModel
	Token     string    `json:"-" gorm:"index;type:varchar(255)"`
	TargetId  uint      `json:"targetId" gorm:"index;not null"`
	Type      string    `json:"-" gorm:"index;not null;type:varchar(255)"`
	Used      bool      `json:"-" gorm:"index;not null;type:boolean;default:false"`
	ExpiresAt time.Time `json:"-" gorm:"index;not null"`
}

func (AppTokenModel) TableName() string {
	return "app_tokens"
}
