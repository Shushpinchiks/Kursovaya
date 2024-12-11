package connection

import (
	"Kursach/database"
	model "Kursach/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"math"
	"net/http"
	"strconv"
)

func entrance(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userID := session.Get("user_id")
	if userID != nil {
		ctx.Redirect(302, base)
		return
	}
	ctx.HTML(200, "home_page.html", nil)
}

func reg(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userID := session.Get("user_id")
	if userID != nil {
		ctx.Redirect(302, base)
		return
	}
	ctx.HTML(200, "register.html", nil)
}

func movies(ctx *gin.Context) {
	pageSize := 8

	page, err := strconv.Atoi(ctx.Param("num"))
	if err != nil {
		page = 1
	}

	offset := (page - 1) * pageSize

	selectedGenres := ctx.QueryArray("genre")

	var movies []model.Film
	query := database.DB.Preload("Genres").Offset(offset).Limit(pageSize)

	// Если есть выбранные жанры, добавляем фильтрацию
	if len(selectedGenres) > 0 {
		query = query.Joins("JOIN genre_film ON genre_film.film_id = film.id").
			Joins("JOIN genre ON genre.id = genre_film.genre_id").
			Where("genre.name IN ?", selectedGenres).
			Group("film.id")
	}

	if err := query.Find(&movies).Error; err != nil {
		log.Println("Error fetching movies from database:", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	for i, movie := range movies {
		content, err := database.FetchFileContent(movie.Description)
		if err != nil {
			log.Println("Error fetching file content:", err)
		}
		movies[i].Description = content
	}

	var totalCount int64
	database.DB.Model(&model.Film{}).Count(&totalCount)

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	session := sessions.Default(ctx)
	userID := session.Get("user_id")
	user := model.User{}
	database.DB.First(&user, userID)
	queryString := ctx.Request.URL.RawQuery
	data := struct {
		Movies      []model.Film // Ваши данные о фильмах
		CurrentPage int
		HasNextPage bool
		HasPrevPage bool
		NextPage    int
		PrevPage    int
		Pages       []int
		QueryString string
	}{
		Movies:      movies,
		CurrentPage: page,
		HasNextPage: page < totalPages,
		HasPrevPage: page > 1,
		NextPage:    page + 1,
		PrevPage:    page - 1,
		Pages:       make([]int, totalPages),
		QueryString: "?" + queryString,
	}

	for i := 0; i < totalPages; i++ {
		data.Pages[i] = i + 1
	}

	genres := []model.Genre{}
	database.DB.Find(&genres)
	root, _ := strconv.Atoi(user.Root)
	ctx.HTML(http.StatusOK, "movies.html", gin.H{
		"Movies":    data.Movies,
		"data":      data,
		"Root":      root,
		"AllGenres": genres,
	})
}

func favorites(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userID := session.Get("user_id")

	if userID == nil {
		ctx.Redirect(302, base)
		return
	}

	movies, err := database.GetFavoriteMoviesByUserID(userID.(int))
	if err != nil {
		log.Fatal(err)
	}

	var user model.User
	database.DB.First(&user, userID)
	ctx.HTML(200, "favorites.html", gin.H{
		"Movies": movies,
		"user":   user,
	})
}

func film(ctx *gin.Context) {
	filmID := ctx.Param("id")
	renderFilmPage(ctx, filmID)
}

func test(ctx *gin.Context) {

}

func admin(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userID := session.Get("user_id")
	user := model.User{}
	database.DB.First(&user, userID)
	ctx.HTML(200, "admin_panel.html", gin.H{
		"admin": user,
	})
}

func manage_users(ctx *gin.Context) {
	user := []model.User{}
	database.DB.Find(&user)
	ctx.HTML(200, "manage_users.html", gin.H{
		"users": user,
	})
}

func exportData(ctx *gin.Context) {
	ctx.HTML(200, "export_data.html", nil)
}

func addMovie(ctx *gin.Context) {
	ctx.HTML(200, "add_film.html", nil)
}

func deleteMovie(ctx *gin.Context) {
	films := []model.Film{}
	database.DB.Find(&films)
	ctx.HTML(200, "delete_film.html", gin.H{
		"films": films,
	})
}
