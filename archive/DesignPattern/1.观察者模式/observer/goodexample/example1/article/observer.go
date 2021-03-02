package article

import "sync"

var (
	obs observer
)

type observer struct {
	ProcessorMap sync.Map `json:"processor_map"`
}

// GetObs ...
func GetObs() *observer {
	return &obs
}

// AddProcessor ...
func (o *observer) AddProcessor(p processor) {
	o.ProcessorMap.Store(p.GetID(), p)
}

// DeleteProcessor ...
func (o *observer) DeleteProcessor(id int64) {
	o.ProcessorMap.Delete(id)
}

// PostEvent ...
func (o *observer) PostEvent(e Event) error {
	if e.ID == 0 {
		return nil
	}
	if e.Type == 0 {
		return nil
	}

	o.ProcessorMap.Range(
		func(k, v interface{}) bool {
			p := v.(processor)
			if e.Type == TypeEventAdd {
				if err := p.EntryAdded(e); err != nil {
					return false
				}
			}
			if e.Type == TypeEventDelete {
				if err := p.EntryDeleted(e); err != nil {
					return false
				}
			}
			if e.Type == TypeEventModify {
				if err := p.EntryModified(e); err != nil {
					return false
				}
			}
			return true
		})

	return nil
}
