package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/godbus/dbus"
)

const (
	fff0UUID = "0000fff0-0000-1000-8000-00805f9b34fb"
	fff1UUID = "0000fff1-0000-1000-8000-00805f9b34fb"
	fff3UUID = "0000fff3-0000-1000-8000-00805f9b34fb"
	fff5UUID = "0000fff5-0000-1000-8000-00805f9b34fb"
)

func printObjs(objs interface{}) {
	blob, err := json.Marshal(objs)
	if err != nil {
		log.Fatal("json.Marshal() failed, ", err)
	}

	buf := new(bytes.Buffer)

	if err := json.Indent(buf, blob, "", "    "); err != nil {
		log.Fatal("json.Indent() failed, ", err)
	}

	fmt.Println(buf.String())
}

func scan(ctx context.Context, conn *dbus.Conn) {
	adapter := NewAdapter(conn, "/org/bluez/hci0")

	if err := adapter.SetPowered(true); err != nil {
		log.Fatal("SetPowered() failed, ", err)
	}

	filter := map[string]interface{}{
		"Transport":     "le",
		"DuplicateData": true,
	}
	if err := adapter.SetDiscoveryFilter(filter); err != nil {
		log.Fatal("SetDiscoveryFilter() failed, ", err)
	}

	d, err := adapter.StartDiscovery(ctx)
	if err != nil {
		log.Fatal("StartDiscovery() failed, ", err)
	}

	log.Print(d)
}

func startNotifications(d *Device) error {
	s, err := d.Service(fff0UUID)
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

func tempHandler(path dbus.ObjectPath) SignalHandler {
	return func(s *dbus.Signal) {
		log.Printf("SignalMatcher.Match(signal=%v)", s)

		if path != s.Path {
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

		probeSlice, ok := raw.Value().([]uint8)
		if !ok {
			return
		}

		log.Printf("tempHandler -  Body: %v", probeSlice)
	}
}

func setupSignalMatchers(d *Device) ([]*SignalMatcher, error) {
	matchers := make([]*SignalMatcher, 0)

	tempChar, err := d.Characteristic(fff5UUID)
	if err != nil {
		return nil, err
	}

	tempPath := tempChar.Path()

	// For a description of matching rules see
	// https://dbus.freedesktop.org/doc/dbus-specification.html#message-bus-routing-match-rules
	m := NewSignalMatcher(tempHandler(tempPath),
		dbus.WithMatchObjectPath(tempPath),
		dbus.WithMatchInterface("org.freedesktop.DBus.Properties"),
		dbus.WithMatchMember("PropertiesChanged"),
	)

	matchers = append(matchers, m)

	return matchers, nil
}

func main() {
	conn, err := dbus.ConnectSystemBus()
	if err != nil {
		log.Fatal("ConnectSystemBus() failed, ", err)
	}
	defer conn.Close()

	sigch := make(chan *dbus.Signal, 128)
	conn.Signal(sigch)

	ctx := context.Background()

	manager := NewObjectManager(conn, "/")

	devices, err := findDevices(manager, "BBQ")
	if err != nil {
		log.Fatal("findDevice() failed, ", err)
	}
	if len(devices) == 0 {
		log.Fatal("Do devices detected", err)
	}

	device := devices[0]

	if err := device.Connect(ctx); err != nil {
		log.Fatal("Connect() failed, ", err)
	}
	defer devices[0].Disconnect(ctx)

	if err = startNotifications(device); err != nil {
		log.Fatal("startNotifications() failed, ", err)
	}

	matchers, err := setupSignalMatchers(devices[0])
	if err != nil {
		log.Fatal("setupSignalMatchers() failed, ", err)
	}

	for _, m := range matchers {
		conn.AddMatchSignal(m.MatchOptions()...)
	}

	for s := range sigch {
		for _, m := range matchers {
			m.Match(s)
		}
	}
}
