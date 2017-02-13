// Copyright 2017 Jeff Nickoloff "jeff@allingeek.com"
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package clients

import (
	"context"
	"errors"
	"time"
)

type RetryFunc func() (interface{}, ClientError)
type CancellableFunc func(ctx context.Context) (interface{}, ClientError)

func Retry(f RetryFunc, pw PerishableWaiter) (interface{}, error) {
	pw.Start()
	for {
		result, err := f()
		if err == nil {
			return result, nil
		} else if !err.IsRetriable() {
			return nil, err.Error()
		}

		if e := pw.WaitOrDie(err.Error()); e != nil {
			return nil, err.Error()
		}
	}
}

func RetryPeriodic(f RetryFunc, timeout time.Duration, interval time.Duration, maxJitter time.Duration) (interface{}, error) {
	return Retry(f, &JitteredBackoff{
		TTL:       timeout,
		Initial:   interval,
		MaxJitter: maxJitter,
		Bof:       ConstantBackoff,
		Jf:        Jitter,
	})
}

func RetryLinear(f RetryFunc, timeout time.Duration, interval time.Duration, maxJitter time.Duration) (interface{}, error) {
	return Retry(f, &JitteredBackoff{
		TTL:       timeout,
		Initial:   interval,
		MaxJitter: maxJitter,
		Bof:       LinearBackoff,
		Jf:        Jitter,
	})
}

func RetryExponential(f RetryFunc, timeout time.Duration, initial time.Duration, maxJitter time.Duration) (interface{}, error) {
	return Retry(f, &JitteredBackoff{
		TTL:       timeout,
		Initial:   initial,
		MaxJitter: maxJitter,
		Bof:       ExponentialBackoff,
		Jf:        Jitter,
	})
}

type Waiter interface {
	WaitOrDie(e error) error
}
type Perishable interface {
	Start()
	IsDying() bool
}
type PerishableWaiter interface {
	Perishable
	Waiter
}

type JitteredBackoff struct {
	dead      <-chan time.Time
	round     uint
	TTL       time.Duration
	Initial   time.Duration
	MaxJitter time.Duration
	Bof       BackoffFunc
	Jf        JitterFunc
}

func (w *JitteredBackoff) WaitOrDie(e error) error {
	if w.Bof == nil {
		panic(errors.New(`Bof is nil`))
	}
	if w.Jf == nil {
		panic(errors.New(`Jf is nil`))
	}
	var ej, ei time.Duration
	// treat negative maxJitter like zero jitter
	if w.MaxJitter >= 0 && int(w.MaxJitter) > 0 {
		ej = w.MaxJitter
	}
	// treat negative interval like zero interval
	if w.Initial >= 0 {
		ei = w.Initial
	}

	select {
	case <-w.dead:
		return e
	case <-time.After(w.Bof(w.round, ei) + w.Jf(ej)):
	}
	w.round++
	return nil
}
func (w *JitteredBackoff) Start() {
	w.dead = time.After(w.TTL)
}
func (w *JitteredBackoff) IsDying() bool {
	return w.dead != nil
}
