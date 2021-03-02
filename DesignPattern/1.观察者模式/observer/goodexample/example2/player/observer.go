package player

import (
	"sync"
)

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
	v, ok := o.ProcessorMap.Load(p.GetType())
	if ok {
		pList := v.([]processor)
		pList = append(pList, p)
		o.ProcessorMap.Store(p.GetType(), pList)
	} else {
		o.ProcessorMap.Store(p.GetType(), []processor{p})
	}
}

// DeleteProcessorByType ...
func (o *observer) DeleteProcessorByType(t int) {
	o.ProcessorMap.Delete(t)
}

// PostEvent ...
func (o *observer) PostEvent(e Event) error {
	if e.ID == 0 {
		return nil
	}
	v, ok := o.ProcessorMap.Load(e.Type)
	if ok {
		pList := v.([]processor)
		for _, p := range pList {
			if err := p.DoSomething(e); err != nil {
				return err
			}
		}
	}
	return nil
}
