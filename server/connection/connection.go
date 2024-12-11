package connection

import (
	"Kursach/database"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

var base string = "/movies/page/1"

func methodGet(router *gin.Engine) {
	router.GET("/", entrance)
	router.GET("/reg", reg)
	router.GET("/test", test)
	router.GET("/movies/page/:num", AuthRequired(), movies)
	router.GET("/logout", AuthRequired(), Logout)
	router.GET("/favorites", AuthRequired(), favorites)
	router.GET("/movie/:id", AuthRequired(), film)
	router.GET("/admin_panel", AuthRequired(), admin)
	router.GET("/manage-users", AuthRequired(), manage_users)
	router.GET("//export-data", AuthRequired(), exportData)
	router.GET("/export/csv/:filename", AuthRequired(), exportToCSV)
	router.GET("/export/json/:filename", AuthRequired(), exportToJSON)
	router.GET("/add-movie", AuthRequired(), addMovie)
	router.GET("/delete-movie", AuthRequired(), deleteMovie)
}

func methodPost(router *gin.Engine) {
	router.POST("/", login)
	router.POST("/reg", register)
	router.POST("/movie/:id", addReview)
	router.POST("/favorites/add/:id", addFavorites)
	router.POST("/rate/:id", rateFilm)
	router.POST("/delete_user", deleteUser)
	router.POST("/add-movie", addNewFilm)
	router.POST("/delete_film", deleteFilm)
}

func StartServer() {
	database.InitEnv()
	database.DB = database.Init()
	router := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("kursach", store))

	router.LoadHTMLGlob("templates/*.html")
	router.Static("/resources", "./templates/resources")
	router.Static("/css", "./templates/css")
	router.Static("/js", "./templates/js")
	methodGet(router)
	methodPost(router)
	go database.TimeBackup()
	router.Run(":8080")
}
