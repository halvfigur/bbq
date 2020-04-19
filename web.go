package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type (
	Web struct {
		upgrader websocket.Upgrader
		events   chan Measurement
		server   *http.Server

		mut       sync.RWMutex
		observers []observer
	}

	observer chan<- Measurement

	subject interface {
		addObserver(observer)

		removeObserver(observer)
	}
)

func NewWeb() *Web {
	mux := http.NewServeMux()

	w := &Web{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		mut:       sync.RWMutex{},
		observers: make([]observer, 0),
		server: &http.Server{
			Addr:    ":9000",
			Handler: mux,
		},
	}

	mux.HandleFunc("/", w.handler)

	go w.server.ListenAndServe()

	return w
}

func (w *Web) PushMeasurement(m Measurement) {
	w.notifyAll(m)
}

func (w *Web) Close() {
	w.server.Close()
}

func (w *Web) notifyAll(m Measurement) {
	w.mut.RLock()
	defer w.mut.RUnlock()

	for _, o := range w.observers {
		o <- m
	}
}

func (w *Web) addObserver(o observer) {
	w.mut.Lock()
	defer w.mut.Unlock()

	w.observers = append(w.observers, o)
}

func (w *Web) removeObserver(o observer) {
	w.mut.Lock()
	defer w.mut.Unlock()

	for i, x := range w.observers {
		if x == o {
			w.observers = append(w.observers[:i], w.observers[i+1:]...)
			return
		}
	}
}

func (web *Web) handler(w http.ResponseWriter, r *http.Request) {
	log.Print("Client connected")
	c, err := web.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Upgrade() failed, :", err)
		return
	}
	defer c.Close()

	ch := make(chan Measurement)
	defer close(ch)

	o := make(chan Measurement, 1)
	web.addObserver(o)
	defer web.removeObserver(o)

	for m := range o {
		if err := c.WriteJSON(m); err != nil {
			log.Print("WriteJSON() failed, :", err)
			return
		}
	}
}
