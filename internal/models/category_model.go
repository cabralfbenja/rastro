package models

type CategoryModel struct {
	BaseModel
	Name     string `json:"name" gorm:"type:varchar(100);not null;unique"`
	Slug     string `json:"slug" gorm:"type:varchar(100);not null;unique"`
	IsCustom bool   `json:"isCustom" gorm:"default:false;type:bool"`
}

func (CategoryModel) TableName() string {
	return "categories"
}
