package utils

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"neovim-tips/database"
	"neovim-tips/models"
)

func PopulateTips() {
	tips := []model.Tip{}

	for _, tip := range tips {
		if !TipExists(tip.Content) {
			database.DB.Create(&tip)
		}
	}
}

func TipExists(content string) bool {
	var count int64
	database.DB.Model(&model.Tip{}).Where("content = ?", content).Count(&count)
	return count > 0
}

func CreateSuperUser(username, password string) {
	var count int64
	database.DB.Model(&model.User{}).Where("username = ?", username).Count(&count)
	if count == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("Failed to hash password:", err)
		}
		superUser := model.User{
			Username: username,
			Password: string(hashedPassword),
			IsSuper:  true,
		}
		database.DB.Create(&superUser)
	}
}
