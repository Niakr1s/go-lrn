package event_test

import (
	"lrn/event"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testEvent struct{}

type countHandler struct {
	counter int
}

func (ch *countHandler) handle(h *event.Handler, evt event.Event) {
	ch.counter++
}

func Test_Event(t *testing.T) {
	e := event.Emitter{}

	ch := countHandler{}
	// simple handler
	h := e.Add(
		ch.handle,
	)
	e.Emit(struct{}{})
	assert.Equal(t, 1, ch.counter)
	e.Emit(testEvent{})
	assert.Equal(t, 2, ch.counter)

	assert.Equal(t, 1, e.Len())
	h.Done()
	assert.Equal(t, 0, e.Len())

	e.Emit(testEvent{})
	assert.Equal(t, 2, ch.counter)
}

func Test_EventOnce(t *testing.T) {
	e := event.Emitter{}

	ch := countHandler{}
	// once
	h := e.Add(event.Combine(
		ch.handle,
		event.Once(),
	))

	assert.Equal(t, 1, e.Len())
	e.Emit(struct{}{})
	assert.Equal(t, 0, e.Len())
	assert.Equal(t, 1, ch.counter)

	e.Emit(testEvent{})
	assert.Equal(t, 1, ch.counter)

	assert.Equal(t, 0, e.Len())
	h.Done()
	assert.Equal(t, 0, e.Len())

	e.Emit(testEvent{})
	assert.Equal(t, 1, ch.counter)
}

func Test_EventOnly(t *testing.T) {
	e := event.Emitter{}

	ch := countHandler{}
	// only TestEvent
	h := e.Add(event.Combine(
		ch.handle,
		event.Only(func(e event.Event) bool {
			_, ok := e.(testEvent)
			return ok
		}),
	))

	assert.Equal(t, 1, e.Len())
	e.Emit(struct{}{})
	assert.Equal(t, 1, e.Len())
	assert.Equal(t, 0, ch.counter)

	e.Emit(testEvent{})
	assert.Equal(t, 1, ch.counter)

	assert.Equal(t, 1, e.Len())
	h.Done()
	assert.Equal(t, 0, e.Len())

	e.Emit(testEvent{})
	assert.Equal(t, 1, ch.counter)
}

func Test_EventCombine(t *testing.T) {
	e := event.Emitter{}

	ch := countHandler{}
	// only once after getting TestEvent struct
	h := e.Add(event.Combine(
		ch.handle,
		event.Only(func(e event.Event) bool {
			_, ok := e.(testEvent)
			return ok
		}),
		event.Once(),
	))

	assert.Equal(t, 1, e.Len())
	e.Emit(struct{}{})
	assert.Equal(t, 1, e.Len())
	assert.Equal(t, 0, ch.counter)

	e.Emit(testEvent{})
	assert.Equal(t, 1, ch.counter)

	assert.Equal(t, 0, e.Len())
	h.Done()
	assert.Equal(t, 0, e.Len())

	e.Emit(testEvent{})
	assert.Equal(t, 1, ch.counter)
}
