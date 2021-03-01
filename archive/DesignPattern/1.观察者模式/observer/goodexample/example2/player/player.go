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
	event := NewEvent(TypeLevelUp, p.ID, nil)
	if err := GetObs().PostEvent(event); err != nil {
		return err
	}
	return nil
}
