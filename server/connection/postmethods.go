package connection

import (
	"Kursach/database"
	model "Kursach/models"
	"Kursach/useraction"
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func register(ctx *gin.Context) {
	email := ctx.PostForm("email")
	username := ctx.PostForm("username")
	password := useraction.HashPassword(ctx.PostForm("password"))
	conPassword := useraction.HashPassword(ctx.PostForm("confirm_password"))
	fmt.Println(email, username, password, conPassword)
	res, err := useraction.IsUserGetEmail(email)
	if err != nil {
		log.Fatal(err)
	}
	if res {
		ctx.HTML(200, "register.html", gin.H{
			"Error": "Этот email уже зарегистрирован.",
		})
	} else {
		if password != conPassword {
			ctx.HTML(200, "register.html", gin.H{
				"Error": "Пароли не совпадают.",
			})
			return
		} else {
			user := model.User{Name: username, Email: email, Password: password, Root: "0"}
			database.DB.Create(&user)

			ctx.Redirect(302, "/")
		}
	}
}

func login(ctx *gin.Context) {
	email := ctx.PostForm("email")
	password := useraction.HashPassword(ctx.PostForm("password"))
	res, err := useraction.ChekUserIsTrue(email, password)
	if err != nil {
		log.Fatal(err)
	}
	if res == false {
		ctx.HTML(200, "home_page.html", gin.H{
			"Error": "неправильный email или пароль",
		})
	} else {
		user := model.User{}
		database.DB.Where("email = ?", email).First(&user)
		session := sessions.Default(ctx)
		session.Set("user_id", user.ID)
		session.Set("user_root", user.Root)
		err := session.Save()
		if err != nil {
			log.Fatal(err)
		}
		ctx.Redirect(302, base)
		return
	}
}

func addReview(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userID := session.Get("user_id")
	reviewText := ctx.PostForm("review")
	filmID := ctx.Param("id")

	filmIDInt, _ := strconv.Atoi(filmID)

	checkRev := model.Review{}
	query := database.DB.First(&checkRev, "User_ID = ? AND Film_ID = ?", userID.(int), filmIDInt)
	if query.Error != nil {
		if errors.Is(query.Error, gorm.ErrRecordNotFound) {
			review := model.Review{
				Rating: 0,
				Text:   reviewText,
				UserID: userID.(int),
				FilmID: filmIDInt,
			}
			database.DB.Create(&review)

			user := model.User{}
			database.DB.First(&user, userID)
			database.DB.Model(&user).Association("Reviews").Append(&review)

			movie := model.Film{}
			database.DB.First(&movie, filmID)
			database.DB.Model(&movie).Association("Reviews").Append(&review)

			renderFilmPage(ctx, filmID)
			return
		}

		renderFilmPage(ctx, filmID)
		return
	}

	checkRev.Text = reviewText
	database.DB.Save(&checkRev)
	renderFilmPage(ctx, filmID)
}

func addFavorites(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userID := session.Get("user_id")
	filmID := ctx.Param("id")

	film := model.Film{}
	database.DB.First(&film, "ID = ?", filmID)
	user := model.User{}
	database.DB.First(&user, userID)
	database.DB.Model(&user).Association("Films").Append(&film)
}

func rateFilm(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userID := session.Get("user_id")
	filmID := ctx.Param("id")
	filmIDInt, _ := strconv.Atoi(filmID)

	var input struct {
		Rating string `json:"rating"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	ratingInt, err := strconv.Atoi(input.Rating)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
	}

	var score model.Score
	var user model.User
	database.DB.First(&user, userID)

	query := database.DB.Preload("Users").First(&score, "Film_ID = ?", filmID)
	if query.Error != nil {
		if errors.Is(query.Error, gorm.ErrRecordNotFound) {
			firstScore := model.Score{
				Rating:  ratingInt,
				CntVote: 1,
				FilmID:  filmIDInt,
			}
			if err := database.DB.Create(&firstScore).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save rating"})
				return
			}
			database.DB.Model(&user).Association("Scores").Append(&firstScore)
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
	} else {
		userRated := false
		for _, u := range score.Users {
			if u.ID == userID {
				userRated = true
				break
			}
		}

		if userRated {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "User has already rated this film"})
			return
		} else {
			score.Rating += ratingInt
			score.CntVote++
			if err := database.DB.Save(&score).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save rating"})
				return
			}
			database.DB.Model(&user).Association("Scores").Append(&score)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Rating saved successfully"})
}

func deleteUser(ctx *gin.Context) {
	var input struct {
		UserID int `json:"id"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user := model.User{}
	database.DB.First(&user, input.UserID)
	if err := database.DB.Delete(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "User deleted successfully"})

}

func addNewFilm(ctx *gin.Context) {
	name := ctx.PostForm("name")
	director := ctx.PostForm("director")
	poster := "resources/" + ctx.PostForm("poster")
	description := ctx.PostForm("description")
	year := ctx.PostForm("year")
	input := ctx.PostForm("genres")

	genres := strings.Split(input, ",")
	for i := range genres {
		genres[i] = strings.TrimSpace(genres[i])
	}
	for i := range genres {
		fmt.Println(genres[i])
		fmt.Println("-----")
	}

	filePath := fmt.Sprintf("templates/text/%s.txt", name)
	err := os.WriteFile(filePath, []byte(description), 0644)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write description to file"})
		return
	}
	pathDes := fmt.Sprintf("%s.txt", name)

	movie := model.Film{}
	query := database.DB.First(&movie, "Name = ?", name)
	if query.Error != nil {
		if errors.Is(query.Error, gorm.ErrRecordNotFound) {
			movieNew := model.Film{
				Name:        name,
				Director:    director,
				Poster:      poster,
				Description: pathDes,
				Year:        year,
			}
			database.DB.Create(&movieNew)
			for _, genre := range genres {
				genr := model.Genre{}
				database.DB.First(&genr, "Name = ?", genre)
				database.DB.Model(&movie).Association("Genres").Append(&genr)
			}
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
	} else {
		for _, genre := range genres {
			genr := model.Genre{}
			database.DB.First(&genr, "Name = ?", genre)
			database.DB.Model(&movie).Association("Genres").Append(&genr)
		}
	}

	ctx.HTML(200, "add_film.html", nil)
}

func deleteFilm(ctx *gin.Context) {
	var input struct {
		FilmID int `json:"id"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	film := model.Film{}
	database.DB.First(&film, input.FilmID)

	if err := database.DB.Delete(&film).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete film"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Фильм был удалён!"})
}
