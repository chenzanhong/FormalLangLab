package repository

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// 连接PostgreSQL
func ConnectPG() error {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	DB, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return err
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接的最大生命周期

	return nil
}

func DBInit() error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	migrationsDir := filepath.Join(".", "migrations")

	var err error
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %v", err)
	}

	tx := DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	for _, file := range files {
		// 执行以.init.sql结尾的文件
		if filepath.Ext(file.Name()) != ".init.sql" {
			continue
		}
		// 读取文件内容
		content, err := ioutil.ReadFile(filepath.Join(migrationsDir, file.Name()))
		if err != nil {
			return fmt.Errorf("failed to read file: %v", err)
		}
		// 执行SQL语句
		if err = tx.Exec(string(content)).Error; err != nil {
			return fmt.Errorf("failed to execute SQL: %v", err)
		}
	}

	return nil
}
