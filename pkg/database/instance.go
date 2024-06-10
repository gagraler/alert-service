package database

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

/**
 * @author: gagral.x@gmail.com
 * @time: 2023/11/10 23:48
 * @file: instance.go
 * @description: 数据库连接
 */

var DB *gorm.DB

// Config DataBaseConfig GORM
type Config struct {
	DSN         string
	MaxOpenConn int
	MaxIdleConn int
	MaxLifetime int
	MaxIdleTime int
}

func NewDatabase(c *Config) (*gorm.DB, error) {

	db, err := gorm.Open(mysql.Open(c.DSN), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %s", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %s", err)
	}

	// 设置连接池大小
	sqlDB.SetMaxOpenConns(c.MaxOpenConn)
	sqlDB.SetMaxIdleConns(c.MaxIdleConn)
	sqlDB.SetConnMaxLifetime(time.Duration(c.MaxLifetime) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(c.MaxIdleTime) * time.Second)

	if err != nil {
		return nil, err
	}

	DB = db

	return db, err
}
