package util

import (
	"context"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

// Redis客户端
var RedisClient *redis.Client

// InitRedis 初始化Redis客户端
func InitRedis(fileName string) error {
	// 从配置文件中读取Redis配置
	// 调用loadEnv函数加载配置文件
	loadEnv(fileName)

	// 从环境变量中读取Redis配置
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")

	redisDB := 0
	redisDBStr := os.Getenv("REDIS_DB")
	if redisDBStr != "" {
		redisDB, _ = strconv.Atoi(redisDBStr)
	}

	// 创建Redis客户端
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,     // Redis服务器地址
		Password: redisPassword, // Redis密码
		DB:       redisDB,       // Redis数据库
	})

	// 测试Redis连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return err
	}

	return nil
}

// SetUserSession 设置用户会话
func SetUserSession(userID uint, token string) error {
	if RedisClient == nil {
		return errors.New("redis client not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 设置用户会话，过期时间为30分钟
	err := RedisClient.Set(ctx, "user:"+strconv.Itoa(int(userID)), token, 30*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

// GetUserSession 获取用户会话
func GetUserSession(userID uint) (string, error) {
	if RedisClient == nil {
		return "", errors.New("redis client not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 获取用户会话
	token, err := RedisClient.Get(ctx, "user:"+strconv.Itoa(int(userID))).Result()
	if err != nil {
		return "", err
	}

	return token, nil
}

// DeleteUserSession 删除用户会话
func DeleteUserSession(userID uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 删除用户会话
	err := RedisClient.Del(ctx, "user:"+strconv.Itoa(int(userID))).Err()
	if err != nil {
		return err
	}

	return nil
}
