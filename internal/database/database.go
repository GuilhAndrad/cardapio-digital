package database

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect() *gorm.DB {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	db, err := gorm.Open(sqlite.Open("cardapio.db"), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatal("Fogo no parquinho: não foi possível conectar ao banco de dados: ", err)
	}
	log.Println("Deus abençoe: a conexão com o banco de dados foi estabelecida com sucesso")

	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(10)           
		sqlDB.SetMaxOpenConns(100)          
		sqlDB.SetConnMaxLifetime(time.Hour)
	}
	return db
}