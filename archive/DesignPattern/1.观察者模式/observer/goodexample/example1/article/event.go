package article

import (
	"time"
)

const (
	TypeEventAdd    = 1
	TypeEventDelete = 2
	TypeEventModify = 3
)

// Event ...
type Event struct {
	ID      int64                  `json:"id"`
	Type    int                    `json:"type"`
	Article Article                `json:"article"`
	Time    int64                  `json:"time"`
	Param   map[string]interface{} `json:"param"`
}

// NewEvent ...
func NewEvent(t int, article Article, param map[string]interface{}) Event {
	return Event{
		ID:      time.Now().Unix(),
		Type:    t,
		Article: article,
		Time:    time.Now().Unix(),
		Param:   param,
	}
}
