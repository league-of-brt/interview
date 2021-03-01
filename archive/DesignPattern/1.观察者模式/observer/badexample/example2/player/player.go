package player

import (
	"log"
	"time"
)

// Player ...
type Player struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// NewPlayer ...
func NewPlayer(name string) Player {
	return Player{
		ID:   time.Now().Unix(),
		Name: name,
	}
}

// Attack ...
func (p Player) Attack() error {
	log.Print("attack somebody")
	log.Print("level up")
	log.Print("level up reward change")
	log.Print("level rank change")
	log.Print("level up announce")
	return nil
}
