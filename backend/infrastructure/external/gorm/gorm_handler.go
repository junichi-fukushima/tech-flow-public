package gorm

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DataBaseConfig struct {
	Host       string
	Database   string
	Port       string
	Driver     string
	User       string
	Password   string
	CtxTimeout time.Duration
}

type SQLHandler struct {
	DB *gorm.DB
}

var (
	sqlHandlerInstance *SQLHandler
	sqlHandlerOnce     sync.Once
)

// 初期化を一度だけ行う
func GetSQLHandler() (*SQLHandler, error) {
	var err error
	sqlHandlerOnce.Do(func() {
		cfg := NewSQLHandlerConfig()
		sqlHandlerInstance, err = NewSQLHandler(cfg)
	})
	return sqlHandlerInstance, err
}

func NewSQLHandlerConfig() *DataBaseConfig {
	return &DataBaseConfig{
		Host:       os.Getenv("DB_HOST"),
		Database:   os.Getenv("DB_NAME"),
		Port:       os.Getenv("DB_PORT"),
		Driver:     "mysql",
		User:       os.Getenv("DB_USERNAME"),
		Password:   os.Getenv("DB_PASSWORD"),
		CtxTimeout: 60 * time.Second,
	}
}

// SQL接続を初期化
func NewSQLHandler(cfg *DataBaseConfig) (*SQLHandler, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &SQLHandler{DB: db}, nil
}

// Select(全取得想定)
func (handler *SQLHandler) Select(dest interface{}, query string, args ...interface{}) error {
	if err := handler.DB.Raw(query, args...).Scan(dest).Error; err != nil {
		return err
	}
	return nil
}

func (handler *SQLHandler) Get(dest interface{}, query string, args ...interface{}) error {
	if err := handler.DB.Raw(query, args...).Scan(dest).Error; err != nil {
		return err
	}
	return nil
}

// Insert: レコードの挿入
func (handler *SQLHandler) Insert(value interface{}) error {
	return handler.DB.Create(value).Error
}

// Update: レコードの更新
func (handler *SQLHandler) Update(value interface{}) error {
	return handler.DB.Save(value).Error
}

// First: 単一レコードの取得
func (handler *SQLHandler) First(dest interface{}, conds ...interface{}) error {
	return handler.DB.First(dest, conds...).Error
}

// Where: 条件付きでレコードを取得
func (handler *SQLHandler) Where(query interface{}, args ...interface{}) *gorm.DB {
	return handler.DB.Where(query, args...)
}

// IsRecordNotFound: ErrRecordNotFound の判定
func (handler *SQLHandler) IsRecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
