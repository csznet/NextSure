package snapshot

import (
	"context"
	"io"
	"log"
	"os"
	"strings"

	cdp "github.com/chromedp/chromedp"
)

func Get(urlStr string) bool {
	// 创建新的cdp上下文
	ctx, cancel := cdp.NewContext(context.Background())
	defer cancel()

	// 此处以360搜索首页为例
	//urlStr := `https://www.csz.net`

	domain := FileName(urlStr)

	var buf []byte
	log.Println("正在打开网页")
	if err := cdp.Run(ctx, cdp.EmulateViewport(1920, 0), fullScreenshot(urlStr, 90, &buf)); err != nil {
		log.Fatal(err)
	}
	//判断文件夹是否存在
	dir := "images"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0755); err != nil {
			panic(err)
		}
	}
	ok := true
	// 写入文件
	file, err := os.Create("images/" + domain + ".png")
	if err != nil {
		panic(err)
		ok = false
	}
	defer file.Close()

	_, err = io.WriteString(file, string(buf))
	if err != nil {
		panic(err)
		ok = false
	}
	return ok
}

func FileName(urlStr string) string {
	// 提取域名并替换 "." 为 "-"
	domainOld := strings.ReplaceAll(strings.Split(urlStr, "//")[1], ".", "-")
	domainStr := strings.Split(domainOld, "/")
	domain := domainStr[0]
	domain = strings.ReplaceAll(domain, "/", "")
	if len(domainStr) > 0 {
		domainOld = strings.ReplaceAll(domainOld, "/", "_")
		domain = domain + strings.TrimLeft(domainOld, domain)
	}
	return domain
}

func fullScreenshot(urlstr string, quality int, res *[]byte) cdp.Tasks {
	return cdp.Tasks{
		cdp.Navigate(urlstr),
		cdp.FullScreenshot(res, quality),
	}
}
