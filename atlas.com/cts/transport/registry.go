package transport

import (
	"errors"
	"fmt"
	"sync"
)

type Registry struct {
	mutex sync.Mutex

	transports     map[string]Model
	transportLocks map[string]*sync.RWMutex
}

var registry *Registry
var once sync.Once

func GetRegistry() *Registry {
	once.Do(func() {
		registry = &Registry{}
		registry.transports = make(map[string]Model)
		registry.transportLocks = make(map[string]*sync.RWMutex)
	})
	return registry
}

func getKey(sourceId uint32, destinationId uint32) string {
	return fmt.Sprintf("%d:%d", sourceId, destinationId)
}

func (r *Registry) getLock(sourceId uint32, destinationId uint32) *sync.RWMutex {
	mk := getKey(sourceId, destinationId)
	return r.getLockWithKey(mk)
}

func (r *Registry) getLockWithKey(mk string) *sync.RWMutex {
	if val, ok := r.transportLocks[mk]; ok {
		return val
	} else {
		var mm = &sync.RWMutex{}
		r.mutex.Lock()
		r.transportLocks[mk] = mm
		r.mutex.Unlock()
		return mm
	}
}

func (r *Registry) Add(model Model) {
	k := getKey(model.Source(), model.Destination())
	tl := r.getLockWithKey(k)
	tl.Lock()
	r.transports[k] = model
	tl.Unlock()
}

func (r *Registry) Get(sourceId uint32, destinationId uint32) (Model, error) {
	k := getKey(sourceId, destinationId)
	tl := r.getLockWithKey(k)
	tl.RLock()
	if val, ok := r.transports[k]; ok {
		tl.RUnlock()
		return val, nil
	}
	tl.RUnlock()
	return Model{}, errors.New("not found")
}

func (r *Registry) GetAll() []Model {
	var result []Model
	for _, t := range r.transports {
		result = append(result, t)
	}
	return result
}

func (r *Registry) Update(model Model) error {
	k := getKey(model.Source(), model.Destination())
	tl := r.getLockWithKey(k)
	tl.Lock()
	r.transports[k] = model
	tl.Unlock()
	return nil
}
