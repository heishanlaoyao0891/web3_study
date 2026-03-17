package test

import (
	"fmt"
	"go-blog/util"
	"testing"
)

// 测试类标准写法
func TestMysqlConnect(t *testing.T) {
	err := util.CreateTestDB(".env_local")
	if err != nil {
		t.Fatalf("数据库连接失败: %v", err)
	}

	fmt.Println("数据库连接成功！")

}
