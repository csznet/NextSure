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
)

func main() {
	//url := "https://www.csz.net"
	//snapshot.Get(url)
	//println(title(url))
	web()
	//addLink("https://www.baidu.com")
}

func addLink(url string) {
	var link conf.Link
	link.Url = url
	if !sql.ExistLink(link) {
		link.Title = title(url)
		if snapshot.Get(url) {
			link.Img = snapshot.FileName(url)
		}
		sql.NewLink(link)
	}
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
			"link": sql.GetLink(1, 10),
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
