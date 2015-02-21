// Copyright 2013-2014 Aerospike, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package atomic

import (
	"math"
	"sync"
)

// AtomicQueue is a blocking FIFO queue.
// If the queue is empty, nil is returned.
// if the queue is full, offer will return false
type AtomicQueue struct {
	head, tail uint32
	data       []interface{}
	size       uint32
	lock       sync.Mutex
}

// NewQueue creates a new queue with initial size.
func NewAtomicQueue(size int) *AtomicQueue {
	if size <= 0 {
		panic("Queue size cannot be less than 1")
	}

	return &AtomicQueue{
		data: make([]interface{}, uint32(size)),
		size: uint32(size),
	}
}

// Offer adds an item to the queue unless the queue is full.
// In case the queue is full, the item will not be added to the queue
// and false will be returned
func (q *AtomicQueue) Offer(obj interface{}) bool {
	q.lock.Lock()

	if (q.tail <= math.MaxInt32 && q.head <= math.MaxInt32 && q.head > q.tail+q.size) ||
		(q.tail > math.MaxInt32 && q.head > math.MaxInt32 && q.head-q.tail >= q.size) ||
		(q.tail > math.MaxInt32 && q.head <= math.MaxInt32 && q.tail >= q.head-q.size) {

		q.lock.Unlock()
		return false
	}

	q.head++
	q.data[q.head%q.size] = obj
	q.lock.Unlock()
	return true
}

// Poll removes and returns an item from the queue.
// If the queue is empty, nil will be returned.
func (q *AtomicQueue) Poll() (res interface{}) {
	q.lock.Lock()

	// queue is empty
	if q.tail != q.head {
		q.tail++
		res = q.data[q.tail%q.size]
	}

	q.lock.Unlock()
	return res
}
