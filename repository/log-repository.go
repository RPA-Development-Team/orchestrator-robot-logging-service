package repository

import (
	"os"

	"github.com/khalidzahra/robot-logging-service/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type LogRepository interface {
	Save(log entity.Log)
	Update(log entity.Log)
	Delete(log entity.Log)
	FindAll() []entity.Log
	CloseDB()
}

type database struct {
	connection *gorm.DB
}

func NewLogRepository() LogRepository {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	db.AutoMigrate(&entity.Log{})
	return &database{
		connection: db.Set("gorm:auto_preload", true),
	}
}
func (db *database) CloseDB() {
	sqlDB, _ := db.connection.DB()
	err := sqlDB.Close()
	if err != nil {
		panic("Failed to close database")
	}
}

func (db *database) Save(log entity.Log) {
	db.connection.Create(&log)
}

func (db *database) Update(log entity.Log) {
	db.connection.Save(&log)
}

func (db *database) Delete(log entity.Log) {
	db.connection.Delete(&log)
}

func (db *database) FindAll() []entity.Log {
	var logs []entity.Log
	db.connection.Find(&logs)
	return logs
}
