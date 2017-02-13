package siren

//    Copyright 2017 Jeff Nickoloff "jeff@allingeek.com"
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
//

import (
	"context"
	"fmt"
	"time"
)

var cancelfuncs map[string]context.CancelFunc
var backends map[string]Backend

func Activate(ctx context.Context, monitor Monitor) {
	if _, ok := backends[monitor.Trigger.Backend]; !ok {
		panic(fmt.Errorf("No such registered backend: %v", monitor.Trigger.Backend))
	}

	if cancelfuncs == nil {
		cancelfuncs = map[string]context.CancelFunc{}
	}

	Deactivate(monitor)

	ctx, cf := context.WithCancel(ctx)
	cancelfuncs[monitor.ID] = cf
	go daemon(ctx, monitor)
}

func Deactivate(monitor Monitor) {
	if cancelfuncs == nil {
		return
	}
	if c, ok := cancelfuncs[monitor.ID]; ok {
		c()
		delete(cancelfuncs, monitor.ID)
	}
}

func daemon(ctx context.Context, monitor Monitor) {
	t := monitor.Trigger
	b := backends[t.Backend]

	consecutiveErrors := 0
lup:
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(monitor.Period):
			data, err := b(ctx, t.Query)
			if err != nil {
				consecutiveErrors++
				if consecutiveErrors > 3 {
					panic(fmt.Errorf("Unable to fetch data: %v", err))
				}
				continue
			}
			consecutiveErrors = 0

			if len(data) < monitor.Alarm.Count {
				continue
			}

			trend := 0
			for i := 0; i < len(data); i++ {
				dp := data[i]
				if monitor.Alarm.Comparator()(dp, monitor.Alarm.Threshold) {
					if trend < 0 {
						trend--
					} else {
						trend = -1
					}
				} else {
					if trend > 0 {
						trend++
					} else {
						trend = 1
					}
				}

				if trend == monitor.Alarm.ClearCount {
					// is STATE_CLEAR
					for _, cb := range monitor.ClearActions {
						go cb(ctx, monitor)
					}
					continue lup
				} else if trend == -1*monitor.Alarm.Count {
					// is STATE_ALARM
					for _, cb := range monitor.AlarmActions {
						go cb(ctx, monitor)
					}
					continue lup
				}
			}
			// No trend detected in window - is STATE_FLAPPING
			for _, cb := range monitor.FlappingActions {
				go cb(ctx, monitor)
			}
		}
	}
}

func RegisterBackend(name string, f Backend) {
	if backends == nil {
		backends = map[string]Backend{}
	}
	if len(name) <= 0 {
		panic(fmt.Errorf("Empty backend name provided"))
	}
	if f == nil {
		panic(fmt.Errorf("No backend function provided"))
	}
	if _, ok := backends[name]; ok {
		panic(fmt.Errorf("Duplicate backend name used: %v", name))
	}
	backends[name] = f
}

func ClearAll() {
	for bn := range backends {
		delete(backends, bn)
	}
	for mi, cf := range cancelfuncs {
		if cf != nil {
			cf()
		}
		delete(cancelfuncs, mi)
	}
}

func (c Condition) Comparator() Comparator {
	switch c.ComparatorID {
	case COMP_LESSER: return lesser
	case COMP_LESSEREQUAL: return lesserEqual
	case COMP_GREATER: return greater
	case COMP_GREATEREQUAL: return greaterEqual
	default:
		panic(fmt.Errorf(`No such comparator.`))
	}
}
func lesser(p float64, threshold float64) bool {
	return p < threshold
}
func lesserEqual(p float64, threshold float64) bool {
	return p <= threshold
}
func greater(p float64, threshold float64) bool {
	return p > threshold
}
func greaterEqual(p float64, threshold float64) bool {
	return p <= threshold
}
