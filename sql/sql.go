package sql

import (
	"errors"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"log"
	"nextsure/conf"
	"nextsure/snapshot"
	"os"
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

func GetNoImg() (bool, conf.Link) {
	var noImgLink conf.Link
	db, _ := openDB()
	result := db.Where("img = ?", "loading").First(&noImgLink)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, noImgLink
	}
	return true, noImgLink
}

func DelLink(link conf.Link) {
	db, err := openDB()
	if err != nil {
		log.Fatalln(err)
	}
	var result = db.First(&link, "lid = ?", link.Lid)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Fatalln("Lid不存在")
	} else if result.Error != nil {
		log.Fatalln(result.Error)
	}

	result = db.Delete(&link)
	os.Remove("images/" + snapshot.FileName(link.Url) + ".png")
	if result.Error != nil {
		log.Fatalln(result.Error)
	}
}

func ChangeImg(link conf.Link) {
	db, err := openDB()
	if err != nil {
		log.Fatalln(err)
	}
	var findLink conf.Link
	var result = db.First(&findLink, "lid = ?", link.Lid)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Fatalln("Lid不存在")
	} else if result.Error != nil {
		log.Fatalln(result.Error)
	}
	findLink.Img = link.Img
	result = db.Save(&findLink)
	if result.Error != nil {
		log.Fatalln(result.Error)
	}
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
