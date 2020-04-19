package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	dbus "github.com/godbus/dbus/v5"
	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	client "github.com/influxdata/influxdb1-client/v2"
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

type InfluxDBWrapper struct {
	c client.Client
}

func NewInfluxDBWrapper(addr string) (*InfluxDBWrapper, error) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: fmt.Sprintf("http://%s", addr),
	})
	if err != nil {
		return nil, err
	}

	return &InfluxDBWrapper{
		c: c,
	}, nil
}

func (w *InfluxDBWrapper) PushTemperatures(temps []int16, t time.Time) error {
	bps, err := client.NewBatchPoints(client.BatchPointsConfig{
		Precision:       "ms",
		Database:        "bbq",
		RetentionPolicy: "",
	})
	if err != nil {
		return err
	}

	tags := map[string]string{
		"cut": "brisket",
	}

	fields := map[string]interface{}{
		"probe1": temps[0],
		"probe2": temps[1],
		"probe3": temps[2],
		"probe4": temps[3],
		"probe5": temps[4],
		"probe6": temps[5],
	}

	p, err := client.NewPoint("temperature",
		tags,
		fields,
		t)
	if err != nil {
		return err
	}

	bps.AddPoint(p)

	return w.c.Write(bps)
}

func main() {
	conn, err := dbus.SystemBus()
	if err != nil {
		log.Fatal("SystemBus() failed, ", err)
	}
	defer conn.Close()

	sigch := make(chan *dbus.Signal, 128)
	conn.Signal(sigch)

	manager := NewObjectManager(conn, "/")

	devices, err := findDevices(manager, "BBQ")
	if err != nil {
		log.Fatal("findDevice() failed, ", err)
	}
	if len(devices) == 0 {
		log.Fatal("Do devices detected", err)
	}

	device := devices[0]

	db, err := NewInfluxDBWrapper("localhost:8086")
	if err != nil {
		log.Fatal("NewInfluxDBWrapper() failed, ", err)
	}

	b, err := NewBbq(device)
	if err != nil {
		log.Fatal("NewBbq() failed, ", err)
	}

	matchers := b.SignalMatchers()
	for _, m := range matchers {
		conn.AddMatchSignal(m.MatchOptions()...)
	}

	w := NewWeb()

	for {
		select {
		case s := <-sigch:
			for _, m := range matchers {
				m.Match(s)
			}

		case m := <-b.Measurements():
			if err := db.PushTemperatures(m.Temperatures, m.T); err != nil {
				log.Print("Failed to push temperatures, ", err)
			}

			w.PushMeasurement(m)

		}
	}
}
