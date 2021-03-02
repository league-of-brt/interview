package article

import (
	"log"
	"time"
)

// Article ...
type Article struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// NewArticle ...
func NewArticle(title, content string) Article {
	return Article{
		ID:      time.Now().Unix(),
		Title:   title,
		Content: content,
	}
}

// Add ...
func (a Article) Add() error {
	log.Print("do something")
	log.Print("rank change")
	log.Print("point change")
	return nil
}

// Delete ...
func (a Article) Delete() error {
	log.Print("do something")
	log.Print("rank change")
	log.Print("point change")
	return nil
}

// Modify ...
func (a Article) Modify() error {
	log.Print("do something")
	log.Print("rank change")
	log.Print("point change")
	return nil
}
