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
	"testing"
)

func panicPass(t *testing.T) {
	if r := recover(); r == nil {
		t.Fatal(`Failed to panic`)
	}
}

func panicFail(t *testing.T) {
	if r := recover(); r != nil {
		t.Fatal(`Paniced`)
	}
}

func TestRegisterBackendValidation(t *testing.T) {
	b := func(c context.Context, q string) ([]float64, error) { return nil, nil }

	t.Run(`EmptyName`, func(t *testing.T) {
		defer panicPass(t)
		RegisterBackend(``, b)
	})
	t.Run(`NilBackend`, func(t *testing.T) {
		defer panicPass(t)
		RegisterBackend(`name1`, nil)
	})
	t.Run(`DoubleName`, func(t *testing.T) {
		defer panicPass(t)
		RegisterBackend(`name2`, b)
		RegisterBackend(`name2`, b)
	})
	t.Run(`Three unique`, func(t *testing.T) {
		defer panicFail(t)
		RegisterBackend(`name3`, b)
		RegisterBackend(`name4`, b)
		RegisterBackend(`name5`, b)
	})
	ClearAll()
}

func TestClearAll(t *testing.T) {
	b := func(c context.Context, q string) ([]float64, error) { return nil, nil }
	t.Run(`ReregisterBackend`, func(t *testing.T) {
		defer panicFail(t)
		RegisterBackend(`name2`, b)
		ClearAll()
		RegisterBackend(`name2`, b)
	})
}

func TestDeactivate(t *testing.T) {
	t.Run(`NoPanicOnMissing`, func(t *testing.T) {
		defer panicFail(t)
		m := Monitor{ID: `ray`}
		Deactivate(m)
	})
}

func TestActivate(t *testing.T) {
	t.Run(`NoSuchBackend`, func(t *testing.T) {
		defer panicPass(t)
		m := Monitor{ID: `egon`, Trigger: Metric{Backend: `notthere`}}
		Activate(context.Background(), m)
	})
}
