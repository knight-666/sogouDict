package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/guonaihong/gout"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var baseUrl = "https://pinyin.sogou.com"
var path = "dist"
var recommend = false

func main() {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	header := http.Header{}
	body := ""
	code := -1

	err = gout.GET(baseUrl + "/dict/cate/index").
		BindHeader(&header).
		BindBody(&body).
		Code(&code).Do()
	if err != nil {
		panic(err)
	}

	reader, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		panic(err)
	}

	reader.Find("li.nav_list").Each(func(i int, selection *goquery.Selection) {
		href, exists := selection.Find("a").Attr("href")
		if exists {
			cate(href)
		}
	})

	fmt.Println("下载完成")
}

func cate(cate string) {
	header := http.Header{}
	body := ""
	code := -1

	err := gout.GET(baseUrl + cate).
		BindHeader(&header).
		BindBody(&body).
		Code(&code).Do()
	if err != nil {
		panic(err)
	}

	reader, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		panic(err)
	}

	page := 0
	reader.Find("div#dict_page_list").Find("li").Each(func(i int, selection *goquery.Selection) {
		text := selection.Find("a").Text()
		num := regexp.MustCompile("\\d+").FindAllString(text, -1)

		if len(num) == 1 {
			numTemp, _ := strconv.Atoi(num[0])
			if numTemp > page {
				page = numTemp
			}
		}
	})

	fmt.Println(reader.Find("div.cate_title").Text(), page, "页")

	download(reader)

	pageData(cate, page)
}

func pageData(cate string, page int) {
	for i := 2; i <= page; i++ {
		header := http.Header{}
		body := ""
		code := -1

		err := gout.GET(fmt.Sprintf("%s%d", baseUrl+cate+"/default/", i)).
			BindHeader(&header).
			BindBody(&body).
			Code(&code).Do()
		if err != nil {
			panic(err)
		}

		reader, err := goquery.NewDocumentFromReader(strings.NewReader(body))
		if err != nil {
			panic(err)
		}

		download(reader)
	}
}

func download(reader *goquery.Document) {
	reader.Find("div.dict_dl_btn").Each(func(i int, selection *goquery.Selection) {
		href, exists := selection.Find("a").Attr("href")
		if exists {

			parse, err := url.Parse(href)
			if err != nil {
				panic(err)
			}

			name := parse.Query()["name"][0]

			if recommend && !strings.Contains(name, "官方推荐") {
				return
			}

			var fileByte []byte
			err = gout.GET(href).BindBody(&fileByte).Do()
			if err != nil {
				panic(err)
			}

			file, err := os.Create(path + "/" + getFileName(name))
			if err != nil {
				panic(err)
			}
			defer file.Close()

			_, err = io.WriteString(file, string(fileByte))
			if err != nil {
				panic(err)
			}
		}
	})
}

func getFileName(name string) string {

	fileName := regexp.MustCompile(`[\\\.\*\?\|/:"<>]`).ReplaceAllString(name, "_")
	suffix := ".scel"

	for {
		_, err := os.Stat(fileName + suffix)
		if err == nil {
			fileName += "(1)"
		}
		if os.IsNotExist(err) {
			break
		}
	}

	return fileName + suffix
}
