// Copyright 2025 Scott Friedman
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

package watch

import (
	"sync"
	"time"
)

// Debouncer groups rapid events and triggers callback after quiet period.
type Debouncer struct {
	delay    time.Duration
	callback func()
	timer    *time.Timer
	mu       sync.Mutex
}

// NewDebouncer creates a new debouncer with specified delay.
func NewDebouncer(delay time.Duration, callback func()) *Debouncer {
	return &Debouncer{
		delay:    delay,
		callback: callback,
	}
}

// Trigger resets the debounce timer. Callback fires after delay with no new triggers.
func (d *Debouncer) Trigger() {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Cancel existing timer if any
	if d.timer != nil {
		d.timer.Stop()
	}

	// Start new timer
	d.timer = time.AfterFunc(d.delay, d.callback)
}

// Stop cancels any pending callback.
func (d *Debouncer) Stop() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
		d.timer = nil
	}
}

// Flush immediately triggers callback and cancels timer.
func (d *Debouncer) Flush() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
		d.timer = nil
	}

	d.callback()
}
