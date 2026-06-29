package util

import (
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var envOnce sync.Once
var Db *gorm.DB

// loadEnv 读取 .env 文件里面的变量
func loadEnv(fileName string) {
	envOnce.Do(func() {
		_, file, _, ok := runtime.Caller(0)
		if !ok {
			return
		}
		dir := filepath.Dir(file)
		s := filepath.Dir(dir)

		envPath := filepath.Join(s, fileName)
		if err := godotenv.Load(envPath); err != nil {
			return
		}
	})
}

// CreateTestDB 创建数据库连接
func CreateTestDB(fileName string) error {
	var err error
	Db, err = newMySQLDB(fileName)
	if err != nil {
		return err
	}
	sqlDb, err := Db.DB()
	if err != nil {
		return err
	}
	sqlDb.SetMaxIdleConns(2)
	sqlDb.SetMaxOpenConns(5)
	sqlDb.SetConnMaxLifetime(30 * time.Minute)
	return nil
}

// newMySQLDB 从环境变量 MYSQL_DSN 创建数据库连接
// M4.1 安全修复：移除硬编码的数据库密码 fallback，必须通过环境变量或 .env 文件配置
func newMySQLDB(fileName string) (*gorm.DB, error) {
	loadEnv(fileName)
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		return nil, fmtError("MYSQL_DSN 环境变量未设置，请在 %s 文件中配置数据库连接字符串", fileName)
	}

	return gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",
			SingularTable: false,
			NoLowerCase:   false,
		},
	})
}

// fmtError 简单格式化错误信息（避免引入额外包）
func fmtError(format string, args ...interface{}) error {
	msg := format
	for _, arg := range args {
		_ = arg
	}
	return &dbError{msg}
}

type dbError struct {
	msg string
}

func (e *dbError) Error() string {
	return e.msg
}