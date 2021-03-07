# 实现一个简单的爬虫，说说你的思路

* [实现一个简单的爬虫，说说你的思路](#实现一个简单的爬虫说说你的思路)
  * [1 抓取单个网页内容](#1-抓取单个网页内容)
  * [2 抓取单个网页，解析内容，打印](#2-抓取单个网页解析内容打印)
  * [3 批量抓取多个网页，解析内容，打印](#3-批量抓取多个网页解析内容打印)
  * [4 批量抓取多个网页，解析内容，输出到 README\.md](#4-批量抓取多个网页解析内容输出到-readmemd)

github：https://github.com/xie4ever/practice/tree/master/golang/crawler

阿里数据库月报：http://mysql.taobao.org/monthly/ 简直是一个理想的爬虫试验场。

## 1 抓取单个网页内容

golang简单api的使用。

```golang
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jackdanger/collectlinks"
)

func main() {
	url := "http://mysql.taobao.org/monthly/2014/08/"
	if err := download(url); err != nil {
		log.Fatal(err)
	}
}

func download(url string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	// 自定义Header
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	// 函数结束后关闭相关链接
	defer resp.Body.Close()

	links := collectlinks.All(resp.Body)
	for _, link := range links {
		fmt.Println(link)
	}
	return nil
}
```

## 2 抓取单个网页，解析内容，打印

学会使用goquery解析网页文件。

```golang
package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	url := "http://mysql.taobao.org/monthly/2014/08/"
	body, err := download(url)
	if err != nil {
		log.Fatal(err)
	}

	contentList, err := getContentList(body)
	if err != nil {
		log.Fatal(err)
	}

	for _, content := range contentList {
		fmt.Println(content)
	}
}

func download(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, errors.New("invalid status code")
	}

	dom, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	return dom, nil
}

func getContentList(dom *goquery.Document) ([]string, error) {
	var list []string
	dom.Find("a:contains(·)").Each(func(i int, selection *goquery.Selection) {
		href, _ := selection.Attr("href")
		list = append(list, fmt.Sprintf("%s%s %s", "http://mysql.taobao.org/", href, getTitle(selection.Text())))
	})

	return list, nil
}

func getTitle(str string) string {
	str = strings.Replace(str, "\n", "", -1)
	rs := []rune(str)

	var (
		firstIdx     int
		lastIdx      int
		count        int
		lastPointIdx int
	)

	for idx, r := range rs {
		if firstIdx == 0 && r != ' ' {
			firstIdx = idx
		}

		if r != ' ' {
			lastIdx = idx
		}

		if r == '·' {
			count++
			if count == 2 {
				lastPointIdx = idx
			}
		}
	}

	prefix := string(rs[firstIdx : lastPointIdx+2])
	prefix = strings.Replace(prefix, " ", "", -1)
	suffix := string(rs[lastPointIdx+2 : lastIdx+1])

	return fmt.Sprintf("%s%s", prefix, suffix)
}
```

## 3 批量抓取多个网页，解析内容，打印

批量处理。

```golang
package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	statusCodeError = errors.New("invalid status code")
)

func main() {
	var monthList []string
	monthMap := make(map[string][]string, 8)
	for i := 2014; i <= 2021; i++ {
		month := strconv.Itoa(i)
		monthList = append(monthList, month)
		monthMap[month] = []string{"01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12"}
	}

	for _, month := range monthList {
		idList, ok := monthMap[month]
		if !ok {
			continue
		}
		for _, id := range idList {
			url := fmt.Sprintf("http://mysql.taobao.org/monthly/%s/%s/", month, id)
			body, err := download(url)
			if errors.Is(err, statusCodeError) {
				continue
			}
			if err != nil {
				log.Fatal(err)
			}

			contentList, err := getContentList(body)
			if err != nil {
				log.Fatal(err)
			}

			for _, content := range contentList {
				fmt.Println(content)
			}

			time.Sleep(2 * time.Second)
		}
	}
}

func download(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, statusCodeError
	}

	dom, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	return dom, nil
}

func getContentList(dom *goquery.Document) ([]string, error) {
	var list []string
	dom.Find("a:contains(·)").Each(func(i int, selection *goquery.Selection) {
		href, _ := selection.Attr("href")
		list = append(list, fmt.Sprintf("%s%s %s", "http://mysql.taobao.org/", href, getTitle(selection.Text())))
	})

	return list, nil
}

func getTitle(str string) string {
	str = strings.Replace(str, "\n", "", -1)
	rs := []rune(str)

	var (
		firstIdx     int
		lastIdx      int
		count        int
		lastPointIdx int
	)

	for idx, r := range rs {
		if firstIdx == 0 && r != ' ' {
			firstIdx = idx
		}

		if r != ' ' {
			lastIdx = idx
		}

		if r == '·' {
			count++
			lastPointIdx = idx
		}
	}

	prefix := string(rs[firstIdx : lastPointIdx+2])
	prefix = strings.Replace(prefix, " ", "", -1)
	suffix := string(rs[lastPointIdx+2 : lastIdx+1])

	return fmt.Sprintf("%s%s", prefix, suffix)
}
```

## 4 批量抓取多个网页，解析内容，输出到 README.md

学习基本文件操作。

```golang
package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	statusCodeError = errors.New("invalid status code")
)

func main() {
	fileName := "README.md"
	if _, err := os.Stat(fileName); !os.IsNotExist(err) {
		if err := os.Remove(fileName); err != nil {
			log.Fatal(err)
		}
	}
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	if _, err := file.WriteString("|address|content|\n"); err != nil {
		log.Fatal(err)
	}
	if _, err := file.WriteString("|----|----|\n"); err != nil {
		log.Fatal(err)
	}

	var monthList []string
	monthMap := make(map[string][]string, 8)
	for i := 2014; i <= 2021; i++ {
		month := strconv.Itoa(i)
		monthList = append(monthList, month)
		monthMap[month] = []string{"01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12"}
	}

	for _, month := range monthList {
		idList, ok := monthMap[month]
		if !ok {
			continue
		}

		for _, id := range idList {
			url := fmt.Sprintf("http://mysql.taobao.org/monthly/%s/%s/", month, id)
			body, err := download(url)
			if errors.Is(err, statusCodeError) {
				continue
			}
			if err != nil {
				log.Fatal(err)
			}

			contentList, err := getContentList(body)
			if err != nil {
				log.Fatal(err)
			}

			for _, content := range contentList {
				text := fmt.Sprintf("%s\n", content)
				fmt.Print(text)
				if _, err := file.WriteString(text); err != nil {
					log.Fatal(err)
				}
			}

			time.Sleep(2 * time.Second)
		}
	}
}

func download(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, statusCodeError
	}

	dom, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	return dom, nil
}

func getContentList(dom *goquery.Document) ([]string, error) {
	var list []string
	dom.Find("a:contains(·)").Each(func(i int, selection *goquery.Selection) {
		href, _ := selection.Attr("href")
		list = append(list, fmt.Sprintf("|%s%s|%s|", "http://mysql.taobao.org/", href, getTitle(selection.Text())))
	})

	return list, nil
}

func getTitle(str string) string {
	str = strings.Replace(str, "\n", "", -1)
	rs := []rune(str)

	var (
		firstIdx     int
		lastIdx      int
		count        int
		lastPointIdx int
	)

	for idx, r := range rs {
		if firstIdx == 0 && r != ' ' {
			firstIdx = idx
		}

		if r != ' ' {
			lastIdx = idx
		}

		if r == '·' {
			count++
			lastPointIdx = idx
		}
	}

	prefix := string(rs[firstIdx : lastPointIdx+2])
	prefix = strings.Replace(prefix, " ", "", -1)
	suffix := string(rs[lastPointIdx+2 : lastIdx+1])

	return fmt.Sprintf("%s%s", prefix, suffix)
}
```
