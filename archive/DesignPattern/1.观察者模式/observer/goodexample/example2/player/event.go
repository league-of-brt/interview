package player

import "time"

const TypeLevelUp = 1

// Event ...
type Event struct {
	ID       int64                  `json:"id"`
	Type     int                    `json:"type"`
	PlayerID int64                  `json:"player_id"`
	Time     int64                  `json:"time"`
	Param    map[string]interface{} `json:"param"`
}

// NewEvent ...
func NewEvent(t int, playerID int64, param map[string]interface{}) Event {
	return Event{
		ID:       time.Now().Unix(),
		Type:     t,
		PlayerID: playerID,
		Time:     time.Now().Unix(),
		Param:    param,
	}
}
