package main

import (
	"log"

	dbus "github.com/godbus/dbus/v5"
)

type (
	SignalHandler func(s *dbus.Signal)

	SignalMatcher struct {
		options []dbus.MatchOption
		ch      chan *dbus.Signal
	}
)

func NewSignalMatcher(handler SignalHandler, options ...dbus.MatchOption) *SignalMatcher {

	log.Printf("NewSignalMatcher()")
	ch := make(chan *dbus.Signal, 16)

	go func() {
		for body := range ch {
			handler(body)
		}
	}()

	return &SignalMatcher{
		options: options,
		ch:      ch,
	}
}

func (m *SignalMatcher) Close() {
	close(m.ch)
}

func (m *SignalMatcher) MatchOptions() []dbus.MatchOption {
	return m.options
}

func (m *SignalMatcher) Match(s *dbus.Signal) {
	m.ch <- s
}
