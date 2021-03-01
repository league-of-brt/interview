package player

import (
	"sync"
)

var (
	obs      *observer
	obsMutex sync.Mutex
)

type observer struct {
	ProcessorMap sync.Map `json:"processor_map"`
}

// GetObs ...
func GetObs() *observer {
	obsMutex.Lock()
	if obs == nil {
		obs = &observer{
			ProcessorMap: sync.Map{},
		}
	}
	obsMutex.Unlock()
	return obs
}

// AddProcessor ...
func (o *observer) AddProcessor(p processor) {
	obsMutex.Lock()
	v, ok := o.ProcessorMap.Load(p.GetType())
	if ok {
		pList := v.([]processor)
		pList = append(pList, p)
		o.ProcessorMap.Store(p.GetType(), pList)
	} else {
		o.ProcessorMap.Store(p.GetType(), []processor{p})
	}
	obsMutex.Unlock()
}

// DeleteProcessorByType ...
func (o *observer) DeleteProcessorByType(t int) {
	obsMutex.Lock()
	o.ProcessorMap.Delete(t)
	obsMutex.Unlock()
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
