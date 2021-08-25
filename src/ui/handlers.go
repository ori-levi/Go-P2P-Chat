package ui

import "sync"

type Handler func(interface{})

type Handlers struct {
	lock     sync.RWMutex
	channel  chan interface{}
	handlers []Handler
}

func NewHandlers() *Handlers {
	return &Handlers{
		channel: make(chan interface{}),
	}
}

func (h *Handlers) Start() {
	go h.handleChannel()
}

func (h *Handlers) AddHandler(f Handler) {
	h.lock.Lock()
	defer h.lock.Unlock()

	h.handlers = append(h.handlers, f)
}

func (h *Handlers) handleChannel() {
	for {
		h.startHandlers(<-h.channel)
	}
}

func (h *Handlers) startHandlers(msg interface{}) {
	h.lock.RLock()
	defer h.lock.RUnlock()
	for _, handler := range h.handlers {
		handler(msg)
	}
}
