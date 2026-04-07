package database

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Config is subsection of config.
type Config struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	Database        string `yaml:"database"`
	Charset         string `yaml:"charset"`
	ParseTime       bool   `yaml:"parse_time"`
	Loc             string `yaml:"loc"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`  // seconds
	ConnMaxIdleTime int    `yaml:"conn_max_idle_time"` // seconds
}

// NewClient init database connection
func NewClient(conf Config) (*gorm.DB, error) {
	dsn, err := formatDSN(conf)
	if err != nil {
		return nil, err
	}
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return nil, err
	}

	sqldb, err := db.DB()
	if err != nil {
		return nil, err
	}
	err = sqldb.Ping()
	if err != nil {
		return nil, err
	}

	maxOpen := conf.MaxOpenConns
	if maxOpen == 0 {
		maxOpen = 100
	}
	maxIdle := conf.MaxIdleConns
	if maxIdle == 0 {
		maxIdle = 50
	}
	connMaxLifetime := conf.ConnMaxLifetime
	if connMaxLifetime == 0 {
		connMaxLifetime = 300
	}
	connMaxIdleTime := conf.ConnMaxIdleTime
	if connMaxIdleTime == 0 {
		connMaxIdleTime = 60
	}
	sqldb.SetMaxOpenConns(maxOpen)
	sqldb.SetMaxIdleConns(maxIdle)
	sqldb.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)
	sqldb.SetConnMaxIdleTime(time.Duration(connMaxIdleTime) * time.Second)

	return db, nil
}

// formatDSN 获取 DSN
func formatDSN(conf Config) (string, error) {
	if conf.Host == "" || conf.Port == 0 || conf.Username == "" || conf.Password == "" || conf.Database == "" {
		return "", errors.New("db config should not be empty")
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s", conf.Username, conf.Password, conf.Host,
		conf.Port, conf.Database, conf.Charset, conf.ParseTime, conf.Loc)
	return dsn, nil
}
