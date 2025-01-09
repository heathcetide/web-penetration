package configs

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func InitDB() (*gorm.DB, error) {
	dbConfig := &DBConfig{
		Host:     "localhost",
		Port:     "3306",
		User:     "root",
		Password: "1234",
		DBName:   "pentest_db",
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
	)

	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
