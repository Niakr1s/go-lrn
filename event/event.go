package event

type Event interface{}

type HandlerFunc func(h *Handler, evt Event)

type Handler struct {
	isDone bool
	Hf     HandlerFunc
}

func (h *Handler) Done() {
	h.isDone = true
}

type handlers []*Handler

type HandlerId int

type Emitter struct {
	handlers handlers
}

func (e Emitter) Len() int {
	l := 0
	for _, h := range e.handlers {
		if !h.isDone {
			l++
		}
	}
	return l
}

func (e *Emitter) Emit(evt Event) {
	for _, h := range e.handlers {
		if !h.isDone {
			h.Hf(h, evt)
		}
	}
	e.removeDoneHandlers()
}

func (e *Emitter) removeDoneHandlers() {
	activeHandlers := handlers{}
	for _, h := range e.handlers {
		if !h.isDone {
			activeHandlers = append(activeHandlers, h)
		}
	}
	e.handlers = activeHandlers
}

func (e *Emitter) Add(hf HandlerFunc) *Handler {
	h := &Handler{Hf: hf}
	e.handlers = append(e.handlers, h)
	return h
}

type Operator func(hf HandlerFunc) HandlerFunc

func Combine(hf HandlerFunc, ops ...Operator) HandlerFunc {
	if len(ops) == 0 {
		return hf
	}

	res := hf
	for i := len(ops) - 1; i >= 0; i-- {
		res = ops[i](res)
	}
	return res
}

func Only(f func(Event) bool) Operator {
	return func(hf HandlerFunc) HandlerFunc {
		return func(h *Handler, evt Event) {
			if f(evt) {
				hf(h, evt)
			}
		}
	}
}

func Once() Operator {
	return func(hf HandlerFunc) HandlerFunc {
		return func(h *Handler, evt Event) {
			hf(h, evt)
			h.Done()
		}
	}
}
