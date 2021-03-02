package player

import "log"

type processor interface {
	GetType() int
	DoSomething(e Event) error
}

// LevelUpRankProcessor ...
type LevelUpRankProcessor struct {
	Type int `json:"type"`
}

// GetType ...
func (p LevelUpRankProcessor) GetType() int {
	return p.Type
}

// DoSomething ...
func (p LevelUpRankProcessor) DoSomething(e Event) error {
	log.Printf("LevelUpRankProcessor get the event: %d", e.ID)
	return nil
}

// LevelUpRewardProcessor ...
type LevelUpRewardProcessor struct {
	Type int `json:"type"`
}

// GetType ...
func (p LevelUpRewardProcessor) GetType() int {
	return p.Type
}

// DoSomething ...
func (p LevelUpRewardProcessor) DoSomething(e Event) error {
	log.Printf("LevelUpRewardProcessor get the event: %d", e.ID)
	return nil
}

// LevelUpAnnounceProcessor ...
type LevelUpAnnounceProcessor struct {
	Type int `json:"type"`
}

// GetType ...
func (p LevelUpAnnounceProcessor) GetType() int {
	return p.Type
}

// DoSomething ...
func (p LevelUpAnnounceProcessor) DoSomething(e Event) error {
	log.Printf("LevelUpAnnounceProcessor get the event: %d", e.ID)
	return nil
}
