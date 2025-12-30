package models

type Account struct {
	Username string `json:"username" Gorm:"type:varchar(255); not null"`
	Password string `json:"password" Gorm:"type:varchar(255); not null"`
	Base
}

func (a *Account) TableName() string {
	return "t_account"
}
