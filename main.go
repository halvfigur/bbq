package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/godbus/dbus"
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
	s, err := d.Service("0000fff0-0000-1000-8000-00805f9b34fb")
	if err != nil {
		return err
	}

	fff1, err := s.Characteristic("0000fff1-0000-1000-8000-00805f9b34fb")
	if err != nil {
		return err
	}
	if err = fff1.StartNotify(); err != nil {
		return err
	}

	fff3, err := s.Characteristic("0000fff3-0000-1000-8000-00805f9b34fb")
	if err != nil {
		return err
	}
	if err = fff3.StartNotify(); err != nil {
		return err
	}

	fff5, err := s.Characteristic("0000fff5-0000-1000-8000-00805f9b34fb")
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

	for s := range sigch {
		log.Printf("Signal, sender=%s, path=%v, name=%s, body=%v",
			s.Sender, s.Path, s.Name, s.Body)
	}
}
