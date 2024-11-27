package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"ticker-tracer/config"
	"ticker-tracer/model/entity"
	"time"
)

var dbClient *gorm.DB

func InitDb() error {
	var err error

	cfg := config.GetConfig()
	cnn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Europe/Istanbul",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password,
		cfg.Postgres.DbName, cfg.Postgres.SSLMode)

	dbClient, err = gorm.Open(postgres.Open(cnn), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDb, _ := dbClient.DB()
	err = sqlDb.Ping()
	if err != nil {
		return err
	}

	sqlDb.SetMaxIdleConns(cfg.Postgres.MaxIdleConns)
	sqlDb.SetMaxOpenConns(cfg.Postgres.MaxOpenConns)
	sqlDb.SetConnMaxLifetime(cfg.Postgres.ConnMaxLifetime * time.Minute)

	err = AutoMigrate()
	if err != nil {
		return err
	}

	log.Println("Db connection established")
	return nil
}
func AutoMigrate() error {

	err := dbClient.AutoMigrate(&entity.TicketRequest{}, &entity.TicketRequest{})
	if err != nil {
		_ = fmt.Errorf("error migrating db: %v", err)
		return err
	}
	return nil

}

func GetDb() *gorm.DB {
	return dbClient
}

func CloseDb() {
	con, _ := dbClient.DB()
	con.Close()
}
