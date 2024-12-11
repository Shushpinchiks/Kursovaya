package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       int    `gorm:"primary_key;AUTO_INCREMENT"`
	Name     string `gorm:"size:255"`
	Email    string `gorm:"size:255"`
	Password string `gorm:"size:255"`
	Root     string `gorm:"size:255"`

	Reviews []Review `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Films   []Film   `gorm:"many2many:favourites;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Scores  []Score  `gorm:"many2many:scores;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Review struct {
	gorm.Model
	ID     int    `gorm:"primary_key;AUTO_INCREMENT"`
	Rating int    `gorm:"size:10"`
	Text   string `gorm:"type:text"`
	UserID int    `gorm:"index"`
	FilmID int    `gorm:"index"`

	User User `gorm:"constraint:OnUpdate:CASCADE"`
	Film Film `gorm:"constraint:OnUpdate:CASCADE"`
}

type Film struct {
	gorm.Model
	ID          int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Name        string `gorm:"size:255" json:"name"`
	Director    string `gorm:"size:255" json:"director"`
	Poster      string `gorm:"size:255" json:"poster"`
	Description string `gorm:"type:text" json:"description"`
	Year        string `gorm:"size:255" json:"year"`

	Reviews []Review `gorm:"foreignKey:FilmID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Score   Score    `gorm:"foreignKey:FilmID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Genres  []Genre  `gorm:"many2many:genre_film;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Users   []User   `gorm:"many2many:favourites;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Score struct {
	gorm.Model
	ID      int `gorm:"primary_key;AUTO_INCREMENT"`
	Rating  int `gorm:"size:255"`
	CntVote int `gorm:"size:255"`
	FilmID  int `gorm:"index"`

	Users []User `gorm:"many2many:scores;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Genre struct {
	gorm.Model
	ID    int    `gorm:"primary_key;AUTO_INCREMENT"`
	Name  string `gorm:"size:255"`
	Films []Film `gorm:"many2many:genre_film;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (User) TableName() string {
	return "public.user"
}

func (Review) TableName() string {
	return "public.review"
}

func (Film) TableName() string {
	return "public.film"
}

func (Score) TableName() string {
	return "public.score"
}

func (Genre) TableName() string {
	return "public.genre"
}
