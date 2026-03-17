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

// loadEnv读取.env文件里面的变量
func loadEnv(fileName string) {
	//确保代码块只执行一次，即使在并发场景下也能保证线程安全
	envOnce.Do(func() {
		//获取当前文件所在的目录 获取当前执行代码的位置
		_, file, _, ok := runtime.Caller(0)
		if !ok {
			return
		}
		//多次调用获取父目录，实现相对路径的准确定位
		dir := filepath.Dir(file)
		s := filepath.Dir(dir)

		//加载.env文件
		envPath := filepath.Join(s, fileName)
		if err := godotenv.Load(envPath); err != nil {
			return
		}
	})
}

// 创建数据库连接
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
	sqlDb.SetMaxIdleConns(2)                   // 保持2个空闲连接就绪
	sqlDb.SetMaxOpenConns(5)                   // 允许最多5个并发连接
	sqlDb.SetConnMaxLifetime(30 * time.Minute) // 连接最多重用30分钟
	return nil
}

// 连接字符串从TEST_MYSQL_DSN环境变量或.env文件读取
// 格式: user:password@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local
func newMySQLDB(fileName string) (*gorm.DB, error) {
	loadEnv(fileName)
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		dsn = "root:password@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
	}

	return gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Logger: 设置为logger.Info以在开发中查看所有SQL查询
		Logger: logger.Default.LogMode(logger.Info),

		// NamingStrategy: 自定义表和列命名
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",
			SingularTable: false,
			NoLowerCase:   false,
		},
	})
}
