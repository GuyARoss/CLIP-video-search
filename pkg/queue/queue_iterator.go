package queue

import "time"

type QueueIterator interface {
	Next(interface{})
}

func LongLivedIterator(q Queue, iter QueueIterator) {
	for {
		if !q.IsEmpty() {
			iter.Next(q.Next())
		} else {
			time.Sleep(time.Millisecond * 200)
		}
	}
}
