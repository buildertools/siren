package siren

//   Copyright 2017 Jeff Nickoloff "jeff@allingeek.com"
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

import (
	"context"
	"time"
)

const (
	STATE_CLEAR = iota
	STATE_ALARM
	STATE_FLAPPING
)

const (
	SHAPE_DURATION = iota
	SHAPE_COUNTER
	SHAPE_RATE
)

const (
	REL_ABOVE_OR_EQUAL = iota
	REL_ABOVE
	REL_BELOW
)

const (
	COMP_LESSER = iota
	COMP_LESSEREQUAL
	COMP_GREATER
	COMP_GREATEREQUAL
)

type Monitor struct {
	ID              string
	Metadata        map[string]interface{}
	Trigger         Metric
	Alarm           Condition
	AlarmActions    []Action          `json:"-"`
	ClearActions    []Action          `json:"-"`
	FlappingActions []Action          `json:"-"`
	Suppressed      bool
	Period          time.Duration
}

type Action func(ctx context.Context, monitor Monitor)

type Notice struct {
	MonitorID string
	Metadata  map[string]interface{}
	State     int
}

type Metric struct {
	Shape   int
	Backend string
	Query   string
}

type Condition struct {
	Threshold    float64
	Count        uint
	ClearCount   uint
	ComparatorID int
}
type Comparator func(p float64, threshold float64) bool


type Stateful interface {
	State() int
}

type Suppressable interface {
	Suppress() error
}

type Backend func(context context.Context, query string) ([]float64, error)
