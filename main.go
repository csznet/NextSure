package main

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"nextsure/conf"
	"nextsure/snapshot"
	"nextsure/sql"
	"regexp"
	"strconv"
	"time"
)

func main() {
	//snapshot.Get(url)
	//println(title(url))
	go getImg()
	web()
}

func noImg(Lid uint) {
	var link conf.Link
	link.Lid = Lid
	link.Img = "loading"
	sql.ChangeImg(link)
}

func delLink(Lid uint) {
	var link conf.Link
	link.Lid = Lid
	sql.DelLink(link)
}

func getImg() {
	ok, link := sql.GetNoImg()
	if ok {
		if snapshot.Get(link.Url) {
			link.Img = snapshot.FileName(link.Url)
			sql.ChangeImg(link)
		}
	}
	time.AfterFunc(10*time.Second, getImg)
}

func addLink(url string) {
	var link conf.Link
	link.Url = url
	link.Title = title(url)
	link.Img = "loading"
	//if !sql.ExistLink(link) {
	//	link.Title = title(url)
	//	if snapshot.Get(url) {
	//		link.Img = snapshot.FileName(url)
	//	}
	sql.NewLink(link)
	//}
}

func title(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	re := regexp.MustCompile(`(?i)<title>(.*?)</title>`)
	title := "未知网站"
	if len(re.FindStringSubmatch(string(body))) > 1 {
		title = re.FindStringSubmatch(string(body))[1]
	}
	return title
}

// Web Serve
func web() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.StaticFS("/images", http.Dir("./images"))
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"page": 1,
			"link": sql.GetLink(1, 12),
		})
	})
	r.GET("/del/:lid", func(c *gin.Context) {
		Lid, err := strconv.Atoi(c.Param("lid"))
		msg := "删除失败"
		if err == nil {
			msg = "删除成功"
			delLink(uint(Lid))
		}
		c.HTML(http.StatusOK, "msg.tmpl", gin.H{
			"msg": msg,
			"url": "close",
		})
	})
	r.GET("/ref/:lid", func(c *gin.Context) {
		Lid, err := strconv.Atoi(c.Param("lid"))
		msg := "刷新快照失败"
		if err == nil {
			msg = "已成功加入队列"
			noImg(uint(Lid))
		}
		c.HTML(http.StatusOK, "msg.tmpl", gin.H{
			"msg": msg,
			"url": "close",
		})
	})
	r.GET("/page/:s", func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		page, err := strconv.Atoi(c.Param("s"))
		if err != nil || page < 1 {
			page = 1
		}
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"page": page,
			"link": sql.GetLink(page, 12),
		})
	})
	r.POST("/add", func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		log.Println(c.PostForm("url"))
		msg := "请输入网址"
		if len(c.PostForm("url")) > 0 {
			go addLink(c.PostForm("url"))
			msg = "成功添加到队列:" + c.PostForm("url")
		}
		c.HTML(http.StatusOK, "msg.tmpl", gin.H{
			"msg": msg,
			"url": "/",
		})
	})
	r.Run("0.0.0.0:8088")
}
