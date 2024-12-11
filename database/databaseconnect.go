package database

import (
	model "Kursach/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"os/exec"
	"time"
)

type Config struct {
	host     string
	user     string
	password string
	dbname   string
	port     string
	sslmode  string
}

var DB *gorm.DB

func Init() *gorm.DB {
	var cfg Config = Config{
		host:     os.Getenv("DB_HOST"),
		user:     os.Getenv("DB_USER"),
		password: os.Getenv("DB_PASSWORD"),
		dbname:   os.Getenv("DB_NAME"),
		port:     os.Getenv("DB_PORT"),
		sslmode:  " sslmode=disable",
	}
	dsn := "host=" + cfg.host + " user=" + cfg.user + " password=" + cfg.password + " dbname=" + cfg.dbname + cfg.sslmode
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&model.User{}, &model.Review{}, &model.Film{}, &model.Score{}, &model.Genre{})

	return db
}

func backupDatabase() error {
	var cfg Config = Config{
		host:     os.Getenv("DB_HOST"),
		user:     os.Getenv("DB_USER"),
		password: os.Getenv("DB_PASSWORD"),
		dbname:   os.Getenv("DB_NAME"),
		port:     os.Getenv("DB_PORT"),
		sslmode:  " sslmode=disable",
	}
	// Укажите путь для сохранения бэкапа
	backupFile := fmt.Sprintf("backup_%s.sql", time.Now().Format("20060102_150405"))

	// Формируем команду для выполнения pg_dump
	cmd := exec.Command("pg_dump", "-h", cfg.host, "-p", cfg.port, "-U", cfg.user, "-F", "c", "-b", "-v", "-f", backupFile, cfg.dbname)

	// Устанавливаем переменную окружения для пароля
	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", cfg.password))

	// Выполняем команду
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ошибка выполнения pg_dump: %v, вывод: %s", err, string(output))
	}

	return nil
}
