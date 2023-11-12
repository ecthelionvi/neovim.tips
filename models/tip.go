// model/tip.go
package model

type Tip struct {
    ID      uint   `gorm:"primaryKey"`
    Content string `gorm:"type:text;not null"`
}
