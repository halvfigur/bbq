package main

import (
	"log"
	"reflect"

	"github.com/godbus/dbus"
)

type (
	ObjectManager struct {
		DBusObjectProxy
	}
)

func NewObjectManager(conn *dbus.Conn, path string) *ObjectManager {
	debug("NewObjectManager(%v, %v)", conn, path)

	return &ObjectManager{
		DBusObjectProxy: newDBusObjectProxy(conn, destOrgBluez, "org.freedesktop.DBus.ObjectManager", path),
	}
}

func (m *ObjectManager) GetManagedObjects() (map[string]map[string]map[string]interface{}, error) {

	//v[0] is map[dbus.ObjectPath]map[string]map[string]dbus.Variant meaning
	// object path -> interface -> property -> variant

	objs := make(map[dbus.ObjectPath]map[string]map[string]dbus.Variant)
	if err := m.Call("org.freedesktop.DBus.ObjectManager.GetManagedObjects", 0).Store(objs); err != nil {
		return nil, err
	}

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
	return paths, nil
}
