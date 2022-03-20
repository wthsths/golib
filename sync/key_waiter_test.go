package gl_sync

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAcquireRelease(t *testing.T) {

	waiterAll := NewKeyWaiter()

	tests := []struct {
		waiter *KeyWaiter
		key    string
		count  int
	}{
		{waiter: waiterAll, key: "key-small", count: 10},
		{waiter: waiterAll, key: "key-medium", count: 100},
		{waiter: waiterAll, key: "key-large", count: 1000},
		{waiter: waiterAll, key: "key-large-2", count: 1000},
		{waiter: waiterAll, key: "key-huge", count: 5000},
	}
	for _, test := range tests {
		test := test
		t.Run(test.key, func(t *testing.T) {
			t.Parallel()
			waiter := test.waiter

			// Acquire release a non existing key
			assert.Panics(t, func() {
				waiter.Release(test.key)
			})

			// Acquire a key
			waiter.Acquire(test.key)

			// Release it once
			assert.NotPanics(t, func() {
				waiter.Release(test.key)
			})

			// Release a non acquired key
			assert.Panics(t, func() {
				waiter.Release(test.key)
			})

			// Acquire and release same key multiple times
			wg := sync.WaitGroup{}
			for i := 0; i < test.count; i++ {
				wg.Add(1)
				go func(index int, key string, wgIn *sync.WaitGroup, waiterIn *KeyWaiter) {
					defer wgIn.Done()
					// Acquire a key
					// t.Logf("acquiring %dth  %s", index, key)
					waiterIn.Acquire(key)

					// Release it
					// t.Logf("releasing %dth  %s", index, key)
					assert.NotPanics(t, func() {
						waiter.Release(test.key)
					})
				}(i, test.key, &wg, waiter)
			}

			wg.Wait()
		})

		for _, test := range tests {
			assert.NotContains(t, waiterAll.waitKeys, test.key)
		}
		assert.Empty(t, waiterAll.waitKeys)
	}
}
