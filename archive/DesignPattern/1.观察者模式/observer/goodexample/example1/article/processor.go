package article

import "log"

type processor interface {
	GetID() int64
	EntryAdded(e Event) error
	EntryDeleted(e Event) error
	EntryModified(e Event) error
}

// RankProcessor ...
type RankProcessor struct {
	ID int64 `json:"id"`
}

// GetID ...
func (p RankProcessor) GetID() int64 {
	return p.ID
}

// EntryAdded ...
func (p RankProcessor) EntryAdded(e Event) error {
	log.Printf("rankProcessor EntryAdded, event id :%d", e.ID)
	return nil
}

// EntryDeleted ...
func (p RankProcessor) EntryDeleted(e Event) error {
	log.Printf("rankProcessor EntryDeleted, event id :%d", e.ID)
	return nil
}

// EntryModified ...
func (p RankProcessor) EntryModified(e Event) error {
	log.Printf("rankProcessor EntryModified, event id :%d", e.ID)
	return nil
}

// PointProcessor ...
type PointProcessor struct {
	ID int64 `json:"id"`
}

// GetID ...
func (p PointProcessor) GetID() int64 {
	return p.ID
}

// EntryAdded ...
func (p PointProcessor) EntryAdded(e Event) error {
	log.Printf("pointProcessor EntryAdded, event id :%d", e.ID)
	return nil
}

// EntryDeleted ...
func (p PointProcessor) EntryDeleted(e Event) error {
	log.Printf("pointProcessor EntryAdded, event id :%d", e.ID)
	return nil
}

// EntryModified ...
func (p PointProcessor) EntryModified(e Event) error {
	log.Printf("pointProcessor EntryAdded, event id :%d", e.ID)
	return nil
}
