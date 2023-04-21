package sql

import (
	"errors"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"log"
	"nextsure/conf"
	"time"
)

// 打开数据库
func openDB() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
}

// ExistLink 判断链接是否存在
func ExistLink(link conf.Link) bool {
	db, err := openDB()
	if err != nil {
		log.Fatalf("无法打开数据库: %v\n", err)
		return false
	}
	db.AutoMigrate(&conf.Link{})
	var existingLink conf.Link
	result := db.Where("url = ?", link.Url).First(&existingLink)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}

// NewLink 添加链接收藏
func NewLink(newLink conf.Link) bool {
	db, err := openDB()
	if err != nil {
		log.Fatalf("无法打开数据库: %v\n", err)
		return false
	}
	//生成相应的表
	err = db.AutoMigrate(&conf.Link{})
	if err != nil {
		log.Fatalf("无法迁移表: %v\n", err)
		return false
	}
	newLink.Date = time.Now().Unix()
	if !ExistLink(newLink) {
		result := db.Create(&newLink)
		if result.Error != nil {
			log.Fatalf("添加链接失败: %v\n", result.Error)
			return false
		}
	}
	return true
}

func GetLink(page, pageSize int) []conf.Link {
	var links []conf.Link
	db, _ := openDB()
	offset := (page - 1) * pageSize
	result := db.Offset(offset).Limit(pageSize).Find(&links)
	if result.Error != nil {
		return nil
	}
	return links
}
