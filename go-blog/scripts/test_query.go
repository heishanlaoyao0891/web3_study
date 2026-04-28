package main

import (
	"fmt"
	"go-blog/model"
	"go-blog/util"
)

func main() {
	if err := util.CreateTestDB(".env_local"); err != nil {
		panic(err)
	}

	var questions []model.InterviewQuestion
	result := util.Db.Order("id desc").Find(&questions)
	
	fmt.Printf("查询到 %d 条记录\n", len(questions))
	if result.Error != nil {
		fmt.Println("错误:", result.Error)
	}
	
	for i, q := range questions {
		fmt.Printf("%d. ID:%d [%s] %s\n", i+1, q.ID, q.Category, q.Title)
	}
}
