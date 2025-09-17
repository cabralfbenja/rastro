package models

type UserModel struct {
	BaseModel
	FirstName    *string `gorm:"varchar(255)" json:"firstName"`
	LastName     *string `gorm:"varchar(255)" json:"lastName"`
	Username     string  `gorm:"varchar(255);not null;unique" json:"username"`
	Email        string  `gorm:"varchar(200);unique" json:"email"`
	PasswordHash string  `gorm:"varchar(255);not null" json:"-"`
}

func (receiver UserModel) TableName() string {
	return "users"
}
