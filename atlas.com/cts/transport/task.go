package transport

import (
	"atlas-cts/configuration"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"time"
)

const StateEvaluationTask = "state_evaluation_task"

type stateEvaluationTask struct {
	l        logrus.FieldLogger
	c        *configuration.Configuration
	interval int
}

func NewStateEvaluationTask(l logrus.FieldLogger, c *configuration.Configuration, interval int) *stateEvaluationTask {
	return &stateEvaluationTask{l, c, interval}
}

func (r *stateEvaluationTask) Run() {
	span := opentracing.StartSpan(StateEvaluationTask)
	for _, tc := range r.c.Transports {
		t, err := GetRegistry().Get(tc.Source, tc.Destination)
		if err != nil {
			r.l.WithError(err).Errorf("Unable to retrieve transport from %d to %d.", tc.Source, tc.Destination)
			continue
		}
		if !t.Enabled() {
			continue
		}
		ns := getState(time.Now(), tc)
		if t.State() != ns {
			err = UpdateState(r.l, span)(tc.Source, tc.Destination, ns)
			if err != nil {
				r.l.WithError(err).Errorf("Unable to transistion state of transport.")
				continue
			}
			r.l.Debugf("Transport from %d to %d will transfer from %s state to %s state.", tc.Source, tc.Destination, t.State(), ns)
		}
	}
	span.Finish()
}

func (r *stateEvaluationTask) SleepTime() time.Duration {
	return time.Millisecond * time.Duration(r.interval)
}
