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
	"sync/atomic"
	"testing"
	"time"
)

func TestDebouncer_Trigger(t *testing.T) {
	var callCount atomic.Int32

	debouncer := NewDebouncer(50*time.Millisecond, func() {
		callCount.Add(1)
	})

	// Trigger multiple times rapidly
	debouncer.Trigger()
	debouncer.Trigger()
	debouncer.Trigger()

	// Wait for debounce delay
	time.Sleep(100 * time.Millisecond)

	// Should only be called once
	if count := callCount.Load(); count != 1 {
		t.Errorf("Debouncer called %d times, want 1", count)
	}
}

func TestDebouncer_Stop(t *testing.T) {
	var callCount atomic.Int32

	debouncer := NewDebouncer(50*time.Millisecond, func() {
		callCount.Add(1)
	})

	debouncer.Trigger()
	debouncer.Stop()

	// Wait past debounce delay
	time.Sleep(100 * time.Millisecond)

	// Should not be called since we stopped it
	if count := callCount.Load(); count != 0 {
		t.Errorf("Debouncer called %d times after stop, want 0", count)
	}
}

func TestDebouncer_Flush(t *testing.T) {
	var callCount atomic.Int32

	debouncer := NewDebouncer(50*time.Millisecond, func() {
		callCount.Add(1)
	})

	debouncer.Trigger()
	debouncer.Flush()

	// Should be called immediately
	if count := callCount.Load(); count != 1 {
		t.Errorf("Debouncer called %d times, want 1", count)
	}

	// Wait to ensure no second call
	time.Sleep(100 * time.Millisecond)

	if count := callCount.Load(); count != 1 {
		t.Errorf("Debouncer called %d times after flush, want 1", count)
	}
}

func TestDebouncer_RapidTriggers(t *testing.T) {
	var callCount atomic.Int32

	debouncer := NewDebouncer(100*time.Millisecond, func() {
		callCount.Add(1)
	})

	// Simulate rapid file changes
	for i := 0; i < 10; i++ {
		debouncer.Trigger()
		time.Sleep(20 * time.Millisecond) // Faster than debounce delay
	}

	// Wait for final debounce
	time.Sleep(150 * time.Millisecond)

	// Should only be called once despite 10 triggers
	if count := callCount.Load(); count != 1 {
		t.Errorf("Debouncer called %d times, want 1", count)
	}
}

func TestDebouncer_MultipleCycles(t *testing.T) {
	var callCount atomic.Int32

	debouncer := NewDebouncer(50*time.Millisecond, func() {
		callCount.Add(1)
	})

	// First cycle
	debouncer.Trigger()
	time.Sleep(100 * time.Millisecond)

	// Second cycle
	debouncer.Trigger()
	time.Sleep(100 * time.Millisecond)

	// Should be called twice (once per cycle)
	if count := callCount.Load(); count != 2 {
		t.Errorf("Debouncer called %d times, want 2", count)
	}
}

func TestDebouncer_ConcurrentTriggers(t *testing.T) {
	var callCount atomic.Int32

	debouncer := NewDebouncer(50*time.Millisecond, func() {
		callCount.Add(1)
	})

	// Trigger from multiple goroutines
	done := make(chan bool)
	for i := 0; i < 5; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				debouncer.Trigger()
				time.Sleep(5 * time.Millisecond)
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 5; i++ {
		<-done
	}

	// Wait for final debounce
	time.Sleep(100 * time.Millisecond)

	// Should only be called once despite concurrent triggers
	if count := callCount.Load(); count != 1 {
		t.Errorf("Debouncer called %d times, want 1", count)
	}
}
