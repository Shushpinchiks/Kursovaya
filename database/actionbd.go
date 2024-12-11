package database

import (
	model "Kursach/models"
	"io/ioutil"
)

func FetchFileContent(fileName string) (string, error) {
	filePath := "templates/text/" + fileName

	body, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func GetFavoriteMoviesByUserID(userID int) ([]model.Film, error) {
	var movies []model.Film
	var user model.User
	err := DB.Preload("Films").Preload("Films.Genres").Where("id = ?", userID).Find(&user).Error

	for _, movie := range user.Films {
		movies = append(movies, movie)
	}
	return movies, err
}
