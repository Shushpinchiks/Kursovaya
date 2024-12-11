package connection

import (
	"Kursach/database"
	model "Kursach/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
)

func renderFilmPage(ctx *gin.Context, filmID string) {
	film := model.Film{}
	database.DB.Preload("Genres").Preload("Reviews").Preload("Reviews.User").Preload("Score").First(&film, filmID)

	content, err := database.FetchFileContent(film.Description)
	if err != nil {
		log.Fatal(err)
	}
	film.Description = content

	user := model.User{}
	session := sessions.Default(ctx)
	userID := session.Get("user_id")
	database.DB.Preload("Scores").First(&user, userID)

	count := database.DB.Model(&user).Where("film_id = ?", filmID).Association("Scores").Count()
	userHasRated := count > 0

	score := model.Score{}
	database.DB.First(&score, "Film_ID = ?", filmID)
	var rating float64
	if score.CntVote == 0 {
		rating = 0
	} else {
		rating = float64(score.Rating) / float64(score.CntVote)
	}
	ctx.HTML(200, "film_one.html", gin.H{
		"film":         film,
		"userHasRated": userHasRated,
		"rating":       rating,
	})
}
