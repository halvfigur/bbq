package main

import (
	"context"
	"log"
	"math"
	"time"

	dbus "github.com/godbus/dbus/v5"
)

const (
	fff0UUID = "0000fff0-0000-1000-8000-00805f9b34fb"
	fff1UUID = "0000fff1-0000-1000-8000-00805f9b34fb"
	fff3UUID = "0000fff3-0000-1000-8000-00805f9b34fb"
	fff5UUID = "0000fff5-0000-1000-8000-00805f9b34fb"
)

type (
	BbqDB interface {
		PushTemperatures(temps []int16, t time.Time) error
	}

	Measurement struct {
		Temperatures []int16
		T            time.Time
	}

	Bbq struct {
		dev      *Device
		tempPath dbus.ObjectPath
		events   chan Measurement
		matchers []*SignalMatcher
	}
)

func NewBbq(dev *Device) (*Bbq, error) {
	if err := dev.Connect(context.Background()); err != nil {
		return nil, err
	}

	tempPath, err := dev.Characteristic(fff5UUID)
	if err != nil {
		return nil, err
	}

	b := &Bbq{
		dev:      dev,
		tempPath: tempPath.Path(),
		events:   make(chan Measurement, 1),
	}

	if err := b.setupSignalMatchers(); err != nil {
		return nil, err
	}

	if err := b.startNotifications(); err != nil {
		return nil, err
	}

	return b, nil
}

func (b *Bbq) setupSignalMatchers() error {
	b.matchers = make([]*SignalMatcher, 0)

	tempChar, err := b.dev.Characteristic(fff5UUID)
	if err != nil {
		return err
	}

	tempPath := tempChar.Path()

	// For a description of matching rules see
	// https://dbus.freedesktop.org/doc/dbus-specification.html#message-bus-routing-match-rules
	m := NewSignalMatcher(b.handleTemperatureUpdate,
		dbus.WithMatchObjectPath(tempPath),
		dbus.WithMatchInterface("org.freedesktop.DBus.Properties"),
		dbus.WithMatchMember("PropertiesChanged"),
	)

	b.matchers = append(b.matchers, m)

	return nil
}

func (b *Bbq) startNotifications() error {
	s, err := b.dev.Service(fff0UUID)
	if err != nil {
		return err
	}

	fff1, err := s.Characteristic(fff1UUID)
	if err != nil {
		return err
	}
	if err = fff1.StartNotify(); err != nil {
		return err
	}

	fff3, err := s.Characteristic(fff3UUID)
	if err != nil {
		return err
	}
	if err = fff3.StartNotify(); err != nil {
		return err
	}

	fff5, err := s.Characteristic(fff5UUID)
	if err != nil {
		return err
	}
	if err = fff5.StartNotify(); err != nil {
		return err
	}

	payloads := [][]byte{
		[]byte{0x23, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x23},
		[]byte{0x22, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x22},
		[]byte{0x22, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x22},
		[]byte{0x22, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x22},
		[]byte{0x22, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x22},
		[]byte{0x22, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x22},
		[]byte{0x22, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x22},
		[]byte{0x24, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24},
	}

	fff4, err := s.Characteristic("0000fff4-0000-1000-8000-00805f9b34fb")
	if err != nil {
		return err
	}

	for _, payload := range payloads {
		if err := fff4.WriteValue(payload, nil); err != nil {
			return err
		}
	}

	return nil
}

func (b *Bbq) handleTemperatureUpdate(s *dbus.Signal) {
	log.Printf("SignalMatcher.Match(signal=%v)", s)

	t := time.Now().UTC()

	if b.tempPath != s.Path {
		return
	}

	if s.Name != "org.freedesktop.DBus.Properties.PropertiesChanged" {
		return
	}

	// Properties changed should have a 3 element payload
	if len(s.Body) != 3 {
		return
	}

	// The first element is the name of the interface who's property
	// changed
	iface, ok := s.Body[0].(string)
	if !ok {
		return
	}

	if iface != "org.bluez.GattCharacteristic1" {
		return
	}

	// The second element is a dictionary mapping the names of the
	// properties that have changed to their new values
	props, ok := s.Body[1].(map[string]dbus.Variant)
	if !ok {
		return
	}

	// This is the temperature property
	raw, ok := props["Value"]
	if !ok {
		return
	}

	data, ok := raw.Value().([]uint8)
	if !ok {
		return
	}

	// Only push out the changes if it won't block us
	if len(b.events) == 0 {
		b.events <- b.newMeasurement(data, t)
	}
}

func (b *Bbq) newMeasurement(data []uint8, t time.Time) Measurement {
	temps := make([]int16, 6)
	for i := 0; i < 6; i++ {
		l, h := data[2*i], data[2*i+1]

		if h == 0 {
			temps[i] = int16(l)
			continue
		}

		temps[i] = math.MinInt16
	}

	return Measurement{temps, t}
}

func (b *Bbq) Measurements() chan Measurement {
	return b.events
}

func (b *Bbq) Close() error {
	close(b.events)

	return b.dev.Disconnect(context.Background())
}

func (b *Bbq) SignalMatchers() []*SignalMatcher {
	return b.matchers
}
