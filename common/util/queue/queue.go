package queue

type (
	Queue struct {
		start, end *node
		length int
	}
	node struct {
		value interface{}
		next *node
	}
)

// Create a new queue
func New() *Queue {
	return &Queue{nil,nil,0}
}
// Take the next item off the front of the queue
func (q *Queue) Dequeue() interface{} {
	if q.length == 0 {
		return nil
	}
	n := q.start
	if q.length == 1 {
		q.start = nil
		q.end = nil
	} else {
		q.start = q.start.next
	}
	q.length--
	return n.value
}
// Put an item on the end of a queue
func (q *Queue) Enqueue(value interface{}) {
	n := &node{value,nil}
	if q.length == 0 {
		q.start = n
		q.end = n
	} else {
		q.end.next = n
		q.end = n
	}
	q.length++
}
// Return the number of items in the queue
func (q *Queue) Len() int {
	return q.length
}

// Return the first item in the queue without removing it
func (q *Queue) Peek() interface{} {
	if q.length == 0 {
		return nil
	}
	return q.start.value
}

func (q *Queue) ForEach(f func(e interface{}) bool ) {

	n := q.start

	for i:=0;i<q.length;i++ {

		if n !=nil {

			if br := f(n.value);br {

				break
			}
		}

		n = n.next
	}
}
