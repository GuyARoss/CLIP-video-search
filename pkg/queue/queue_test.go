package queue

import "testing"

func TestQueue(t *testing.T) {
	q := New()
	messages := []string{"is this working?", "is this working x2?"}

	for _, m := range messages {
		q.Add(m)
	}

	item := q.Next()
	if item != messages[0] {
		t.Errorf("expected '%s' got '%s'", messages[0], item)
	}

	item = q.Next()
	if item != messages[1] {
		t.Errorf("expected '%s' got '%s'", messages[1], item)
	}
}
