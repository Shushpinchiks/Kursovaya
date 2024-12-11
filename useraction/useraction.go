package useraction

import (
	"Kursach/database"
	model "Kursach/models"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"gorm.io/gorm"
)

// Хэширует пароль
func HashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	hashedPassword := hash.Sum(nil)
	return hex.EncodeToString(hashedPassword)
}

// Проверяет, есть ли такой email в БД
func IsUserGetEmail(email string) (bool, error) {
	user := model.User{}
	query := database.DB.First(&user, "email = ?", email)

	if query.Error != nil {
		if errors.Is(query.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, query.Error
	}

	return true, nil
}

func checkPasswordHash(password, hash string) bool {
	if password != hash {
		return false
	}

	return true
}

func ChekUserIsTrue(email, password string) (bool, error) {
	user := model.User{}
	query := database.DB.First(&user, "email = ?", email)
	if query.Error != nil {
		if errors.Is(query.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, query.Error
	}

	if checkPasswordHash(password, user.Password) {
		return true, nil
	}

	return false, nil
}
