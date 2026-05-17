package database

import (
	"cardapio-digital/internal/models"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect() (*gorm.DB, error) {
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
		return nil, fmt.Errorf("não foi possível conectar ao banco de dados: %w", err)
	}
	log.Println("Deus abençoe: a conexão com o banco de dados foi estabelecida com sucesso")

	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(10)           
		sqlDB.SetMaxOpenConns(100)          
		sqlDB.SetConnMaxLifetime(time.Hour)
	}
	// Aplica as migrações automáticas para criar as tabelas no banco de dados
	if err := RunMigrations(db); err != nil {
		return nil, err
	}

	return db, nil
}

// RunMigrations é responsável por aplicar as migrações automáticas para criar as tabelas no banco de dados
func RunMigrations(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.Restaurant{},
		&models.Category{},
		&models.Item{},
	)

	if err != nil {
		log.Printf("Fogo no parquinho: as migrações automáticas falharam: %v", err)
		return err
	}

	log.Println("Deus abençoe: as migrações automáticas foram aplicadas com sucesso")
	return nil
}