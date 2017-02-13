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
	"math/rand"
	"time"
)

type BackoffFunc func(r uint, t time.Duration) time.Duration
type JitterFunc func(t time.Duration) time.Duration

func NoBackoff(round uint, zero time.Duration) time.Duration {
	return time.Duration(0)
}
func NoJitter(none time.Duration) time.Duration {
	return time.Duration(0)
}

func ConstantBackoff(round uint, flat time.Duration) time.Duration {
	return flat
}

func LinearBackoff(round uint, initial time.Duration) time.Duration {
	if initial <= 0 {
		return time.Duration(0)
	}
	return time.Duration(round+1) * initial
}

func ExponentialBackoff(round uint, initial time.Duration) time.Duration {
	if initial <= 0 {
		return time.Duration(0)
	}
	r := initial << round
	if r < 0 {
		return initial
	}
	return time.Duration(r)
}

func Jitter(max time.Duration) time.Duration {
	if max <= 0 {
		return time.Duration(0)
	}
	return time.Duration(rand.Intn(int(max/time.Millisecond))) * time.Millisecond
}
