package queue

type Queue interface {
	Next() interface{}
	IsEmpty() bool
	Add(interface{})
	Length() int
}

type SimpleQueue struct {
	data []interface{}
}

func (s *SimpleQueue) Next() interface{} {
	item := s.data[0]
	s.data = s.data[1:]

	return item
}

func (s *SimpleQueue) Length() int {
	return len(s.data)
}

func (s *SimpleQueue) IsEmpty() bool {
	return len(s.data) == 0
}

func (s *SimpleQueue) Add(item interface{}) {
	s.data = append(s.data, item)
}

func New() Queue {
	return &SimpleQueue{
		data: make([]interface{}, 0),
	}
}
