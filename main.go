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

	if err := adapter.SetProperty(ctx, "Powered", true); err != nil {
		log.Fatal("SetPropert() failed, ", err)
	}

	filter := map[string]interface{}{
		"Transport":     "le",
		"DuplicateData": true,
	}
	if err := adapter.SetDiscoveryFilter(ctx, filter); err != nil {
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
}
