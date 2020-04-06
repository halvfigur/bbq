package main

import (
	"context"
	"log"
	"reflect"

	"github.com/godbus/dbus"
)

type (
	ObjectManager struct {
		DBusObjectProxy
		iface string
	}
)

func NewObjectManager(conn *dbus.Conn, path string) *ObjectManager {
	debug("NewObjectManager(%v, %v)", conn, path)

	return &ObjectManager{
		DBusObjectProxy: newDBusObjectProxy(conn, destOrgBluez, path),
		iface:           "org.freedesktop.DBus.ObjectManager",
	}
}

func (m *ObjectManager) GetManagedObjects(ctx context.Context) map[string]map[string]map[string]interface{} {

	v, _ := m.call(ctx, "org.freedesktop.DBus.ObjectManager.GetManagedObjects", 0)

	//v[0] is map[dbus.ObjectPath]map[string]map[string]dbus.Variant meaning
	// object path -> interface -> property -> variant

	objs := v[0].(map[dbus.ObjectPath]map[string]map[string]dbus.Variant)

	paths := make(map[string]map[string]map[string]interface{})

	for opath, oifaces := range objs {
		ifaces := make(map[string]map[string]interface{})

		for oiface, oattrs := range oifaces {
			attrs := make(map[string]interface{})

			for attr, value := range oattrs {
				attrs[attr] = value.Value()
			}

			ifaces[oiface] = attrs
		}

		paths[string(opath)] = ifaces
	}

	log.Println("paths -> ", reflect.TypeOf(paths))
	return paths
}
