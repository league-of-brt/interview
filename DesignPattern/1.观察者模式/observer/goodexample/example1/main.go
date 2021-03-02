package main

import (
	"time"

	"design-pattern/observer/goodexample/example1/article"
)

func main() {
	rankProcessor := article.RankProcessor{
		ID: 1,
	}
	pointProcessor := article.PointProcessor{
		ID: 2,
	}
	obs := article.GetObs()
	obs.AddProcessor(rankProcessor)
	obs.AddProcessor(pointProcessor)

	a := article.NewArticle("title", "content")
	_ = a.Add()
	time.Sleep(time.Second)
	_ = a.Modify()
	time.Sleep(time.Second)
	_ = a.Delete()
}
