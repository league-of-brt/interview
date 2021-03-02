package main

import (
	"time"

	"design-pattern/observer/badexample/example1/article"
)

func main() {
	a := article.NewArticle("title", "content")
	_ = a.Add()
	time.Sleep(time.Second)
	_ = a.Modify()
	time.Sleep(time.Second)
	_ = a.Delete()
}
