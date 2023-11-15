package model

type Tip struct {
	ID      uint   `gorm:"primaryKey"`
	Content string `gorm:"type:text;not null"`
}

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	IsSuper  bool
}
