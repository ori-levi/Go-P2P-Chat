package ui

import (
	"sync"
	"time"
)

type ConsumerHandler func(interface{})
type ProviderHandler func() interface{}

type Handlers struct {
	channel       chan interface{}
	consumersLock sync.RWMutex
	consumers     []ConsumerHandler
	providersLock sync.RWMutex
	providers     []ProviderHandler
}

func NewHandlers() *Handlers {
	return &Handlers{
		channel: make(chan interface{}),
	}
}

func (h *Handlers) Start() {
	go h.provideChannel()
	go h.consumeChannel()
}

func (h *Handlers) AddConsumer(f ConsumerHandler) {
	h.consumersLock.Lock()
	defer h.consumersLock.Unlock()

	h.consumers = append(h.consumers, f)
}

func (h *Handlers) AddProvider(f ProviderHandler) {
	h.providersLock.Lock()
	defer h.providersLock.Unlock()

	h.providers = append(h.providers, f)
}

func (h *Handlers) consumeChannel() {
	for {
		h.startConsumers(<-h.channel)
	}
}

func (h *Handlers) startConsumers(msg interface{}) {
	h.consumersLock.RLock()
	defer h.consumersLock.RUnlock()
	for _, consumer := range h.consumers {
		consumer(msg)
	}
}

func (h *Handlers) provideChannel() {
	for {
		messages := h.getMessages()
		if len(messages) > 0 {
			for _, msg := range messages {
				h.channel <- msg
			}

		} else {
			time.Sleep(time.Second / 2)
		}
	}
}

func (h *Handlers) getMessages() []interface{} {
	h.consumersLock.RLock()
	defer h.consumersLock.RUnlock()

	messages := make([]interface{}, len(h.consumers))
	for _, provider := range h.providers {
		messages = append(messages, provider())
	}

	return messages
}
