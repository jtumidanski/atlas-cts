package topic

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"sync"
)

type registry struct {
	topics map[string]string
	lock   sync.RWMutex
}

var once sync.Once
var r *registry

func GetRegistry() *registry {
	once.Do(func() {
		r = &registry{
			topics: make(map[string]string, 0),
			lock:   sync.RWMutex{},
		}
	})
	return r
}

func (r *registry) Get(l logrus.FieldLogger, span opentracing.Span, token string) string {
	r.lock.RLock()
	if val, ok := r.topics[token]; ok {
		r.lock.RUnlock()
		return val
	} else {
		r.lock.RUnlock()
		r.lock.Lock()
		if val, ok = r.topics[token]; ok {
			r.lock.Unlock()
			return val
		}
		td, err := getTopic(token)(l, span)
		if err != nil {
			r.lock.Unlock()
			l.WithError(err).Fatalf("Unable to locate topic for token %s.", token)
			return ""
		}
		attr := td.Data().Attributes

		r.topics[token] = attr.Name
		r.lock.Unlock()
		return attr.Name
	}
}
