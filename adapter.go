package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/godbus/dbus"
)

type (
	Adapter struct {
		dobj dbus.BusObject
	}
)

var ErrContext = errors.New("context done")

func call(ctx context.Context, dobj dbus.BusObject, method string, flags dbus.Flags, args ...interface{}) ([]interface{}, error) {
	debug("call(ctx, %v, %v, %v)", method, flags, args)

	ch := make(chan *dbus.Call, 1)

	dobj.Go(method, flags, ch, args...)
	select {
	case <-ctx.Done():
		return nil, ErrContext
	case c := <-ch:
		return c.Body, c.Err
	}
}

func NewAdapter(conn *dbus.Conn, iface string) *Adapter {
	debug("NewAdapter(%v, %v)", conn, iface)

	return &Adapter{
		dobj: conn.Object("org.bluez", dbus.ObjectPath(fmt.Sprintf("/org/bluez/%s", iface))),
	}
}

func (a *Adapter) call(ctx context.Context, method string, flags dbus.Flags, args ...interface{}) ([]interface{}, error) {
	debug("Adapter.call(ctx, %v, %v, %v)", method, flags, args)
	/*
		ch := make(chan *dbus.Call, 1)

		a.dobj.Go(method, flags, ch, args...)
		select {
		case <-ctx.Done():
			return nil, ErrContext
		case c := <-ch:
			return c.Body, c.Err
		}
	*/
	return call(ctx, a.dobj, method, flags, args...)
}

func (a *Adapter) SetProperty(ctx context.Context, key string, value interface{}) error {
	debug("Adapter.SetProperty(%v, %v)", key, value)

	const setProperty = "org.freedesktop.DBus.Properties.Set"

	_, err := a.call(ctx, setProperty, 0, "org.bluez.Adapter1", key, dbus.MakeVariant(value))

	return err
}

func (a *Adapter) SetDiscoveryFilter(ctx context.Context, filter map[string]interface{}) error {
	debug("Adapter.SetDiscoveryFilter(ctx, %v)", filter)

	const setDiscoveryFilter = "org.bluez.Adapter1.SetDiscoveryFilter"

	_, err := a.call(ctx, setDiscoveryFilter, 0, filter)

	return err
}

func (a *Adapter) StartDiscovery(ctx context.Context) ([]interface{}, error) {
	debug("Adapter.StartDiscovery(ctx)")

	const startDiscovery = "org.bluez.Adapter1.StartDiscovery"

	return a.call(ctx, startDiscovery, 0)
}

type ObjectManager struct {
	dobj dbus.BusObject
}

func NewObjectManager(conn *dbus.Conn, path string) *ObjectManager {
	debug("NewObjectManager(%v, %v)", conn, path)

	return &ObjectManager{
		dobj: conn.Object("org.bluez", dbus.ObjectPath(path)),
	}
}

func (o *ObjectManager) GetManagedObjects(ctx context.Context) []interface{} {
	const getManagedObjects = "org.freedesktop.DBus.ObjectManager.GetManagedObjects"

	objs, _ := call(ctx, o.dobj, getManagedObjects, 0)

	return objs
}
