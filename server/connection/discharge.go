package connection

import (
	"Kursach/database"
	model "Kursach/models"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func exportToCSV(ctx *gin.Context) {
	directory := "export/"
	filename := directory + ctx.Param("filename") + ".csv"
	file, err := os.Create(filename)
	if err != nil {
		log.Print(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"ID", "Name", "Director", "Poster", "Description", "Year"})

	movies := []model.Film{}
	database.DB.Preload("Genres").Preload("Score").Preload("Reviews").Find(&movies)
	for _, movie := range movies {
		writer.Write([]string{fmt.Sprint(movie.ID), movie.Name, movie.Director, movie.Poster, movie.Description, movie.Year})
	}
	ctx.Redirect(302, "/admin_panel")
}

func exportToJSON(ctx *gin.Context) {
	directory := "export/"
	filename := directory + ctx.Param("filename") + ".json"
	file, err := os.Create(filename)
	if err != nil {
		fmt.Print(err)
	}
	defer file.Close()

	movies := []model.Film{}
	database.DB.Find(&movies)

	encoder := json.NewEncoder(file)
	encoder.Encode(movies)
	ctx.Redirect(302, "/admin_panel")
}
