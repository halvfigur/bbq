package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

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

func main() {
	conn, err := dbus.ConnectSystemBus()
	if err != nil {
		log.Fatal("ConnectSystemBus() failed, ", err)
	}
	defer conn.Close()

	ctx := context.Background()

	manager := NewObjectManager(conn, "/")
	/*
		objs := manager.GetManagedObjects(ctx)

		printObjs(objs)
	*/

	devices := findDevice(ctx, manager, "BBQ")
	printObjs(devices)

	if err := devices[0].Connect(ctx); err != nil {
		log.Fatal("Connect() failed, ", err)
	}
	defer devices[0].Disconnect(ctx)

	log.Print("Sleeping for 5")
	time.Sleep(5 * time.Second)
}
